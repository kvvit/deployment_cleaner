package metrics

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchAndUpdateDeploymentMetrics(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
	deploymentMetrics *prometheus.GaugeVec) {
	secrets, err := clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Println("Error fetching secrets:", err)
		return
	}

	for _, secret := range secrets.Items {
		if secret.Type != "helm.sh/release.v1" || secret.Labels["version"] != "1" {
			continue
		}
		deploymentMetrics.WithLabelValues(secret.Labels["name"]).Set(float64(time.Now().Sub(secret.CreationTimestamp.Time).Seconds()))
	}
}
