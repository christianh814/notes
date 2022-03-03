package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	goyaml "gopkg.in/yaml.v2"
)

var TheYaml string = `apiVersion: v1
kind: Namespace
metadata:
  name: bgd
---
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: bgd
  name: bgd
spec:
  replicas: 1
  selector:
    matchLabels:
      app: bgd
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: bgd
    spec:
      containers:
      - image: quay.io/redhatworkshops/bgd:latest
        name: bgd
        env:
        - name: COLOR
          value: "blue"
        env:
        - name: BREAKAGE
          value: "trying --- to break things"
        resources: {}
---
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    app: bgd
  name: bgd
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: bgd
`

func main() {
	splittedYaml, err := SplitYAML([]byte(TheYaml))
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(splittedYaml[0]))

	/*
		for _, v := range splittedYaml {
			fmt.Println(string(v))
		}
	*/

}

func SplitYAML(resources []byte) ([][]byte, error) {

	dec := goyaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := goyaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}
