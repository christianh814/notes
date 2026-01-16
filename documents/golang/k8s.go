package main

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {

	/*
	   // Create in-cluster config:
	       `InClusterConfig()`` automatically:
	         * Detects it is running in a Pod
	         * Uses the mounted ServiceAccount token
	         * Uses the cluster CA
	         * Talks to https://kubernetes.default.svc
	       No extra configuration required.
	*/
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// List pods in the current namespace: Service account of the pod must have permissions to list pods
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d pods in the namespace:\n", len(pods.Items))
	for _, pod := range pods.Items {
		log.Printf("Pod: %s, Status: %s\n", pod.Name, pod.Status.Phase)
	}
}
