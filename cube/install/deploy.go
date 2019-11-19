package install

import (
	"fmt"
	"os"
	"path/filepath"
)

type (
	ExistFile func(string) bool
	MkdirAll  func(string, os.FileMode) error
	Rename    func(string, string) error

	DeployFileParam struct {
		ExistFile ExistFile
		MkdirAll  MkdirAll
		Rename    Rename
		Source    string
		Dest      string
		DirPerm   os.FileMode
	}
)

func DeployFile(params DeployFileParam) error {
	if params.ExistFile(params.Dest) {
		return nil
	}
	dir := filepath.Dir(params.Dest)
	if !params.ExistFile(dir) {
		// create parent directory
		if err := params.MkdirAll(dir, params.DirPerm); err != nil {
			return fmt.Errorf("failed to create a directory: %s : %w", dir, err)
		}
	}
	// move
	if err := os.Rename(params.Source, params.Dest); err != nil {
		return fmt.Errorf("failed to copy a file %s to %s: %w", params.Source, params.Dest, err)
	}
	return nil
}
