package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"kopia-provisioner/version"

	_ "github.com/josephspurrier/goversioninfo"
)

func main() {
	v := version.Get()

	b := v.Version.Build
	v.Version.Build = v.Version.Build.ShortCommit()
	ver := v.Version.String()

	args := []string{
		"run",
		"github.com/josephspurrier/goversioninfo/cmd/goversioninfo",
		// General
		"-skip-versioninfo",
		"-o=../../version/versioninfo_windows.syso",
		// StringFileInfo.CompanyName
		"-company=ariaci",
		// StringFileInfo.InternalName
		"-internal-name=kopia-provisioner",
		// StringFileInfo.ProductName
		"-product-name=kopia-provisioner",
		// StringFileInfo.LegalCopyright
		"-copyright=© ariaci. Licensed under MIT.",
		// StringFileInfo.OriginalFilename
		"-original-name=kopia-provisioner.exe",
		// StringFileInfo.PrivateBuild
		fmt.Sprint("-private-build=", b),
		// FileVersion
		fmt.Sprintf("-ver-major=%d", v.Version.Core.Major),
		fmt.Sprintf("-ver-minor=%d", v.Version.Core.Minor),
		fmt.Sprintf("-ver-patch=%d", v.Version.Core.Patch),
		"-ver-build=0",
		// ProductVersion
		fmt.Sprintf("-product-ver-major=%d", v.Version.Core.Major),
		fmt.Sprintf("-product-ver-minor=%d", v.Version.Core.Minor),
		fmt.Sprintf("-product-ver-patch=%d", v.Version.Core.Patch),
		"-product-ver-build=0",
		// StringFileInfo.FileVersion
		fmt.Sprint("-file-version=", ver),
		// StringFileInfo.ProductVersion
		fmt.Sprint("-product-version=", ver),
		// StringFileInfo.FileDescription
		fmt.Sprint("-description=kopia-provisioner ", ver),
	}

	if runtime.GOARCH != "386" {
		args = append(args, "-64")
	}

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
