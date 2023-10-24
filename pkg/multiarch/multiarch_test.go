package multiarch

import (
	"testing"
)

func TestInvalidArch(t *testing.T) {
	testarch := "x86"
	cpuarch, err := GetCPUArch(testarch)
	if err == nil {
		t.Errorf("testarch %s, should not be a valid CPUArch %v", testarch, cpuarch)
	}
}

func TestEmptyArch(t *testing.T) {
	testarch := ""
	cpuarch, err := GetCPUArch(testarch)
	if err == nil {
		t.Errorf("testarch %s, should not be a valid CPUArch %v", testarch, cpuarch)
	}
}
func TestCreateArch(t *testing.T) {
	rpmarch := []string{"x86_64", "aarch64", "s390x", "ppc64le", "i386"}
	goarch := []string{"amd64", "arm64", "s390x", "ppc64le", "386"}

	for i, s := range rpmarch {
		cpuarch, err := GetCPUArch(s)
		if err != nil {
			t.Errorf("%s", err)
		}
		wantrpm := rpmarch[i]
		wantgo := goarch[i]

		if wantrpm != cpuarch.RPMArch() {
			t.Errorf("RPMArch fails %s doesn't match %s ", wantrpm, cpuarch.RPMArch())
		}
		if wantgo != cpuarch.GoArch() {
			t.Errorf("GoArch fails %s doesn't match %s ", wantrpm, cpuarch.GoArch())
		}
	}

	for i, s := range goarch {
		cpuarch, err := GetCPUArch(s)
		if err != nil {
			t.Errorf("%s", err)
		}
		wantrpm := rpmarch[i]
		wantgo := goarch[i]

		if wantrpm != cpuarch.RPMArch() {
			t.Errorf("RPMArch fails %s doesn't match %s ", wantrpm, cpuarch.RPMArch())
		}
		if wantgo != cpuarch.GoArch() {
			t.Errorf("GoArch fails %s doesn't match %s ", wantrpm, cpuarch.GoArch())
		}
	}
}
func TestImageArch(t *testing.T) {
	maImage := "quay.io/multi-arch/sushy-tools:muiltarch"
	x86Image := "quay.io/multi-arch/sushy-tools:x86"
	armImage := "quay.io/multi-arch/sushy-tools:arm"

	pullSecret := "{}"
	// MF x86 + arm64
	arch, err := GetImagePlatforms(maImage, pullSecret)
	if err != nil {
		t.Errorf("maImage: %s %s ", arch, err)
	}

	// Single arch x86
	arch, err = GetImagePlatforms(x86Image, pullSecret)
	if err != nil {
		t.Errorf("x86Image: %s %s ", arch, err)
	}

	// Single arch arm64
	arch, err = GetImagePlatforms(armImage, pullSecret)
	if err != nil {
		t.Logf("armImage: %s %s ", arch, err)
	}
}
