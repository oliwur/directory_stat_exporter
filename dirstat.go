package main

import (
	"fmt"
	"github.com/codestoke/directory_stat_exporter/cfg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// globals
var config cfg.Config

func main() {
	config = cfg.GetConfig()

	http.HandleFunc("/metrics", handleMetrics)
	if err := http.ListenAndServe(":9999", nil); err != nil {
		panic(err)
	}
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	for _, dir := range config.Directories {
		if dir.Recursive {
			w.Write([]byte(getDirMetric("dirstat", "files_count", dir.Path, int64(getFileCountInDirRecursively(dir.Path)))))
			w.Write([]byte(getDirMetric("dirstat", "oldest_file_age", dir.Path, int64(getOldestAgeInDirRecursively(dir.Path)))))
		} else {
			w.Write([]byte(getDirMetric("dirstat", "files_count", dir.Path, int64(getFileCountInDir(dir.Path)))))
			w.Write([]byte(getDirMetric("dirstat", "oldest_file_age", dir.Path, int64(getOldestAgeInDir(dir.Path)))))
		}

	}
}

func getModTime(file string) int64 {
	info, err := os.Stat(file)
	if err != nil {
		log.Print(err)
		return 0
	}
	return info.ModTime().Unix()
}

func getOldestAgeInDirRecursively(dir string) int64 {
	var maxAge int64 = 0
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("an error occurred! %v\n", err)
		}
		if !info.IsDir() {
			age := time.Now().Unix() - getModTime(path)
			if age > maxAge {
				maxAge = age
			}
		}
		return nil
	})
	return maxAge
}

func getOldestAgeInDir(dir string) int64 {
	var files, _ = ioutil.ReadDir(dir)
	var maxAge int64 = 0
	for _, file := range files {
		//fmt.Println(file)
		if !file.IsDir() {
			age := time.Now().Unix() - getModTime(dir+string(os.PathSeparator)+file.Name())
			if age > maxAge {
				maxAge = age
			}
		}
	}
	return maxAge
}

func getFileCountInDirRecursively(dir string) int {
	count := 0
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("an erro occurred: %v\n", err)
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count
}

func getFileCountInDir(dir string) int {
	files, _ := ioutil.ReadDir(dir)
	return len(files)
}

func getDirMetric(namespace string, metricName string, dir string, value int64) string {
	str := fmt.Sprintf("# HELP %s_%s\n", namespace, metricName)
	str += fmt.Sprintf("# TYPE %s_%s counter\n", namespace, metricName)
	str += fmt.Sprintf("%s_%s{dir=\"%s\"} %d\n", namespace, metricName, dir, value)
	return str
}
