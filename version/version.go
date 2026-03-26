package version

import (
	_ "embed"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
)

//go:embed VERSION
var version SemVer

type Build RevisionInfo
type RevisionInfo struct {
	Commit string
	Dirty  bool
}

type PreRelease string
type SemVer string

type Info struct {
	Version  SemVer
	Revision RevisionInfo
	Time     string
}

func Get() Info {
	v := Info{
		Version: version,
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

func (v SemVer) component(n int) uint {
	parts := strings.SplitN(string(v), ".", 3)
	if len(parts) <= n {
		return 0
	}
	if i, err := strconv.Atoi(parts[n]); err == nil {
		return uint(i)
	}

	return 0
}

func (v SemVer) Major() uint {
	return v.Core().component(0)
}

func (v SemVer) Minor() uint {
	return v.Core().component(1)
}

func (v SemVer) Patch() uint {
	return v.Core().component(2)
}

func (v SemVer) Core() SemVer {
	parts := strings.SplitN(string(v), "-", 2)
	return SemVer(parts[0])
}

func (v SemVer) PreRelease() PreRelease {
	parts := strings.SplitN(string(v), "-", 2)
	if len(parts) < 2 {
		return ""
	}
	return PreRelease(parts[1])
}

func (p PreRelease) String() string {
	if p == "" {
		return ""
	}

	return "-" + string(p)
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

func (i RevisionInfo) Build() Build {
	return Build(i)
}

func (b Build) String() string {
	if b.Commit == "" {
		return ""
	}

	s := "+git." + b.Commit
	if b.Dirty {
		s += ".dirty"
	}

	return s
}

func (i RevisionInfo) String() string {
	if i.Commit == "" {
		return "unknown"
	}

	s := i.Commit
	if i.Dirty {
		s += ".dirty"
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
	return fmt.Sprint(v.Version, v.Revision.Build())
}

func (v Info) Print() {
	s := fmt.Sprint(v.Version, " build: ", v.Revision)
	if v.Time != "" {
		s += fmt.Sprint(" (", v.Time, ")")
	}

	fmt.Println(s)
}
