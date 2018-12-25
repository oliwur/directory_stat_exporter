package main

import (
	"fmt"
	"github.com/codestoke/directory_stat_exporter/cfg"
	"net/http"
	"strings"
	"time"
)

// globals
var config cfg.Config

const (
	namespace              = "dirstat"
	metricFilesInDir       = "files_in_dir"
	metricOldestFileTime   = "oldest_file_time"
	metricCurrentTimestamp = "current_timestamp"
)

type metricValue struct {
	dir       string
	name      string
	value     int64
	recursive bool
}
type metric struct {
	metricName   string
	metricHelp   string
	metricType   string
	metricValues map[string]metricValue
}

var metricRegister map[string]metric
var currentTimestamp metric

func main() {
	config = cfg.GetConfig()

	metricRegister = make(map[string]metric)

	metricRegister[metricFilesInDir] = metric{
		metricName:   metricFilesInDir,
		metricType:   "gauge",
		metricHelp:   "this counts all the files in a directory",
		metricValues: make(map[string]metricValue),
	}
	metricRegister[metricOldestFileTime] = metric{
		metricName:   metricOldestFileTime,
		metricType:   "gauge",
		metricHelp:   "displays the timestamp in unix time of the oldes file",
		metricValues: make(map[string]metricValue),
	}

	http.HandleFunc("/metrics", handleMetrics)
	if err := http.ListenAndServe(":"+config.ServicePort, nil); err != nil {
		panic(err)
	}
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	for _, dir := range config.Directories {
		if dir.Recursive {
			metricRegister[metricFilesInDir].metricValues[dir.Path] = metricValue{
				value:     int64(getFileCountInDirRecursively(dir.Path)),
				recursive: dir.Recursive,
				name:      dir.Name,
			}
			metricRegister[metricOldestFileTime].metricValues[dir.Path] = metricValue{
				value:     int64(getOldestAgeInDirRecursively(dir.Path)),
				recursive: dir.Recursive,
				name:      dir.Name,
			}
		} else {
			metricRegister[metricFilesInDir].metricValues[dir.Path] = metricValue{
				value:     int64(getFileCountInDir(dir.Path)),
				recursive: dir.Recursive,
				name:      dir.Name,
			}
			metricRegister[metricOldestFileTime].metricValues[dir.Path] = metricValue{
				value:     int64(getOldestAgeInDir(dir.Path)),
				recursive: dir.Recursive,
				name:      dir.Name,
			}
		}
	}

	currentTimestamp = metric{
		metricName:   metricCurrentTimestamp,
		metricHelp:   "the current timestamp in unix time.",
		metricType:   "gauge",
		metricValues: map[string]metricValue{"ts": {value: time.Now().Unix()}},
	}
	_, _ = w.Write([]byte(sprintCurrentTimestamp(currentTimestamp)))

	for _, value := range metricRegister {
		_, _ = w.Write([]byte(sprintDirMetric(value)))
	}
}

// this should be replaced with one more generic generator.
func sprintCurrentTimestamp(m metric) string {
	str := ""
	str += fmt.Sprintf("# HELP %s_%s %s\n", namespace, m.metricName, m.metricHelp)
	str += fmt.Sprintf("# TYPE %s_%s %s\n", namespace, m.metricName, m.metricType)
	for _, v := range m.metricValues {
		str += fmt.Sprintf("%s_%s %v\n", namespace, m.metricName, v.value)
	}
	return str
}

func sprintDirMetric(m metric) string {
	str := ""
	str += fmt.Sprintf("# HELP %s_%s %s\n", namespace, m.metricName, m.metricHelp)
	str += fmt.Sprintf("# TYPE %s_%s %s\n", namespace, m.metricName, m.metricType)
	for _, v := range m.metricValues {
		str += fmt.Sprintf("%s_%s{dir=\"%s\",recursive=\"%t\"} %v\n", namespace, m.metricName, v.name, v.recursive, v.value)
	}
	return str
}

func sprintMetric(ns string, name string, value int64, labels map[string]string) string {
	strLbls := ""
	if labels != nil {
		lblArr := make([]string, 0)
		strLbls += "{"
		for k, v := range labels {
			lblArr = append(lblArr, fmt.Sprintf("%s=\"%s\"", k, v))
		}
		strLbls += strings.Join(lblArr, ",")
		strLbls += "}"
	}
	str := fmt.Sprintf("%s_%s%s %v", ns, name, strLbls, value)
	return str
}
