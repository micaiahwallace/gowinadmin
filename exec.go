package gowinadmin

import (
	"bytes"
	"os/exec"
)

// Take the cli arguments and execute command returning stdout and stderr strings
func executeCommand(command string, args []string) (string, string, error) {
	cmd := exec.Command(command, args...)
	var outData, errData bytes.Buffer
	cmd.Stdout = &outData
	cmd.Stderr = &errData
	err := cmd.Run()
	if err != nil {
		return "", "", err
	}
	return outData.String(), errData.String(), nil
}
