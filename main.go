package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v68/github"
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

var cutRelease = flag.Bool("cut-release", false, "cut a github release")

var githubOwner = flag.String("github-owner", "", "github owner")
var githubRepos = flag.String("github-repos", "", "github repos")

var releaseTag = flag.String("release-tag", "", "github release tag")
var releaseName = flag.String("release-name", "", "github release name")
var releaseMessage = flag.String("release-message", "", "github release message")
var releaseDraft = flag.Bool("release-draft", true, "github release draft")
var releasePrerelease = flag.Bool("release-prerelease", false, "github release prerelease")
var releaseCommitish = flag.String("release-commitish", "", "github release commitish")

var skipBuild = flag.Bool("skip-build", false, "skip the build step")

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

	if !*skipBuild {

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

			log.Printf("build operation took %s\n", time.Since(start))
		}
	}

	start := time.Now()
	log.Printf("cutting release %s\n%s", *releaseTag, *releaseCommitish)
	if *cutRelease {
		err := cut(*githubOwner, *githubRepos)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		log.Printf("cut operation took %s\n", time.Since(start))
	}

}

func cut(owner, repos string) error {

	token := os.Getenv("GITHUB_TOKEN") // Best practice: store token in environment variable
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable not set")
	}

	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(token)
	release, resp, err := client.Repositories.CreateRelease(ctx, owner, repos, &github.RepositoryRelease{
		TagName:         releaseTag,
		TargetCommitish: releaseCommitish,
		Name:            releaseName,
		Body:            releaseMessage,
		Draft:           releaseDraft,
		Prerelease:      releasePrerelease,
	})

	if err != nil {
		return fmt.Errorf("github creating release: %v", err)
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("github creating release: %v", resp.Status)
	}

	log.Printf("created release %s\n", *release.HTMLURL)

	// upload each file in the outDir
	f := os.DirFS(*outDir)
	err = fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		file, err := os.Open(filepath.Join(*outDir, path))
		if err != nil {
			return fmt.Errorf("opening file: %v", err)
		}
		defer file.Close()

		log.Printf("uploading %s\n", path)

		asset, resp, err := client.Repositories.UploadReleaseAsset(ctx, owner, repos, *release.ID, &github.UploadOptions{
			Name: path,
		}, file)

		if err != nil {
			return fmt.Errorf("github uploading asset: %v", err)
		}

		if resp.StatusCode != 201 {
			return fmt.Errorf("github uploading asset: %v", resp.Status)
		}

		log.Printf("uploaded %s to %s\n", path, *asset.BrowserDownloadURL)
		return nil
	})

	if err != nil {
		return fmt.Errorf("walking directory: %v", err)
	}

	return nil
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
