package common

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CreateFile(filePath string, data []byte) error {
	err := ioutil.WriteFile(filePath, data, 6000)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create file '%s', got error '%s'", filePath, err.Error()))
	}

	return nil
}

func CreateTarball(tarballFilePath string, filePaths []string) error {
	file, err := os.Create(tarballFilePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create tarball file '%s', got error '%s'", tarballFilePath, err.Error()))
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, filePath := range filePaths {
		err := addFileToTarWriter(filePath, tarWriter)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not add file '%s', to tarball, got error '%s'", filePath, err.Error()))
		}
	}

	return nil
}

func ReadJsonTestFixtures(packageName, fileName string) ([]byte, error) {
	wd, _ := os.Getwd()
	return ioutil.ReadFile(filepath.Join(wd, "../..", "test_fixtures", packageName, fileName))
}

// Private methods

func addFileToTarWriter(filePath string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open file '%s', got error '%s'", filePath, err.Error()))
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not get stat for file '%s', got error '%s'", filePath, err.Error()))
	}

	header := &tar.Header{
		Name:    filePath,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not write header for file '%s', got error '%s'", filePath, err.Error()))
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not copy the file '%s' data to the tarball, got error '%s'", filePath, err.Error()))
	}

	return nil
}
