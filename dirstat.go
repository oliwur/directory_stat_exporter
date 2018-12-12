package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/metrics", handleMetrics)
	if err := http.ListenAndServe(":9999", nil); err != nil {
		panic(err)
	}
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	dir := "\\"
	w.Write([]byte(getDirMetric("dirstat", "files_count", dir, int64(getFileCountInDir(dir)))))
	w.Write([]byte(getDirMetric("dirstat", "oldest_file_age", dir, int64(getOldestAgeInDir(dir)))))
}

func getModTime(file string) int64 {
	info, err := os.Stat(file)
	if err != nil {
		log.Print(err)
		return 0
	}
	return info.ModTime().Unix()
}

func getOldestAgeInDir(dir string) int64 {
	var files, _ = ioutil.ReadDir(dir)
	var maxAge int64 = 0
	for _, file := range files {
		//fmt.Println(file)
		age := time.Now().Unix() - getModTime(dir + string(os.PathSeparator) + file.Name())
		if age > maxAge {
			maxAge = age
		}
	}
	return maxAge
}

func getFileCountInDir(dir string) int {
	files, _ := ioutil.ReadDir(dir)
	return len(files)
}

func getDirMetric(namespace string, metricName string, dir string, value int64) string {
	str := fmt.Sprintf("# HELP %s_%s\n", namespace, metricName)
	str += fmt.Sprintf("# TYPE %s_%s\n", namespace, metricName)
	str += fmt.Sprintf("%s_%s{dir=\"%s\"} %d\n", namespace, metricName, dir, value)
	return str
}
