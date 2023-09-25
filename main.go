package main

import (
	"context"
	"log"
	"time"

	"github.com/kvvit/deployment_cleaner/pkg/clientset"
	"github.com/kvvit/deployment_cleaner/pkg/deleteobjects"
	"github.com/kvvit/deployment_cleaner/pkg/loadvars"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	ctx := context.Background()
	envvars := loadvars.LoadVars()
	clientset := clientset.GetClientset()
	log.Println("Cleaning service has been started")

	ticker := time.NewTicker(10 * time.Minute)
	for ; true; <-ticker.C {
		timeNow := metav1.Now()
		timeWork := timeNow.Hour()
		if timeWork >= envvars.WorkStart && timeWork < envvars.WorkEnd {
			log.Println("Now is working time, pass changes")
		} else {
			deleteobjects.DeleteOldHelmReleases(ctx, clientset, envvars.NameSpace, envvars.TimeToDelete, envvars.DeploymentName, envvars.IsDryRun)
		}
	}
}
