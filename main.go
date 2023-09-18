package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kvvit/deployment_cleaner/pkg/clientset"
	"github.com/kvvit/deployment_cleaner/pkg/deleteobjects"
	"github.com/kvvit/deployment_cleaner/pkg/loadvars"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	ctx := context.Background()
	envvars := loadvars.LoadVars()

	fmt.Println(envvars)

	workStart := 10
	workEnd := 19

	//One day is 86400
	var timeToDelete int64 = 86400

	varCheck := os.Getenv("NAMESPACE")
	if len(varCheck) == 0 {
		log.Fatalln("Please set NAMESPACE env variable")
	}
	nameSpace := os.Getenv("NAMESPACE")

	clientset := clientset.GetClientset()

	ticker := time.NewTicker(10 * time.Minute)
	for ; true; <-ticker.C {
		timeNow := metav1.Now()
		timeWork := timeNow.Hour()
		if timeWork >= workStart && timeWork < workEnd {
			log.Println("Now is working time, pass changes")
		} else {
			deleteobjects.DeleteOldHelmReleases(ctx, clientset, nameSpace, timeToDelete)
		}
	}
}
