package version

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
