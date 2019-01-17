package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func getModTime(file string) int64 {
	info, err := os.Stat(file)
	if err != nil {
		return time.Now().Unix()
	}
	return info.ModTime().Unix()
}

func getOldestFileModTimestamp(dir string, recursive bool) int64 {
	if recursive {
		return getOldestAgeInDirRecursively(dir)
	} else {
		return getOldestAgeInDir(dir)
	}
}

func getOldestAgeInDirRecursively(dir string) int64 {
	var oldestTs int64 = time.Now().Unix()
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			ts := getModTime(path)
			if ts < oldestTs {
				oldestTs = ts
			}
		}
		return nil
	})
	return oldestTs
}

func getOldestAgeInDir(dir string) int64 {
	var files, _ = ioutil.ReadDir(dir)
	var oldestTs int64 = time.Now().Unix()
	for _, file := range files {
		if !file.IsDir() {
			ts := getModTime(dir + string(os.PathSeparator) + file.Name())
			if ts < oldestTs {
				oldestTs = ts
			}
		}
	}
	return oldestTs
}
