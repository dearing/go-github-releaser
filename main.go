package main

import (
	"archive/zip"
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var version = "0.0.0"

// csv file build-matrix
var csvFile = flag.String("csv-file", "go-github-releaser.csv", "csv file with build information")

// go source directory and binary output directory
var srcDir = flag.String("src-dir", ".", "go source directory")
var outDir = flag.String("out-dir", "build", "binary output directory")

// produce sum txt files for each binary
var sumMD5 = flag.Bool("sum-md5", false, "create md5 sum file")
var sumSHA1 = flag.Bool("sum-sha1", false, "create sha1 sum file")
var sumSHA256 = flag.Bool("sum-sha256", false, "create sha256 sum file")

// produce a zip file for each binary
var zipFile = flag.Bool("zip", false, "create zip file")

func main() {

	flag.Parse()

	fmt.Printf("github.com/dearing/go-github-releaser version %s\n", version)

	// look for a go-github-releaser.csv file
	file, err := os.Open(*csvFile)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// read the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			fmt.Printf("invalid line: %s\n", line)
			continue
		}
		goOS := parts[0]
		goARCH := parts[1]
		name := parts[2]

		target := fmt.Sprintf("%s/%s", *outDir, name)

		do(goOS, goARCH, target)
	}

}

func do(goOS, goARCH, target string) {

	log.Printf("building %s/%s %s\n", goOS, goARCH, target)
	err := doBuild(goOS, goARCH, target)
	if err != nil {
		fmt.Printf("error building: %s: %v\n", target, err)
		return
	}

	if *sumMD5 {
		log.Printf("creating md5sum for %s\n", target)
		err = doMD5(target)
		if err != nil {
			fmt.Printf("error creating md5: %s: %v\n", target, err)
		}
	}

	if *sumSHA1 {
		log.Printf("creating sha1sum for %s\n", target)
		err = doSHA1(target)
		if err != nil {
			fmt.Printf("error creating sha1sum: %s: %v\n", target, err)
		}
	}

	if *sumSHA256 {
		log.Printf("creating sha256sum for %s\n", target)
		err = doSHA256(target)
		if err != nil {
			fmt.Printf("error creating sha256sum: %s: %v\n", target, err)
		}
	}

	if *zipFile {
		log.Printf("creating zip for %s\n", target)
		err = doZip(target)
		if err != nil {
			fmt.Printf("error creating zip: %s: %v\n", target, err)
		}
	}
}

func doBuild(goOS, goARCH, target string) error {

	cmd := exec.Command("go", "build", "-o", target, *srcDir)

	// set environment variables
	cmd.Env = append(os.Environ(),
		"GOOS="+goOS,
		"GOARCH="+goARCH,
	)

	// Save the original stdout and stderr
	origStdout := os.Stdout
	origStderr := os.Stderr

	// Restore stdout and stderr after completion
	defer func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
	}()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

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
