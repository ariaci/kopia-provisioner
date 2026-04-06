package main

import (
	"fmt"
	"log"
	"os"

	"kopia-provisioner/version"

	winres "github.com/tc-hib/winres"
	winres_version "github.com/tc-hib/winres/version"
)

var archNames = map[winres.Arch]string{
	winres.ArchI386:  "winres_windows_386",
	winres.ArchAMD64: "winres_windows_amd64",
	winres.ArchARM:   "winres_windows_arm",
	winres.ArchARM64: "winres_windows_arm64",
}

func createWinResVersionInfo() winres_version.Info {
	v := version.Get()

	b := v.Version.Build
	v.Version.Build = v.Version.Build.ShortCommit()
	ver := v.Version.String()

	vi := winres_version.Info{}

	vi.SetFileVersion(fmt.Sprintf("%d.%d.%d.0", v.Version.Core.Major, v.Version.Core.Minor, v.Version.Core.Patch))
	vi.SetProductVersion(fmt.Sprintf("%d.%d.%d.0", v.Version.Core.Major, v.Version.Core.Minor, v.Version.Core.Patch))

	vi.Set(winres.LCIDDefault, "CompanyName", "ariaci")
	vi.Set(winres.LCIDDefault, "FileDescription", fmt.Sprintf("kopia-provisioner %s", ver))
	vi.Set(winres.LCIDDefault, "FileVersion", ver)
	vi.Set(winres.LCIDDefault, "InternalName", "kopia-provisioner")
	vi.Set(winres.LCIDDefault, "LegalCopyright", "© ariaci. Licensed under MIT.")
	vi.Set(winres.LCIDDefault, "OriginalFilename", "kopia-provisioner.exe")
	vi.Set(winres.LCIDDefault, "PrivateBuild", b.String())
	vi.Set(winres.LCIDDefault, "ProductName", "kopia-provisioner")
	vi.Set(winres.LCIDDefault, "ProductVersion", ver)

	return vi
}

func writeWinResResources(rs winres.ResourceSet, a winres.Arch) {
	fn, ok := archNames[a]
	if !ok {
		log.Fatalf("unsupported architecture: %v", a)
	}

	out, err := os.Create(fmt.Sprintf("../../version/%s.syso", fn))
	if err != nil {
		log.Fatalf("unable to create output file: %v", err)
	}

	defer func() {
		if err := out.Close(); err != nil {
			log.Fatalf("unable to close output file: %v", err)
		}
	}()

	rs.WriteObject(out, a)
}

func main() {
	rs := winres.ResourceSet{}
	rs.SetVersionInfo(createWinResVersionInfo())

	for _, a := range []winres.Arch{
		winres.ArchI386,
		winres.ArchAMD64,
		winres.ArchARM64,
	} {
		writeWinResResources(rs, a)
	}
}
