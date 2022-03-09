package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

var MyDeploy string = `apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: mydeploy
  name: mydeploy
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mydeploy
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: mydeploy
    spec:
      containers:
      - image: quay.io/redhatworkshops/welcome-php
        name: welcome-php
        resources: {}
`

func main() {
	// getting new client
	var kubeConfig = flag.String("kubeconfig", "", "kubeconfig file")
	flag.Parse()
	client, err := newClient(*kubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	// just testing things for dynbamic in case I need it
	dynamicClient, err := newDynamicClient(*kubeConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dynamicClient)

	// create the deployment via DoSSA and wait until it's running
	ns := "default"
	depl := "mydeploy"

	// DoSSA requires a pointer to a rest config...need to change that one day
	rc, err := newRestConfig(*kubeConfig)
	if err != nil {
		log.Fatal(err)
	}
	DoSSA(context.TODO(), rc, []byte(MyDeploy))

	// wait for the deployment
	if err = waitForDeployment(client, ns, depl, 90*time.Second); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Running")
	}
}

func newRestConfig(kubeConfigPath string) (*rest.Config, error) {
	if kubeConfigPath == "" {
		kubeConfigPath = os.Getenv("KUBECONFIG")
	}
	if kubeConfigPath == "" {
		kubeConfigPath = clientcmd.RecommendedHomeFile // use default path(.kube/config)
	}
	kc, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	return kc, nil
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

//check to see if the named deployment is running
func isDeploymentRunning(c kubernetes.Interface, ns string, depl string) wait.ConditionFunc {

	return func() (bool, error) {

		dep, err := c.AppsV1().Deployments(ns).Get(context.TODO(), depl, v1.GetOptions{})
		if err != nil {
			return false, err
		}

		if dep.Status.ReadyReplicas == 0 {
			fmt.Println("waiting")
			return false, nil
		}

		return true, nil

	}
}

// Poll up to timeout seconds for pod to enter running state.
func waitForDeployment(c kubernetes.Interface, namespace string, deployment string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, isDeploymentRunning(c, namespace, deployment))
}

// DoSSA  does service side apply with the given YAML as a []byte
func DoSSA(ctx context.Context, cfg *rest.Config, yaml []byte) error {
	// Read yaml into a slice of byte
	yml := yaml

	// get the RESTMapper for the GVR
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// create dymanic client
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return err
	}

	// read YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(yml, nil, obj)
	if err != nil {
		return err
	}

	// Get the GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// Get the REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	// Create object into JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// Create or Update the obj with service side apply
	//     types.ApplyPatchType indicates service side apply
	//     FieldManager specifies the field owner ID.
	_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, v1.PatchOptions{
		FieldManager: "doSSAApply",
	})

	return err
}
