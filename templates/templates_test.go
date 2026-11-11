package templates

import (
	"io/fs"
	"testing"
)

func TestFS_ContainsBootstrapTemplates(t *testing.T) {
	files := []string{
		"initial/flake.nix",
		"initial/darwin.nix",
		"initial/home.nix",
		"initial/bin/bootstrap",
		"initial/bin/install_nix",
	}

	for _, name := range files {
		if _, err := fs.Stat(FS, name); err != nil {
			t.Errorf("expected embedded file %s: %v", name, err)
		}
	}
}

func TestFS_ContainsRebuildTemplates(t *testing.T) {
	files := []string{
		"files/flake.nix",
		"files/mac.nix",
		"files/linux.nix",
		"files/modules/common.nix",
	}

	for _, name := range files {
		if _, err := fs.Stat(FS, name); err != nil {
			t.Errorf("expected embedded file %s: %v", name, err)
		}
	}
}

func TestFS_ExcludesFlakeExamples(t *testing.T) {
	if _, err := fs.Stat(FS, "flakes"); err == nil {
		t.Error("expected flakes/ examples not to be embedded")
	}
}
