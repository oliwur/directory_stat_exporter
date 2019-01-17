package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func getFileCount(dir string, recursive bool) int {
	if recursive {
		return getFileCountInDirRecursively(dir)
	} else {
		return getFileCountInDir(dir)
	}
}

func getFileCountInDirRecursively(dir string) int {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return -1
	}
	count := 0
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			count++
		}
		return nil
	})
	return count
}

func getFileCountInDir(dir string) int {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return -1
	}
	count := 0
	for _, f := range files {
		if !f.IsDir() {
			count++
		}
	}
	return count
}
