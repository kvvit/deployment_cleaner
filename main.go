package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kvvit/deployment_cleaner/pkg/clientset"
	"github.com/kvvit/deployment_cleaner/pkg/deleteobjects"
	"github.com/kvvit/deployment_cleaner/pkg/loadvars"
	"github.com/kvvit/deployment_cleaner/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/apimachinery/pkg/util/wait"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	deploymentMetrics = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deployment_live_time_seconds",
			Help: "Live time of deployments in seconds",
		},
		[]string{"deployment"},
	)
)

func main() {
	ctx := context.Background()
	envvars := loadvars.LoadVars()
	clientset := clientset.GetClientset()
	var wg sync.WaitGroup
	log.Println("Cleaning service has been started")

	prometheus.MustRegister(deploymentMetrics)

	ticker := time.NewTicker(1 * time.Minute)
	for ; true; <-ticker.C {
		wg.Add(6)
		timeNow := metav1.Now()
		timeWork := timeNow.Hour()
		go wait.Forever(func() {
			metrics.FetchAndUpdateDeploymentMetrics(ctx, clientset, envvars.NameSpace, deploymentMetrics)
			wg.Done()
		}, 30*time.Second)
		port := "8080"
		if timeWork >= envvars.WorkStart && timeWork < envvars.WorkEnd {
			log.Println("Now is working time, pass changes")
		} else {
			go deleteobjects.DeleteOldHelmReleases(ctx, clientset, envvars.TimeToDelete, envvars.DeploymentName, envvars.NameSpace, envvars.IsDryRun)
			wg.Done()
		}
		http.Handle("/metrics", promhttp.Handler())
		log.Println("listening on", port)
		go log.Fatal(http.ListenAndServe(":"+port, nil))
		wg.Done()
	}
	wg.Wait()
}
