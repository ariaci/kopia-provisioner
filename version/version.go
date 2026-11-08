package version

//go:generate go run -C ../tools/genversioninfo -buildvcs .

import (
	_ "embed"

	rawverx "github.com/ariaci/rawverx"
)

//go:embed VERSION
var version string

func Get() rawverx.Info {
	i, _ := rawverx.NewInfo(version)
	return i
}
