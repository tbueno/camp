package system

import (
	"os/exec"
	"strings"
)

// GetSystemInfo retrieves system information
func GetSystemInfo() (*System, error) {
	arch, err := getMachineArchitecture()
	if err != nil {
		return nil, err
	}

	os, err := getOperatingSystem()
	if err != nil {
		return nil, err
	}

	return &System{
		OS:           os,
		Architecture: arch,
	}, nil
}

// getMachineArchitecture gets the machine architecture
func getMachineArchitecture() (string, error) {
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getOperatingSystem gets the operating system from the hosting machine
func getOperatingSystem() (string, error) {
	cmd := exec.Command("uname", "-s")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	osName := strings.TrimSpace(strings.ToLower(string(output)))
	return osName, nil
}
