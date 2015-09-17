package test_helpers

import (
	"os/exec"
	"path/filepath"
)

func RunSLStemcells(rootPath string, args []string) ([]byte, error) {
	slStemcellsPath := filepath.Join(rootPath, "out/sl_stemcells")

	cmd := exec.Command(slStemcellsPath, args...)
	output, err := cmd.Output()

	if err != nil {
		return []byte{}, err
	}

	return output, nil
}
