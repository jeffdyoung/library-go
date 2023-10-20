package multiarch

import (
	"testing"
)

func TestInvalidArch(t *testing.T) {
	testarch := "x86"
	cpuarch, err := GetArchitecture(testarch)
	if err == nil {
		t.Errorf("Shouldn't create CPUArchitecture from string %s", testarch)
	}
	t.Logf("Can't create CPUArchitecture %s from string %s", cpuarch, testarch)
}
func TestEmptyArch(t *testing.T) {
	testarch := ""
	cpuarch, err := GetArchitecture(testarch)
	if err == nil {
		t.Errorf("Shouldn't create CPUArchitecture from string %s", testarch)
	}
	t.Logf("Can't create CPUArchitecture %s from string %s", cpuarch, testarch)
}
func TestCreateArch(t *testing.T) {
	rpmarch := []string{"x86_64", "aarch64", "s390x", "ppc64le"}
	goarch := []string{"amd64", "arm64", "s390x", "ppc64le"}

	for i, s := range rpmarch {
		cpuarch, err := GetArchitecture(s)
		if err != nil {
			t.Errorf("Can't create CPUArchiture from string %s", s)
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
		cpuarch, err := GetArchitecture(s)
		if err != nil {
			t.Errorf("Can't create CPUArchiture from string %s", s)
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
	arch, mflist, err := GetImagePlatforms(maImage, pullSecret)
	if err != nil {
		t.Errorf("maImage got %s expected %s", arch, mflist, err)
	}

	// Single arch x86
	arch, mflist, err = GetImagePlatforms(x86Image, pullSecret)
	if err != nil {
		t.Errorf("x86Image %s %s", arch, mflist, err)
	}

	// Single arch arm64
	arch, mflist, err = GetImagePlatforms(armImage, pullSecret)
	if err != nil {
		t.Logf("armImage %s %b %s", arch, mflist, err)
	}
}
