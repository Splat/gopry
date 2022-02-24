package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

type sysPkg struct {
	ID          int
	Name        string
	Description string
	Version     string
}

func GetSystemPkgs() ([]sysPkg, error) {
	var sysPkgs []sysPkg

	pkgConfigCmd := exec.Command("pkg-config", "--list-all")

	if pkgConfigStdOut, err := pkgConfigCmd.Output(); err != nil {
		fmt.Println(err.Error())
		return nil, err
	} else {
		// extra empty line at the end of pkg-config output
		pkgConfigStdOut = pkgConfigStdOut[:len(pkgConfigStdOut)-1]

		for idx, pkg := range strings.Split(string(pkgConfigStdOut), "\n") {
			// gotta handle some weird formatting from stdout
			pkgInfoRawStr := strings.Split(pkg, " ")

			pkgName, pkgDesc := pkgInfoRawStr[0], pkgInfoRawStr[1:]

			if pkgVersion, versionErr := GetSystemPkgVersion(pkgName); versionErr != nil {
				return nil, versionErr
			} else {
				sysPkgs = append(sysPkgs, sysPkg{
					ID:          idx,
					Name:        pkgName,
					Description: strings.TrimSpace(strings.Join(pkgDesc, "")),
					Version:     pkgVersion,
				})
			}
		}

		return sysPkgs, nil
	}
}

func GetSystemPkgVersion(pkgName string) (string, error) {
	//versionCmd := exec.Command("pkg-config", "--print-provides", fmt.Sprintf("\"%s\"", pkgName))
	versionCmd := exec.Command("pkg-config", "--print-provides", pkgName)

	if versionOutput, err := versionCmd.Output(); err != nil {
		fmt.Println(pkgName)
		fmt.Println(err.Error())
		return "", err
	} else {
		versionStr := string(versionOutput)
		versionStr = strings.ReplaceAll(versionStr, "= ", "")
		versionStr = strings.ReplaceAll(versionStr, "\n", "")

		return versionStr, err
	}
}
