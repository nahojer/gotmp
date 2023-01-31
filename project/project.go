package project

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Root represents the strategy to use to find the root directory of a project.
type Root int

// Supported project root strategies.
const (
	// Uses the parent directory of .git/
	Git Root = iota
	// Uses the parent directory of the path pointed to by the GOMOD environment
	// variable.
	GoModule
)

func (r Root) GetDir() (string, error) {
	switch r {
	case Git:
		return gitRootDir()
	case GoModule:
		return goModuleDir()
	default:
		return "", fmt.Errorf("invalid Root %d", r)
	}
}

func gitRootDir() (string, error) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Stdout = writer

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to find git root path: %w", err)
	}

	return strings.TrimSpace(buf.String()), nil
}

func goModuleDir() (string, error) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	cmd := exec.Command("go", "env", "--json")
	cmd.Stdout = writer

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to find go mod root path: %w", err)
	}

	var goEnv struct {
		GOMOD string
	}
	if err := json.Unmarshal(buf.Bytes(), &goEnv); err != nil {
		return "", fmt.Errorf("failed to get GOMOD environment variable: %w", err)
	}

	return filepath.Dir(goEnv.GOMOD), nil
}
