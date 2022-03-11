package main

import (
	"context"
	"flag"
	"log"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeConfig = flag.String("kubeconfig", "", "kubeconfig file")
	flag.Parse()
	client, err := newClient(*kubeConfig)
	if err != nil {
		log.Fatal(err)
	}
	// set the variables
	ns := "openshift-ingress"
	deployment := "router-default"
	ctx := context.TODO()

	// Load the router deployment and get the env var
	rd, err := client.AppsV1().Deployments(ns).Get(ctx, deployment, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	ev := []corev1.EnvVar{
		{
			Name:  "ROUTER_SUBDOMAIN",
			Value: "${name}-${namespace}.apps.127.0.0.1.nip.io",
		},
		{
			Name:  "ROUTER_ALLOW_WILDCARD_ROUTES",
			Value: "true",
		},
		{
			Name:  "ROUTER_OVERRIDE_HOSTNAME",
			Value: "true",
		},
	}

	ev = append(ev, rd.Spec.Template.Spec.Containers[0].Env...)

	rd.Spec.Template.Spec.Containers[0].Env = ev
	client.AppsV1().Deployments(ns).Update(ctx, rd, metav1.UpdateOptions{})

}

func newClient(kubeConfigPath string) (kubernetes.Interface, error) {
	if kubeConfigPath == "" {
		kubeConfigPath = os.Getenv("KUBECONFIG")
	}
	if kubeConfigPath == "" {
		kubeConfigPath = clientcmd.RecommendedHomeFile // use default path(.kube/config)
	}
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(kubeConfig)
}

func newDynamicClient(kubeConfigPath string) (dynamic.Interface, error) {
	if kubeConfigPath == "" {
		kubeConfigPath = os.Getenv("KUBECONFIG")
	}
	if kubeConfigPath == "" {
		kubeConfigPath = clientcmd.RecommendedHomeFile // use default path(.kube/config)
	}
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(kubeConfig)
}
