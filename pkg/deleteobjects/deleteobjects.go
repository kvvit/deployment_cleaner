package deleteobjects

import (
	"context"
	"log"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DeleteOldHelmReleases(ctx context.Context, clientset *kubernetes.Clientset, namespace string, timeToDelete int64) {
	secrets, err := clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing Secrets: %v\n", err)
	}

	for _, secret := range secrets.Items {
		if secret.Type != "helm.sh/release.v1" || secret.Labels["version"] != "1" {
			continue
		}
		if int64(time.Now().Sub(secret.CreationTimestamp.Time).Seconds()) < timeToDelete {
			continue
		}
		log.Printf("Helm release %s is older than 24 hours and will be deleted\n", secret.Labels["name"])
		DeleteObjectsWithCommonName(ctx, clientset, namespace, secret.Labels["name"])
	}
}

func DeleteObjectsWithCommonName(ctx context.Context, clientset *kubernetes.Clientset, namespace, commonName string) {
	objectTypes := []string{"Deployments", "Secrets", "ConfigMaps", "Services", "Ingress"}

	for _, objectType := range objectTypes {
		switch objectType {
		case "Deployments":
			deploymentList, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				log.Fatalf("Error listing Deployments: %v\n", err)
			}
			matchingObjects := []string{}
			for _, deployment := range deploymentList.Items {
				if strings.Contains(deployment.Name, commonName) {
					matchingObjects = append(matchingObjects, deployment.Name)
				}
			}

			for _, objectName := range matchingObjects {
				//clientset.AppsV1().Deployments(namespace).Delete(ctx, objectName, metav1.DeleteOptions{})
				log.Printf("Deployment name: %s in namespace: %s has been deleted\n", objectName, namespace)
			}
		case "Secrets":
			secretList, err := clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				log.Fatalf("Error listing Secrets: %v\n", err)
			}
			matchingObjects := []string{}
			for _, secret := range secretList.Items {
				if strings.Contains(secret.Name, commonName) {
					matchingObjects = append(matchingObjects, secret.Name)
				}
			}
			for _, objectName := range matchingObjects {
				//clientset.CoreV1().Secrets(namespace).Delete(ctx, objectName, metav1.DeleteOptions{})
				log.Printf("Secret name: %s in namespace: %s has been deleted\n", objectName, namespace)
			}
		case "ConfigMaps":
			configMapList, err := clientset.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				log.Fatalf("Error listing ConfigMaps: %v\n", err)
			}
			matchingObjects := []string{}
			for _, configMap := range configMapList.Items {
				if strings.Contains(configMap.Name, commonName) {
					matchingObjects = append(matchingObjects, configMap.Name)
				}
			}
			for _, objectName := range matchingObjects {
				//clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, objectName, metav1.DeleteOptions{})
				log.Printf("ConfigMap name: %s in namespace: %s has been deleted\n", objectName, namespace)
			}
		case "Services":
			serviceList, err := clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				log.Fatalf("Error listing Services: %v\n", err)
			}
			matchingObjects := []string{}
			for _, service := range serviceList.Items {
				if strings.Contains(service.Name, commonName) {
					matchingObjects = append(matchingObjects, service.Name)
				}
			}
			for _, objectName := range matchingObjects {
				//clientset.CoreV1().Services(namespace).Delete(ctx, objectName, metav1.DeleteOptions{})
				log.Printf("Service name: %s in namespace: %s has been deleted\n", objectName, namespace)
			}
		case "Ingress":
			ingressList, err := clientset.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				log.Fatalf("Error listing Ingresses: %v\n", err)
			}
			matchingObjects := []string{}
			for _, ingress := range ingressList.Items {
				if strings.Contains(ingress.Name, commonName) {
					matchingObjects = append(matchingObjects, ingress.Name)
				}
			}
			for _, objectName := range matchingObjects {
				//clientset.NetworkingV1().Ingresses(namespace).Delete(ctx, objectName, metav1.DeleteOptions{})
				log.Printf("Ingress name: %s in namespace: %s has been deleted\n", objectName, namespace)
			}
		}
	}
}
