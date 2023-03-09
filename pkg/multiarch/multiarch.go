package multiarch

import (
	"fmt"

	"github.com/pkg/errors"
)

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
