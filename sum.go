package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func doMD5(target string) error {
	file, err := os.Open(target)
	if err != nil {
		return errors.New("error opening file")
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return errors.New("error hashing file")
	}

	// Remove the build directory prefix from the target
	relativeTarget, err := filepath.Rel(*outDir, target)
	if err != nil {
		return errors.New("error getting relative path")
	}

	sum := hash.Sum(nil)

	// write the sum to a file
	sumFile, err := os.Create(target + ".md5.txt")
	if err != nil {
		return errors.New("error creating sum file")
	}
	defer sumFile.Close()

	content := fmt.Sprintf("%x  %s\n", sum, relativeTarget)
	_, err = sumFile.WriteString(content)
	if err != nil {
		return errors.New("error writing sum file")
	}

	return nil
}

func doSHA1(target string) error {
	file, err := os.Open(target)
	if err != nil {
		return errors.New("error opening file")
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return errors.New("error hashing file")
	}

	// Remove the build directory prefix from the target
	relativeTarget, err := filepath.Rel(*outDir, target)
	if err != nil {
		return errors.New("error getting relative path")
	}

	sum := hash.Sum(nil)

	// write the sum to a file
	sumFile, err := os.Create(target + ".sha1.txt")
	if err != nil {
		return errors.New("error creating sum file")
	}
	defer sumFile.Close()

	content := fmt.Sprintf("%x  %s\n", sum, relativeTarget)
	_, err = sumFile.WriteString(content)
	if err != nil {
		return errors.New("error writing sum file")
	}

	return nil
}

func doSHA256(target string) error {
	file, err := os.Open(target)
	if err != nil {
		return errors.New("error opening file")
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return errors.New("error hashing file")
	}

	// Remove the build directory prefix from the target
	relativeTarget, err := filepath.Rel(*outDir, target)
	if err != nil {
		return errors.New("error getting relative path")
	}

	sum := hash.Sum(nil)

	// write the sum to a file
	sumFile, err := os.Create(target + ".sha256.txt")
	if err != nil {
		return errors.New("error creating sum file")
	}
	defer sumFile.Close()

	content := fmt.Sprintf("%x  %s\n", sum, relativeTarget)
	_, err = sumFile.WriteString(content)
	if err != nil {
		return errors.New("error writing sum file")
	}

	return nil
}
