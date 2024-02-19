# filesystem <!-- omit in toc -->

- [How to use ?](#how-to-use-)
- [Features](#features)

## How to use ?

```sh
go get github.com/kilianpaquier/filesystem@latest
```

## Features

The filesystem package exposes some useful function around files, for instance a simple function as `func Exists(src) bool` to verify a file existence easily.

It also exposes `CopyFile(src, dst) error` and `CopyFileWithPerm(src, dst, perm) error` which copy a given src file to dst with either specific permissions or not.

It also exposes `CopyDir(srcdir, destdir)` to copy a full directory at another place. The destination directory will be created if it doesn't already.

The package also exposes some constants around permissions.