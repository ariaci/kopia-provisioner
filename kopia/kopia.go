package kopia

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
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

// Discard ignores all incoming lines.
func Discard(string) {}

func Run(linehandler func(string), repoConfigPath string, args ...string) error {
	if linehandler == nil {
		return fmt.Errorf("linehandler must not be nil")
	}

	base := []string{}
	if len(repoConfigPath) > 0 {
		base = append(base, "--config-file", repoConfigPath)
	}

	cmd := exec.Command("kopia", append(base, args...)...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start kopia: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		linehandler(scanner.Text())
	}

	errBuf := new(bytes.Buffer)
	io.Copy(errBuf, stderr)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("kopia command failed: %w\n%s", err, errBuf.String())
	}

	return nil
}
