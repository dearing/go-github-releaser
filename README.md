# go-github-releaser

simple cli tool to cross-compile a go program, optionally sum the files, compress and cut a release with the output

> [!TIP]
> only shells out to run `go build`; hasing, zipping and github api are built in

> [!WARNING]
> CGO and complex builds not supported.

- go version go1.23.5

```
go install github.com/dearing/go-github-releaser@latest
```
---
```
export GITHUB_TOKEN=blah
```
## about
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

## cross-compose list
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
## example on windows11/amd64
```
PS D:\projects\go-github-releaser> go-github-releaser -sum-sha1 -github-owner dearing -github-repos go-github-releaser -release-commitish 1cdc3f41b8500d1585d4ca318b77f82b268b2049 -release-name "version 0.2.1" -release-tag v0.2.1 -cut-release
github.com/dearing/go-github-releaser version 0.2.0
2025/02/01 02:09:43 building darwin/amd64 build/go-github-releaser-darwin-amd64
2025/02/01 02:09:44   creating sha1sum for build/go-github-releaser-darwin-amd64
2025/02/01 02:09:44 build operation took 600.4671ms
2025/02/01 02:09:44 building darwin/arm64 build/go-github-releaser-darwin-arm64
2025/02/01 02:09:45   creating sha1sum for build/go-github-releaser-darwin-arm64
2025/02/01 02:09:45 build operation took 595.0914ms
2025/02/01 02:09:45 building linux/amd64 build/go-github-releaser-linux-amd64
2025/02/01 02:09:45   creating sha1sum for build/go-github-releaser-linux-amd64
2025/02/01 02:09:45 build operation took 597.4899ms
2025/02/01 02:09:45 building linux/arm64 build/go-github-releaser-linux-arm64
2025/02/01 02:09:46   creating sha1sum for build/go-github-releaser-linux-arm64
2025/02/01 02:09:46 build operation took 582.1851ms
2025/02/01 02:09:46 building windows/amd64 build/go-github-releaser-windows-amd64.exe
2025/02/01 02:09:46   creating sha1sum for build/go-github-releaser-windows-amd64.exe
2025/02/01 02:09:46 build operation took 604.1977ms
2025/02/01 02:09:46 building windows/arm64 build/go-github-releaser-windows-arm64.exe
2025/02/01 02:09:47   creating sha1sum for build/go-github-releaser-windows-arm64.exe
2025/02/01 02:09:47 build operation took 590.7436ms
2025/02/01 02:09:47 cutting release v0.2.1
1cdc3f41b8500d1585d4ca318b77f82b268b2049
2025/02/01 02:09:48 created release https://github.com/dearing/go-github-releaser/releases/tag/untagged-5c87b091733e1206efcb
2025/02/01 02:09:48 uploading go-github-releaser-darwin-amd64
2025/02/01 02:09:50 uploaded go-github-releaser-darwin-amd64 to https://github.com/dearing/go-github-releaser/releases/download/untagged-5c87b091733e1206efcb/go-github-releaser-darwin-amd64
2025/02/01 02:09:50 uploading go-github-releaser-darwin-amd64.sha1.txt
2025/02/01 02:09:51 uploaded go-github-releaser-darwin-amd64.sha1.txt to https://github.com/dearing/go-github-releaser/releases/download/untagged-5c87b091733e1206efcb/go-github-releaser-darwin-amd64.sha1.txt
2025/02/01 02:09:51 uploading go-github-releaser-darwin-arm64
2025/02/01 02:09:53 uploaded go-github-releaser-darwin-arm64 to https://github.com/dearing/go-github-releaser/releases/download/untagged-5c87b091733e1206efcb/go-github-releaser-darwin-arm64
2025/02/01 02:09:53 uploading go-github-releaser-darwin-arm64.sha1.txt
2025/02/01 02:09:53 uploaded go-github-releaser-darwin-arm64.sha1.txt to https://github.com/dearing/go-github-releaser/releases/download/untagged-5c87b091733e1206efcb/go-github-releaser-darwin-arm64.sha1.txt
2025/02/01 02:09:53 uploading go-github-releaser-linux-amd64
2025/02/01 02:09:55 uploaded go-github-releaser-linux-amd64 to https://github.com/dearing/go-github-releaser/releases/download/untagged-5c87b091733e1206efcb/go-github-releaser-linux-amd64
2025/02/01 02:09:55 uploading go-github-releaser-linux-amd64.sha1.txt
2025/02/01 02:09:56 uploaded go-github-releaser-linux-amd64.sha1.txt to https://github.com/dearing/go-github-releaser/releases/download/untagged-5c87b091733e1206efcb/go-github-releaser-linux-amd64.sha1.txt
2025/02/01 02:09:56 uploading go-github-releaser-linux-arm64
2025/02/01 02:09:58 uploaded go-github-releaser-linux-arm64 to https://github.com/dearing/go-github-releaser/releases/download/untagged-5c87b091733e1206efcb/go-github-releaser-linux-arm64
2025/02/01 02:09:58 uploading go-github-releaser-linux-arm64.sha1.txt
2025/02/01 02:09:58 uploaded go-github-releaser-linux-arm64.sha1.txt to https://github.com/dearing/go-github-releaser/releases/download/untagged-5c87b091733e1206efcb/go-github-releaser-linux-arm64.sha1.txt
2025/02/01 02:09:58 uploading go-github-releaser-windows-amd64.exe
2025/02/01 02:10:01 uploaded go-github-releaser-windows-amd64.exe to https://github.com/dearing/go-github-releaser/releases/download/untagged-5c87b091733e1206efcb/go-github-releaser-windows-amd64.exe
2025/02/01 02:10:01 uploading go-github-releaser-windows-amd64.exe.sha1.txt
2025/02/01 02:10:01 uploaded go-github-releaser-windows-amd64.exe.sha1.txt to https://github.com/dearing/go-github-releaser/releases/download/untagged-5c87b091733e1206efcb/go-github-releaser-windows-amd64.exe.sha1.txt
2025/02/01 02:10:01 uploading go-github-releaser-windows-arm64.exe
2025/02/01 02:10:03 uploaded go-github-releaser-windows-arm64.exe to https://github.com/dearing/go-github-releaser/releases/download/untagged-5c87b091733e1206efcb/go-github-releaser-windows-arm64.exe
2025/02/01 02:10:03 uploading go-github-releaser-windows-arm64.exe.sha1.txt
2025/02/01 02:10:04 uploaded go-github-releaser-windows-arm64.exe.sha1.txt to https://github.com/dearing/go-github-releaser/releases/download/untagged-5c87b091733e1206efcb/go-github-releaser-windows-arm64.exe.sha1.txt
2025/02/01 02:10:04 cut operation took 16.5850329s
PS D:\projects\go-github-releaser> cat .\build\*.txt
e0d1dc56d3b8e28974862891f23ae0528d7a5332  go-github-releaser-darwin-amd64
301f42f224bbd05f49d417ed927130dc97562843  go-github-releaser-darwin-arm64
26985fd24e7d1885dda8fd94b47cdac883b6b134  go-github-releaser-linux-amd64
796b938b86f2192a536c4f829af1ae35d90f16b1  go-github-releaser-linux-arm64
9b85871cd2be180da79807925c321aa99b9fc55d  go-github-releaser-windows-amd64.exe
a86722243b67aef3b43d351f876c8e2aef9167b0  go-github-releaser-windows-arm64.exe
```
