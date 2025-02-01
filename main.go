package main

import (
	"archive/zip"
	"bufio"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v68/github"
)

var version = "0.2.1"

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

// skip because we already built the binaries
var skipBuild = flag.Bool("skip-build", false, "skip the build step")

// cut a github release
var cutRelease = flag.Bool("cut-release", false, "cut a github release")

// github information
var githubOwner = flag.String("github-owner", "", "github owner")
var githubRepos = flag.String("github-repos", "", "github repos")

// release information
var releaseTag = flag.String("release-tag", "", "github release tag")
var releaseName = flag.String("release-name", "", "github release name")
var releaseMessage = flag.String("release-message", "", "github release message")
var releaseDraft = flag.Bool("release-draft", true, "github release draft")
var releasePrerelease = flag.Bool("release-prerelease", false, "github release prerelease")
var releaseCommitish = flag.String("release-commitish", "", "github release commitish") // this is a dumb name from github, consider renaming

func main() {

	flag.Parse()

	fmt.Printf("github.com/dearing/go-github-releaser version %s\n", version)

	if !*skipBuild {

		file, err := os.Open(*csvFile)
		if err != nil {
			fmt.Printf("error opening file: %v\n", err)
			return
		}
		defer file.Close()

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

	if *cutRelease {
		start := time.Now()
		log.Printf("cutting release %s\n%s", *releaseTag, *releaseCommitish)
		err := cut(*githubOwner, *githubRepos)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		log.Printf("cut operation took %s\n", time.Since(start))
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

func doBuild(goOS, goARCH, target string) error {

	cmd := exec.Command("go", "build", "-o", target, *srcDir)

	cmd.Env = append(os.Environ(),
		"GOOS="+goOS,
		"GOARCH="+goARCH,
	)

	origStdout := os.Stdout
	origStderr := os.Stderr

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

	relativeTarget, err := filepath.Rel(*outDir, target)
	if err != nil {
		return errors.New("error getting relative path")
	}

	sum := hash.Sum(nil)

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

	relativeTarget, err := filepath.Rel(*outDir, target)
	if err != nil {
		return errors.New("error getting relative path")
	}

	sum := hash.Sum(nil)

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

	relativeTarget, err := filepath.Rel(*outDir, target)
	if err != nil {
		return errors.New("error getting relative path")
	}

	sum := hash.Sum(nil)

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

func cut(owner, repos string) error {

	token := os.Getenv("GITHUB_TOKEN")
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
