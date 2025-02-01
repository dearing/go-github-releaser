# go-github-releaser

> [!INFORMATION]
> only shells out to run `go build`, sums, zipping and github are all internal to the binary

Simple cli tool to cross-compile a go program, optionally sum the files, compress and cut a release with the output.

> [!WARNING]
> CGO and complex builds not supported.

- go version go1.23.5

---

```
export GITHUB_TOKEN=blah
```
---
```
$ go-github-releaser -h
Usage of go-github-releaser:
  -csv-file string
        csv file with build information (default "go-github-releaser.csv")
  -cut-release
        cut a github release
  -github-owner string
        github owner
  -github-repos string
        github repos
  -out-dir string
        binary output directory (default "build")
  -release-commitish string
        github release commitish
  -release-draft
        github release draft (default true)
  -release-message string
        github release message
  -release-name string
        github release name
  -release-prerelease
        github release prerelease
  -release-tag string
        github release tag
  -skip-build
        skip the build step
  -src-dir string
        go source directory (default ".")
  -sum-md5
        create md5 sum file
  -sum-sha1
        create sha1 sum file
  -sum-sha256
        create sha256 sum file
  -zip
        create zip file
```

---

```
$ go tool dist list
aix/ppc64
android/386
android/amd64
android/arm
android/arm64
darwin/amd64
darwin/arm64
dragonfly/amd64
freebsd/386
freebsd/amd64
freebsd/arm
freebsd/arm64
freebsd/riscv64
illumos/amd64
ios/amd64
ios/arm64
js/wasm
linux/386
linux/amd64
linux/arm
linux/arm64
linux/loong64
linux/mips
linux/mips64
linux/mips64le
linux/mipsle
linux/ppc64
linux/ppc64le
linux/riscv64
linux/s390x
netbsd/386
netbsd/amd64
netbsd/arm
netbsd/arm64
openbsd/386
openbsd/amd64
openbsd/arm
openbsd/arm64
openbsd/ppc64
openbsd/riscv64
plan9/386
plan9/amd64
plan9/arm
solaris/amd64
wasip1/wasm
windows/386
windows/amd64
windows/arm
windows/arm64

```
