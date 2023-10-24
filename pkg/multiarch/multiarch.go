package multiarch

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/manifest"
	imagetypes "github.com/containers/image/v5/types"
	"github.com/sirupsen/logrus"
)

type ImagePlatform struct {
	PlatformOS
	CPUArch
	CPUArchVarient
}
type PlatformOS struct {
	slug string
}
type CPUArchVarient struct {
	slug string
}
type CPUArch struct {
	goArch
	rpmArch
}
type goArch struct {
	slug string
}
type rpmArch struct {
	slug string
}

var (
	CPUArchAMD64 = CPUArch{goArch{AMD64}, rpmArch{X8664}}
	CPUArchARM64 = CPUArch{goArch{ARM64}, rpmArch{AARCH64}}
	CPUArch386   = CPUArch{goArch{_386}, rpmArch{I386}}
)

const (
	// GOARCH values from:https://go.dev/src/go/build/syslist.go
	AMD64       = "amd64"
	ARM64       = "arm64"
	_386        = "386"
	AMD64P32    = "amd64p32"
	ARM         = "arm"
	ARMBE       = "armbe"
	ARM64BE     = "arm64be"
	LOONG64     = "loong64"
	MIPS        = "mips"
	MIPSLE      = "mipsle"
	MIPS64      = "mips64"
	MIPS64LE    = "mips64le"
	MIPS64P32   = "mips64p32"
	MIPS64P32LE = "mips64p32le"
	PPC         = "ppc"
	PPC64       = "ppc64"
	PPC64LE     = "ppc64le"
	RISCV       = "riscv"
	RISCV64     = "riscv64"
	S390        = "s390"
	S390X       = "s390x"
	SPARC       = "sparc"
	SPARC64     = "sparc64"
	WASM        = "wasm"

	// RPMArch (aka 'uname -m') values that are different from GOARCH
	// https://github.com/torvalds/linux/tree/master/arch look for values of UTS_MACHINE
	// similar maping from deb to rpm arch here: https://github.com/torvalds/linux/blob/master/scripts/package/mkdebian#L21
	X8664   = "x86_64"
	AARCH64 = "aarch64"
	I386    = "i386"

	// GOOS values from: https://go.dev/src/go/build/syslist.go
	AIX       = "aix"
	ANDROID   = "android"
	DARWIN    = "darwin"
	DRAGONFLY = "dragonfly"
	FREEBSD   = "freebsd"
	HURD      = "hurd"
	ILLUMOS   = "illumos"
	IOS       = "ios"
	JS        = "js"
	LINUX     = "linux"
	NACL      = "nacl"
	NETBSD    = "netbsd"
	OPENBSD   = "openbsd"
	PLAN9     = "plan9"
	SOLARIS   = "solaris"
	WASIP1    = "wasip1"
	WINDOWS   = "windows"
	Z0S       = "zos"
)

func (image ImagePlatform) GetFullVariant() string {
	fullVariant := image.GetShortVariant()
	if image.Varient() != "" {
		fullVariant = fullVariant + "/" + image.Varient()
	}
	return fullVariant
}
func (image ImagePlatform) GetShortVariant() string {
	return image.OS() + "/" + image.CPUArch.GoArch()
}
func (image ImagePlatform) OS() string {
	return fmt.Sprintf("%s", image.PlatformOS.slug)
}
func (image ImagePlatform) Varient() string {
	return fmt.Sprintf("%s", image.CPUArchVarient.slug)
}
func (a CPUArch) GoArch() string {
	return fmt.Sprintf("%s", a.goArch.slug)
}
func (a CPUArch) RPMArch() string {
	return fmt.Sprintf("%s", a.rpmArch.slug)
}

func GetOs(os string) (PlatformOS, error) {
	switch os {
	case AIX, ANDROID, DARWIN, DRAGONFLY, FREEBSD, HURD, ILLUMOS, IOS, JS, LINUX, NACL, NETBSD, OPENBSD, PLAN9, SOLARIS, WASIP1, WINDOWS, Z0S:
		return PlatformOS{os}, nil
	default:
		return PlatformOS{}, fmt.Errorf("%s is not a known OS", os)
	}
}

// GetArchitecture: Returns a CPUARCH object with goArch, and rpmArch
func GetCPUArch(arch string) (CPUArch, error) {
	switch arch {
	case X8664, AMD64:
		return CPUArchAMD64, nil
	case AARCH64, ARM64:
		return CPUArchARM64, nil
	case _386, I386:
		return CPUArch386, nil
	case AMD64P32, ARM, ARMBE, ARM64BE, LOONG64, MIPS, MIPSLE, MIPS64, MIPS64LE, MIPS64P32, MIPS64P32LE, PPC, PPC64, PPC64LE, RISCV, RISCV64, S390, S390X, SPARC, SPARC64, WASM:
		return CPUArch{goArch{arch}, rpmArch{arch}}, nil
	default:
		return CPUArch{}, fmt.Errorf("%s is not a known CPUArch", arch)
	}
}

// GetImagePlatforms returns an image's platforms linux/amd64, linux/arm/v7, linux/arm64/v8, linux/ppc64le, linux/s390x
func GetImagePlatforms(containerImage, pullSecret string) (platforms []ImagePlatform, err error) {
	ctx := context.TODO()
	platforms = []ImagePlatform{}

	// Create the registry-config file
	ps, err := os.CreateTemp("", "registry-config")
	if err != nil {
		return []ImagePlatform{}, err
	}
	_, err = ps.Write([]byte(pullSecret))
	if err != nil {
		return []ImagePlatform{}, err
	}
	err = ps.Close()
	if err != nil {
		return []ImagePlatform{}, err
	}
	defer func() {
		os.Remove(ps.Name())
	}()
	// Parse the releaseImage
	ref, err := docker.ParseReference(fmt.Sprintf("//%s", containerImage))
	if err != nil {
		logrus.Warnf("Error parsing the image reference for the image: %v", err)
		return []ImagePlatform{}, err
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
		return []ImagePlatform{}, err
	}
	defer func(src imagetypes.ImageSource) {
		src.Close()
	}(src)

	// Get the raw manifest
	rawManifest, _, err := src.GetManifest(ctx, nil)
	if err != nil {
		logrus.Warnf("Error getting the image manifest: %v", err)
		return []ImagePlatform{}, err
	}

	// Verify if the release image is a fat manifest
	if manifest.MIMETypeIsMultiImage(manifest.GuessMIMEType(rawManifest)) {

		index, err := manifest.OCI1IndexFromManifest(rawManifest)
		if err != nil {
			logrus.Error(err, "Error parsing the OCI index from the raw manifest of the image")
		}
		for _, m := range index.Manifests {
			arch, err := GetCPUArch(m.Platform.Architecture)
			if err != nil {
				return []ImagePlatform{}, err
			}
			os, err := GetOs(m.Platform.OS)
			if err != nil {
				return []ImagePlatform{}, err
			}
			platform := ImagePlatform{os, arch, CPUArchVarient{m.Platform.Variant}}
			platforms = append(platforms, platform)
		}
		return platforms, nil
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
		return []ImagePlatform{}, err
	}
	arch, err := GetCPUArch(config.Architecture)
	if err != nil {
		return []ImagePlatform{}, err
	}
	os, err := GetOs(config.OS)
	if err != nil {
		return []ImagePlatform{}, err
	}
	platform := ImagePlatform{os, arch, CPUArchVarient{config.Variant}}
	platforms = append(platforms, platform)
	return platforms, err
}
