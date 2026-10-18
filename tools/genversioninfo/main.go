package main

import (
	"log"

	genverx "github.com/ariaci/rawverx/genverx"
)

func main() {
	i, err := genverx.NewInfoWithVersionFromFile("../../version/VERSION", genverx.AdditionalInfo{
		ProductName: "kopia-provisioner",
		Author:      "ariaci",
		Copyright:   "© ariaci. Licensed under MIT.",
	})
	if err != nil {
		log.Fatalf("Error creating version info: %v", err)
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
