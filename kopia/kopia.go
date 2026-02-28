package kopia

import (
	"fmt"
	"os/exec"
	"strings"
)

type KopiaAction uint8

const (
	KopiaActionAdd KopiaAction = iota
	KopiaActionRemove
	KopiaActionUpdate
)

func (a KopiaAction) String() string {
	switch a {
	case KopiaActionAdd:
		return "add"
	case KopiaActionUpdate:
		return "set"
	case KopiaActionRemove:
		return "delete"
	}
	return ""
}

func Run(repoConfigPath string, args ...string) ([]string, error) {
	base := []string{}
	if len(repoConfigPath) > 0 {
		base = append(base, "--config-file", repoConfigPath)
	}

	cmd := exec.Command("kopia", append(base, args...)...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("kopia command failed: %w\n%s", err, out)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines, nil
}
