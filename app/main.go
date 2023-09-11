package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	ctx := context.TODO()

	workStart := 10
	workEnd := 20

	//One day is 86400
	var seconds int64 = 86400

	varCheck := os.Getenv("NAMESPACES")
	if len(varCheck) == 0 {
		fmt.Println("Please set NAMESPACES env variable")
		return
	}
	nameSpaces := []string{os.Getenv("NAMESPACES")}

	configPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	var config *rest.Config
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("failed to build local config: %s\n", err)
			return
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			log.Fatalf("failed to build in-cluster config: %s\n", err)
			return
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create endpoints watcher: %s\n", err)
		return
	}

	for {
		timeNow := metav1.Now()
		timeCheck := timeNow.Unix()
		timeWork := timeNow.Hour()
		if timeWork < workStart || timeWork >= workEnd {
			for _, nameSpace := range nameSpaces {
				secrets, err := clientset.CoreV1().Secrets(nameSpace).List(ctx, metav1.ListOptions{})
				if err != nil {
					fmt.Printf("Error listing Secrets: %v\n", err)
					return
				}
				for _, secret := range secrets.Items {
					if secret.Type == "helm.sh/release.v1" {
						version := (secret.GetLabels())["version"]
						if version == "1" {
							creationTime := (secret.CreationTimestamp).Unix()
							diffTime := timeCheck - creationTime
							log.Printf("Difference in time is %d for deployment %s\n", diffTime, secret.GetLabels()["name"])
							if diffTime >= seconds {
								fmt.Println("Deployment is older than 24 hours and will be deleted")
								nameCommon := (secret.GetLabels())["name"]
								fmt.Println(nameCommon)
								objects := []string{"Deployments", "Secrets", "ConfigMaps", "Services", "Ingress"}
								for _, objectItem := range objects {
									switch objectItem {
									case "Deployments":
										deploymentList, err := clientset.AppsV1().Deployments(nameSpace).List(ctx, metav1.ListOptions{})
										if err != nil {
											fmt.Printf("Error listing Deployments: %v\n", err)
											return
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
											fmt.Printf("Error listing Secrets: %v\n", err)
											return
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
											fmt.Printf("Error listing ConfigMaps: %v\n", err)
											return
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
											fmt.Printf("Error listing Services: %v\n", err)
											return
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
											fmt.Printf("Error listing Ingresses: %v\n", err)
											return
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
		} else {
			log.Println("Now is working time, pass changes")
		}
		time.Sleep(600 * time.Second)
	}
}
