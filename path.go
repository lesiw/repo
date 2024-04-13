package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

func MkdirAll(path string, perm os.FileMode) (paths []string, err error) {
	var dir fs.FileInfo
	if dir, err = os.Stat(path); err == nil {
		if dir.IsDir() {
			return
		}
		err = &os.PathError{Op: "mkdir", Path: path, Err: syscall.ENOTDIR}
		return
	}

	i := len(path) - 1
	for i >= 0 && os.IsPathSeparator(path[i]) {
		i--
	}
	for i >= 0 && !os.IsPathSeparator(path[i]) {
		i--
	}
	if i < 0 {
		i = 0
	}

	var paths1 []string
	if parent := path[:i]; len(parent) > len(filepath.VolumeName(path)) {
		if paths1, err = MkdirAll(parent, perm); err != nil {
			return
		}
		paths = append(paths, paths1...)
	}

	err = os.Mkdir(path, perm)
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			err = nil
			return
		}
		return
	}
	paths = append(paths, path)
	return
}
