package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"kopia-provisioner/version"
)

func createVersionInfoJSON() string {
	v := version.Get()

	b := v.Version.Build
	v.Version.Build = v.Version.Build.ShortCommit()
	ver := v.Version.String()

	content := fmt.Sprintf(`{
  "RT_VERSION": {
    "#1": {
      "0000": {
        "fixed": {
          "file_version": "%d.%d.%d.0",
          "product_version": "%d.%d.%d.0"
        },
        "info": {
          "0409": {
            "CompanyName": "ariaci",
            "FileDescription": "kopia-provisioner %s",
            "FileVersion": "%s",
            "InternalName": "kopia-provisioner",
            "LegalCopyright": "© ariaci. Licensed under MIT.",
            "OriginalFilename": "kopia-provisioner.exe",
            "PrivateBuild": "%s",
            "ProductName": "kopia-provisioner",
            "ProductVersion": "%s"
          }
        }
      }
    }
  }
}`,
		v.Version.Core.Major, v.Version.Core.Minor, v.Version.Core.Patch,
		v.Version.Core.Major, v.Version.Core.Minor, v.Version.Core.Patch,
		ver, ver, b, ver)

	tmpFile, err := os.CreateTemp("", "kopia-provisioner-versioninfo-*.winres.json")
	if err != nil {
		panic(err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		panic(err)
	}

	path, err := filepath.Abs(tmpFile.Name())
	if err != nil {
		panic(err)
	}

	return path
}

func main() {
	versioninfo_json := createVersionInfoJSON()
	defer os.Remove(versioninfo_json)

	args := []string{
		"run",
		"github.com/tc-hib/go-winres",
		"make",
		"--in", versioninfo_json,
		"--out", "../../version/winres",
	}

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
