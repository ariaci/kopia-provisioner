package main

import (
	"log"
	"rawverx-app/version"

	genverx "github.com/ariaci/rawverx/genverx"
)

func main() {
	i := genverx.Info{
		Version:     version.Get().Version,
		ProductName: "kopia-provisioner",
		Author:      "ariaci",
		Copyright:   "© ariaci. Licensed under MIT.",
	}

	for _, a := range []genverx.Arch{
		genverx.ArchI386,
		genverx.ArchAMD64,
		genverx.ArchARM64,
	} {
		if err := genverx.Generate(i, a, "../../version/winres"); err != nil {
			log.Fatalf("Error generating version info for architecture %v: %v", a, err)
		}
	}
}
