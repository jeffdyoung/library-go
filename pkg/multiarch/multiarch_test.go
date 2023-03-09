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
