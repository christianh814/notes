package main

import (
	"fmt"
	"log"

	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

var TheURL string = "https://github.com/christianh814/example-openshift-go-repo/cluster-XXXX/bootstrap/overlays/default"

func main() {
	theYaml, err := RunKustomize(TheURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(theYaml))
}

// RunKustomize runs kustomize on a specific dir/url and returns a []byte
func RunKustomize(dir string) ([]byte, error) {
	// set up where to run kustomize, how to write it, and which file to create
	kustomizeDir := dir
	fSys := filesys.MakeFsOnDisk()

	// The default options are fine for our use case
	k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())

	// Run Kustomize
	m, err := k.Run(fSys, kustomizeDir)
	if err != nil {
		return nil, err
	}

	// try to Convert to YAML  and returl the results
	return m.AsYaml()
}
