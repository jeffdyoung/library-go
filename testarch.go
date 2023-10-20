package main

import (
	"fmt"
	"os"

	"github.com/openshift/library-go/pkg/multiarch"
)

func main() {

	fileBytes, err := os.ReadFile("/run/user/1000/containers/auth.json")
	if err != nil {
		panic(err)
	}
	pullSecret := string(fileBytes)
	image := "docker.io/library/ubuntu"
	image = "quay.io/openshift-release-dev/ocp-release:4.15.0-ec.1-multi"
	// image = "quay.io/fedora/fedora"
	// image = "quay.io/multi-arch/sushy-tools:arm"
	// image = "quay.io/centos/centos"

	fmt.Println("image ", image)
	platforms, mflist, err := multiarch.GetImagePlatforms(image, pullSecret)
	fmt.Println("errors ", err)
	fmt.Println("platforms ", platforms)

	if mflist {
		fmt.Println("is MF listed")
	} else {
		fmt.Println("not MF listed")
	}
}
