package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var version = "0.1.0"

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

func usage() {
	fmt.Println("go-github-releaser is a tool to build go binaries for multiple platforms and create checksums and zip files")
	fmt.Fprintf(os.Stderr, "usage: go-github-releaser [flags]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {

	flag.Usage = usage
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

		start := time.Now()
		err := do(goOS, goARCH, target)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		log.Printf("  operation took %s\n", time.Since(start))
	}

}

func do(goOS, goARCH, target string) error {

	log.Printf("building %s/%s %s\n", goOS, goARCH, target)
	err := doBuild(goOS, goARCH, target)
	if err != nil {
		return fmt.Errorf("building: %s: %v", target, err)
	}

	if *sumMD5 {
		log.Printf("  creating md5sum for %s\n", target)
		err = doMD5(target)
		if err != nil {
			return fmt.Errorf("creating md5: %s: %v", target, err)
		}
	}

	if *sumSHA1 {
		log.Printf("  creating sha1sum for %s\n", target)
		err = doSHA1(target)
		if err != nil {
			return fmt.Errorf("creating sha1: %s: %v", target, err)
		}
	}

	if *sumSHA256 {
		log.Printf("  creating sha256sum for %s\n", target)
		err = doSHA256(target)
		if err != nil {
			return fmt.Errorf("creating sha256: %s: %v", target, err)
		}
	}

	if *zipFile {
		log.Printf("  creating zip for %s\n", target)
		err = doZip(target)
		if err != nil {
			return fmt.Errorf("creating zip: %s: %v", target, err)
		}
	}
	return nil
}
