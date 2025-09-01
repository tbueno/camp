package system

import "testing"

func TestGetSystemInfo(t *testing.T) {
	sysInfo, err := GetSystemInfo()
	if err != nil {
		t.Errorf("GetSystemInfo() failed: %v", err)
	}

	if sysInfo == nil {
		t.Fatal("Expected system info to not be nil")
	}

	if sysInfo.Architecture == "" {
		t.Error("Expected architecture to not be empty")
	}

	if sysInfo.OS == "" {
		t.Error("Expected operating system to not be empty")
	}

	validArchs := []string{"x86_64", "arm64", "aarch64", "i386", "armv7l"}
	isValidArch := false
	for _, validArch := range validArchs {
		if sysInfo.Architecture == validArch {
			isValidArch = true
			break
		}
	}

	if !isValidArch {
		t.Logf("Got architecture: %s (may be valid but not in our test list)", sysInfo.Architecture)
	}

	validOS := []string{"darwin", "linux"}
	isValidOS := false
	for _, validOSName := range validOS {
		if sysInfo.OS == validOSName {
			isValidOS = true
			break
		}
	}

	if !isValidOS {
		t.Logf("Got OS: %s (may be valid but not in our test list)", sysInfo.OS)
	}
}

func TestGetMachineArchitecture(t *testing.T) {
	arch, err := getMachineArchitecture()
	if err != nil {
		t.Errorf("getMachineArchitecture() failed: %v", err)
	}

	if arch == "" {
		t.Error("Expected architecture to not be empty")
	}
}

func TestGetOperatingSystem(t *testing.T) {
	os, err := getOperatingSystem()
	if err != nil {
		t.Errorf("getOperatingSystem() failed: %v", err)
	}

	if os == "" {
		t.Error("Expected operating system to not be empty")
	}
}
