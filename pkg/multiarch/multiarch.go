package multiarch

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/manifest"
	imagetypes "github.com/containers/image/v5/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ImagePlatform struct {
	ImageOS
	ImageCPUArchitecture
	ImageVarient
}
type ImageOS struct {
	slug string
}
type ImageVarient struct {
	slug string
}
type ImageCPUArchitecture struct {
	slug string
}
type CPUArchitecture struct {
	goArchitecture
	rpmArchitecture
}
type goArchitecture struct {
	slug string
}
type rpmArchitecture struct {
	slug string
}

var (
	CPUArchitectureAMD64   = CPUArchitecture{goArchitecture{architectureAMD64}, rpmArchitecture{architectureX8664}}
	CPUArchitectureARM64   = CPUArchitecture{goArchitecture{architectureARM64}, rpmArchitecture{architectureAARCH64}}
	CPUArchitecturePPC64LE = CPUArchitecture{goArchitecture{architecturePPC64LE}, rpmArchitecture{architecturePPC64LE}}
	CPUArchitectureS390X   = CPUArchitecture{goArchitecture{architectureS390X}, rpmArchitecture{architectureS390X}}
)

const (
	architectureX8664   = "x86_64"
	architectureAARCH64 = "aarch64"
	architectureAMD64   = "amd64"
	architectureARM64   = "arm64"
	architecturePPC64LE = "ppc64le"
	architectureS390X   = "s390x"
)

func (a CPUArchitecture) GoArch() string {
	return fmt.Sprintf("%s", a.goArchitecture.slug)
}
func (a CPUArchitecture) RPMArch() string {
	return fmt.Sprintf("%s", a.rpmArchitecture.slug)
}

func GetArchitecture(arch string) (CPUArchitecture, error) {

	switch arch {
	case architectureX8664, architectureAMD64:
		return CPUArchitectureAMD64, nil
	case architectureAARCH64, architectureARM64:
		return CPUArchitectureARM64, nil
	case architecturePPC64LE:
		return CPUArchitecturePPC64LE, nil
	case architectureS390X:
		return CPUArchitectureS390X, nil
	default:
		return CPUArchitecture{}, errors.Errorf("This is not a valid Architecture: %s", arch)
	}
}

// GetImagePlatforms returns an image's platforms linux/amd64, linux/arm/v7, linux/arm64/v8, linux/ppc64le, linux/s390x
func GetImagePlatforms(containerImage, pullSecret string) (platforms []ImagePlatform, isFatManifest bool, err error) {
	ctx := context.TODO()
	platforms = []ImagePlatform{}

	// Create the registry-config file
	ps, err := os.CreateTemp("", "registry-config")
	if err != nil {
		return []ImagePlatform{}, false, err
	}
	_, err = ps.Write([]byte(pullSecret))
	if err != nil {
		return []ImagePlatform{}, false, err
	}
	err = ps.Close()
	if err != nil {
		return []ImagePlatform{}, false, err
	}
	defer func() {
		os.Remove(ps.Name())
	}()
	// Parse the releaseImage
	ref, err := docker.ParseReference(fmt.Sprintf("//%s", containerImage))
	if err != nil {
		logrus.Warnf("Error parsing the image reference for the image: %v", err)
		return []ImagePlatform{}, false, err
	}
	sys := &imagetypes.SystemContext{
		AuthFilePath: ps.Name(),
		/*SystemRegistriesConfPath:  ...,
		SystemRegistriesConfDirPath: ...,
		SignaturePolicyPath:         ...,
		DockerPerHostCertDirPath:    ...,*/
		// TODO creating the above files could also address mirror registries (when reachable), custom CAs, and such...
	}
	src, err := ref.NewImageSource(ctx, sys)
	if err != nil {
		logrus.Warnf("Error creating the image source: %v", err)
		return []ImagePlatform{}, false, err
	}
	defer func(src imagetypes.ImageSource) {
		src.Close()
	}(src)

	// Get the raw manifest
	rawManifest, _, err := src.GetManifest(ctx, nil)
	if err != nil {
		logrus.Warnf("Error getting the image manifest: %v", err)
		return []ImagePlatform{}, false, err
	}

	// Verify if the release image is a fat manifest
	if manifest.MIMETypeIsMultiImage(manifest.GuessMIMEType(rawManifest)) {

		index, err := manifest.OCI1IndexFromManifest(rawManifest)
		if err != nil {
			logrus.Error(err, "Error parsing the OCI index from the raw manifest of the image")
		}
		for _, m := range index.Manifests {
			platform := ImagePlatform{ImageOS{m.Platform.OS}, ImageCPUArchitecture{m.Platform.Architecture}, ImageVarient{m.Platform.Variant}}
			platforms = append(platforms, platform)
		}
		return platforms, true, nil
	}

	// Parse the architecture of the non-manifest-list image.
	parsedImage, err := image.FromUnparsedImage(ctx, sys, image.UnparsedInstance(src, nil))
	if err != nil {
		logrus.Warnf("Error parsing the manifest of the image: %v", err)
	}
	config, err := parsedImage.OCIConfig(ctx)
	if err != nil {
		// Ignore errors due to invalid images at this stage
		logrus.Warnf("Error parsing the OCI config of the image: %v", err)
		return []ImagePlatform{}, false, err
	}

	platform := ImagePlatform{ImageOS{config.OS}, ImageCPUArchitecture{config.Architecture}, ImageVarient{config.Variant}}
	platforms = append(platforms, platform)
	return platforms, false, err
}
