package version

import (
	_ "embed"
	"fmt"
	"runtime/debug"
)

//go:embed VERSION
var Version string

type RevisionInfo struct {
	Commit string
	Dirty  bool
}

type Info struct {
	Version  string
	Revision RevisionInfo
	Time     string
}

func Get() Info {
	v := Info{
		Version: Version,
	}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return v
	}

	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			v.Revision.Commit = s.Value
		case "vcs.modified":
			v.Revision.Dirty = (s.Value == "true")
		case "vcs.time":
			v.Time = s.Value
		}
	}

	return v
}

func (i RevisionInfo) Short() RevisionInfo {
	c := i.Commit
	if len(c) > 7 {
		c = c[:7]
	}

	return RevisionInfo{
		Commit: c,
		Dirty:  i.Dirty,
	}
}

func (i RevisionInfo) String() string {
	if i.Commit == "" {
		return "unknown"
	}

	s := i.Commit
	if i.Dirty {
		s += "-dirty"
	}

	return s
}

func (v Info) Short() Info {
	return Info{
		Version:  v.Version,
		Revision: v.Revision.Short(),
		Time:     v.Time,
	}
}

func (v Info) String() string {
	s := fmt.Sprintf("%s build: %s", v.Version, v.Revision)
	if v.Time != "" {
		s += fmt.Sprintf(" (%s)", v.Time)
	}

	return s
}

func (v Info) Print() {
	fmt.Println(v)
}
