package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	var err error
	if err = exec.Command(
		"docker",
		"run",
		"-dt",
		"--rm",
		"-p",
		"8080:8080",
		"--name",
		"welcome-php",
		"--label",
		"chx=runner",
		"quay.io/redhatworkshops/welcome-php:latest",
	).Run(); err != nil {
		log.Fatal(err)
	}

	cmdOutPut, err := exec.Command(
		"docker",
		"ps",
		"--filter",
		"label=chx=runner",
		"--format",
		fmt.Sprintf(`table {{.Names}}\t{{.Image}}\t{{.Status}}`),
	).Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(cmdOutPut))

}
