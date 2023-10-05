package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/kvvit/deployment_cleaner/pkg/clientset"
	"github.com/kvvit/deployment_cleaner/pkg/deleteobjects"
	"github.com/kvvit/deployment_cleaner/pkg/loadvars"
	"github.com/kvvit/deployment_cleaner/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	port := "8080"
	log.Println("Cleaning service has been started")

	prometheus.MustRegister(deploymentMetrics)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Web server start listening on", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal("Error starting HTTP server:", err)
		}
	}()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	timeNow := metav1.Now()
	timeWork := timeNow.Hour()

	for {
		select {
		case <-ticker.C:
			metrics.FetchAndUpdateDeploymentMetrics(ctx, clientset, envvars.NameSpace, deploymentMetrics)
			if timeWork >= envvars.WorkStart && timeWork < envvars.WorkEnd {
				log.Println("Now is working time, pass changes")
			} else {
				deleteobjects.DeleteOldHelmReleases(ctx, clientset, envvars.TimeToDelete, envvars.DeploymentName, envvars.NameSpace, envvars.IsDryRun)
			}
		}
	}
}
