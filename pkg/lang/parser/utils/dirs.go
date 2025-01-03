package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/pouya-eghbali/posh/pkg/constants"
)

func CreateTempDir() (string, error) {
	tempDir, err := os.MkdirTemp("", "posh-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %v", err)
	}

	return tempDir, nil
}

func CompileTempDir(tempDir string, output string) error {
	// Run go mod init
	cmd := exec.Command("go", "mod", "init", "main")
	cmd.Dir = tempDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run go mod init: %v, output: %s", err, string(output))
	}

	// Run go get github.com/pouya-eghbali/posh@version
	cmd = exec.Command("go", "get", fmt.Sprintf("github.com/pouya-eghbali/posh@v%s", constants.Version))
	cmd.Dir = tempDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run go get: %v, output: %s", err, string(output))
	}

	// Run go mod tidy
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = tempDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %v, output: %s", err, string(output))
	}

	// Run go build
	cmd = exec.Command("go", "build", "-ldflags", "-s -w", "-o", "main")
	cmd.Dir = tempDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to compile temp dir: %v, output: %s", err, string(output))
	}

	// Move the compiled binary to the output path
	err := os.Rename(path.Join(tempDir, "main"), output)
	if err != nil {
		return fmt.Errorf("failed to move compiled binary: %v", err)
	}

	// Remove the temp dir
	err = os.RemoveAll(tempDir)
	if err != nil {
		return fmt.Errorf("failed to remove temp dir: %v", err)
	}

	return nil
}
