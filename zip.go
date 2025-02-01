package main

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

func doZip(target string) error {
	zipFile, err := os.Create(target + ".zip")
	if err != nil {
		return errors.New("error creating zip file")
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileToZip, err := os.Open(target)
	if err != nil {
		return errors.New("error opening file to zip")
	}
	defer fileToZip.Close()

	w, err := zipWriter.Create(filepath.Base(target))
	if err != nil {
		return errors.New("error creating zip entry")
	}

	if _, err := io.Copy(w, fileToZip); err != nil {
		return errors.New("error writing to zip")
	}

	return nil
}
