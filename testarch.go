package main

import (
	"fmt"
	"os"

	"github.com/openshift/library-go/pkg/multiarch"
)

func main() {

	pullSecretFile, err := os.ReadFile("/run/user/1000/containers/auth.json")
	if err != nil {
		panic(err)
	}
	pullSecret := string(pullSecretFile)

	myvarient := "linux/amd64"

	images := []string{"docker.io/library/ubuntu",
		"quay.io/openshift-release-dev/ocp-release:4.15.0-ec.1-multi",
		"quay.io/fedora/fedora",
		"quay.io/multi-arch/sushy-tools:arm",
		"docker.io/library/busybox",
		"mcr.microsoft.com/windows/servercore:ltsc2022",
		"quay.io/multi-arch/doesntexist:badimage",
	}
	for _, image := range images {
		fmt.Println("image ", image)
		platforms, err := multiarch.GetImagePlatforms(image, pullSecret)
		if err != nil {
			fmt.Println("errors ", err)
		}

		// fmt.Println("platforms ", platforms)
		for _, platform := range platforms {
			fmt.Println("\t" + platform.GetFullVariant())
			if platform.GetShortVariant() == myvarient {
				fmt.Println("\t\tthis works with my os/arch: ", myvarient)
				fmt.Println("\t\tgoArch: ", platform.CPUArch.GoArch())
				fmt.Println("\t\trpmArch: ", platform.CPUArch.RPMArch())
			}
		}
	}
}
