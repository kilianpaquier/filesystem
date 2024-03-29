<!-- This file is safe to edit. Once it exists it will not be overwritten. -->

# filesystem <!-- omit in toc -->

<p align="center">
  <img alt="GitHub Actions" src="https://img.shields.io/github/actions/workflow/status/kilianpaquier/filesystem/integration.yml?branch=main&style=for-the-badge">
  <img alt="GitHub Release" src="https://img.shields.io/github/v/release/kilianpaquier/filesystem?include_prereleases&sort=semver&style=for-the-badge">
  <img alt="GitHub Issues" src="https://img.shields.io/github/issues-raw/kilianpaquier/filesystem?style=for-the-badge">
  <img alt="GitHub License" src="https://img.shields.io/github/license/kilianpaquier/filesystem?style=for-the-badge">
  <img alt="Coverage" src="https://img.shields.io/codecov/c/github/kilianpaquier/filesystem/main?style=for-the-badge">
  <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/kilianpaquier/filesystem/main?style=for-the-badge&label=Go+Version">
  <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/kilianpaquier/filesystem?style=for-the-badge">
</p>

---

- [How to use ?](#how-to-use-)
- [Features](#features)

## How to use ?

```sh
go get -u github.com/kilianpaquier/filesystem@latest
```

## Features

The filesystem package exposes some useful function around files, for instance a simple function as `func Exists(src) bool` to verify a file existence easily.

It also exposes `CopyFile(src, dst) error` and `CopyFileWithPerm(src, dst, perm) error` which copy a given src file to dst with either specific permissions or not.

It also exposes `CopyDir(srcdir, destdir)` to copy a full directory at another place. The destination directory will be created if it doesn't already.

The package also exposes some constants around permissions.