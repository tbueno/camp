package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type TemplDir interface {
	ReadFile(name string) ([]byte, error)
	ReadDir(name string) ([]fs.DirEntry, error)
}

func CopyFile(src, dst string) error {
	bytesRead, err := os.ReadFile(src)

	if err != nil {
		return err
	}

	return SaveFile(bytesRead, dst)
}

func SaveFile(content []byte, file string) error {
	return os.WriteFile(file, content, 0644)
}

// ReplaceInContent replaces all occurrences of old with new in the content
func ReplaceInContent(content []byte, old, new string) []byte {
	c := string(content)
	c = strings.ReplaceAll(c, old, new)
	return []byte(c)
}

func CopyTemplDir(templDir TemplDir, src, dst string) error {
	// Create the destination directory if it does not exist
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dst, err)
	}

	// Iterate over the contents of the source directory
	entries, err := templDir.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", src, err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy the directory
			if err := CopyTemplDir(templDir, srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy the file
			if err := CopyConfFile(templDir, srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyConfFile(templDir TemplDir, srcPath, dstPath string) error {
	src, err := templDir.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", srcPath, err)
	}
	return SaveFile(src, dstPath)

}

// BackupFile backs up a file by renaming it to filename.before-nix
func BackupFile(filepath string) error {
	if _, err := os.Stat(filepath); err == nil {
		if err = RunCommand("sudo", "mv", filepath, filepath+".before-nix"); err != nil {
			return err
		}
	}
	return nil
}

// HostName returns the hostname of the user's machine by using the native calls to the OS
func HostName() string {
	var hostname string
	var err error
	if runtime.GOOS == "darwin" {
		hostname, err = CommandReturn("scutil", "--get", "LocalHostName")
	} else {
		hostname, err = CommandReturn("hostname")
	}
	if err != nil {
		panic("Could not get hostname: " + err.Error())
	}

	return strings.TrimSpace(hostname)
}
