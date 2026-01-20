package project

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const trackingDir = ".camp"
const trackingFile = "env_vars"

// GetTrackingFilePath returns the path to the env tracking file for the given project
func GetTrackingFilePath(projectPath string) string {
	return filepath.Join(projectPath, trackingDir, trackingFile)
}

// LoadTrackedVars loads the list of previously exported variable names from the tracking file
// Returns an empty slice if the file doesn't exist
func LoadTrackedVars(projectPath string) ([]string, error) {
	trackingPath := GetTrackingFilePath(projectPath)

	file, err := os.Open(trackingPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var vars []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			vars = append(vars, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return vars, nil
}

// SaveTrackedVars saves the list of exported variable names to the tracking file
func SaveTrackedVars(projectPath string, vars []string) error {
	trackingPath := GetTrackingFilePath(projectPath)

	// Create .camp directory if it doesn't exist
	trackingDirPath := filepath.Dir(trackingPath)
	if err := os.MkdirAll(trackingDirPath, 0755); err != nil {
		return err
	}

	file, err := os.Create(trackingPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, v := range vars {
		if _, err := file.WriteString(v + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// GetRemovedVars returns the variable names that were previously tracked but are no longer in currentVars
func GetRemovedVars(trackedVars []string, currentVars map[string]string) []string {
	var removed []string
	for _, v := range trackedVars {
		if _, exists := currentVars[v]; !exists {
			removed = append(removed, v)
		}
	}
	return removed
}
