package utils

import (
	"os"
	"os/exec"
	"strings"
)

func RunCommand(command string, workDir string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workDir
	return cmd.Run()
}

func RunCommandSilent(command string, workDir string, arg ...string) (string, error) {
	cmd := exec.Command(command, arg...)
	b := strings.Builder{}
	cmd.Stdout = &b
	cmd.Stderr = &b
	cmd.Dir = workDir
	err := cmd.Run()
	return b.String(), err
}
