package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	ctx := context.Background()

	workStart := 10
	workEnd := 20

	//One day is 86400
	var seconds int64 = 86400

	varCheck := os.Getenv("NAMESPACES")
	if len(varCheck) == 0 {
		log.Fatalln("Please set NAMESPACES env variable")
	}
	nameSpaces := []string{os.Getenv("NAMESPACES")}

	clientset := getClientset()

	ticker := time.NewTicker(10 * time.Minute)
	for ; true; <-ticker.C {
		timeNow := metav1.Now()
		timeCheck := timeNow.Unix()
		timeWork := timeNow.Hour()
		if timeWork >= workStart && timeWork < workEnd {
			log.Println("Now is working time, pass changes")
		} else {
			for _, nameSpace := range nameSpaces {
				secrets, err := clientset.CoreV1().Secrets(nameSpace).List(ctx, metav1.ListOptions{})
				if err != nil {
					log.Fatalf("Error listing Secrets: %v\n", err)
				}
				for _, secret := range secrets.Items {
					if secret.Type == "helm.sh/release.v1" {
						version := (secret.GetLabels())["version"]
						if version == "1" {
							creationTime := (secret.CreationTimestamp).Unix()
							diffTime := timeCheck - creationTime
							log.Printf("Difference in time is %d for deployment %s\n", diffTime, secret.GetLabels()["name"])
							if diffTime >= seconds {
								log.Printf("Helm release %s is older than 24 hours and will be deleted\n", secret.GetLabels()["name"])
								nameCommon := (secret.GetLabels())["name"]
								objects := []string{"Deployments", "Secrets", "ConfigMaps", "Services", "Ingress"}
								for _, objectItem := range objects {
									switch objectItem {
									case "Deployments":
										deploymentList, err := clientset.AppsV1().Deployments(nameSpace).List(ctx, metav1.ListOptions{})
										if err != nil {
											log.Fatalf("Error listing Deployments: %v\n", err)
										}
										matchingObjects := []string{}
										for _, deployment := range deploymentList.Items {
											if strings.Contains(deployment.Name, nameCommon) {
												matchingObjects = append(matchingObjects, deployment.Name)
											}
										}

										for _, objectName := range matchingObjects {
											//clientset.AppsV1().Deployments(nameSpace).Delete(ctx, objectName, metav1.DeleteOptions{})
											log.Printf("Deployment name: %s in namespace: %s has been deleted\n", objectName, nameSpace)
										}
									case "Secrets":
										secretList, err := clientset.CoreV1().Secrets(nameSpace).List(ctx, metav1.ListOptions{})
										if err != nil {
											log.Fatalf("Error listing Secrets: %v\n", err)
										}
										matchingObjects := []string{}
										for _, secret := range secretList.Items {
											if strings.Contains(secret.Name, nameCommon) {
												matchingObjects = append(matchingObjects, secret.Name)
											}
										}
										for _, objectName := range matchingObjects {
											//clientset.CoreV1().Secrets(nameSpace).Delete(ctx, objectName, metav1.DeleteOptions{})
											log.Printf("Secret name: %s in namespace: %s has been deleted\n", objectName, nameSpace)
										}
									case "ConfigMaps":
										configMapList, err := clientset.CoreV1().ConfigMaps(nameSpace).List(ctx, metav1.ListOptions{})
										if err != nil {
											log.Fatalf("Error listing ConfigMaps: %v\n", err)
										}
										matchingObjects := []string{}
										for _, configMap := range configMapList.Items {
											if strings.Contains(configMap.Name, nameCommon) {
												matchingObjects = append(matchingObjects, configMap.Name)
											}
										}
										for _, objectName := range matchingObjects {
											//clientset.CoreV1().ConfigMaps(nameSpace).Delete(ctx, objectName, metav1.DeleteOptions{})
											log.Printf("ConfigMap name: %s in namespace: %s has been deleted\n", objectName, nameSpace)
										}
									case "Services":
										serviceList, err := clientset.CoreV1().Services(nameSpace).List(ctx, metav1.ListOptions{})
										if err != nil {
											log.Fatalf("Error listing Services: %v\n", err)
										}
										matchingObjects := []string{}
										for _, service := range serviceList.Items {
											if strings.Contains(service.Name, nameCommon) {
												matchingObjects = append(matchingObjects, service.Name)
											}
										}
										for _, objectName := range matchingObjects {
											//clientset.CoreV1().Services(nameSpace).Delete(ctx, objectName, metav1.DeleteOptions{})
											log.Printf("Service name: %s in namespace: %s has been deleted\n", objectName, nameSpace)
										}
									case "Ingress":
										ingressList, err := clientset.NetworkingV1().Ingresses(nameSpace).List(ctx, metav1.ListOptions{})
										if err != nil {
											log.Fatalf("Error listing Ingresses: %v\n", err)
										}
										matchingObjects := []string{}
										for _, ingress := range ingressList.Items {
											if strings.Contains(ingress.Name, nameCommon) {
												matchingObjects = append(matchingObjects, ingress.Name)
											}
										}
										for _, objectName := range matchingObjects {
											//clientset.NetworkingV1().Ingresses(nameSpace).Delete(ctx, objectName, metav1.DeleteOptions{})
											log.Printf("Ingress name: %s in namespace: %s has been deleted\n", objectName, nameSpace)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
