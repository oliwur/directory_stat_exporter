package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/codestoke/directory_stat_exporter/cfg"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/exporter-toolkit/web"
)

type metricValue struct {
	labels    map[string]string
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

const (
	namespace              = "dirstat"
	metricFilesInDir       = "files_in_dir"
	metricOldestFileTime   = "oldest_file_time"
	metricCurrentTimestamp = "current_timestamp"
)

var (
	config           cfg.Config
	metricRegister   map[string]metric
	currentTimestamp metric
	cache            []byte
	lastRequest      time.Time
	cacheLock        sync.Mutex
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config.file", "config.yml", "provide a custom config file")
	var (
		tlsConfigFile = flag.String("web.config", "", "Path to config yaml file that can enable TLS or authentication.")
	)
	flag.Parse()

	config = cfg.GetConfig(configFile)

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
		metricHelp:   "displays the timestamp in unix time of the oldest file",
		metricValues: make(map[string]metricValue),
	}

	lastRequest = time.Unix(0, 0)
	cache = []byte("# dirstat")
	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)

	http.HandleFunc("/metrics", handleMetrics)

	server := &http.Server{Addr: ":" + config.ServicePort}
	if err := web.ListenAndServe(server, *tlsConfigFile, logger); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
		os.Exit(1)
	}
}

func handleMetrics(w http.ResponseWriter, _ *http.Request) {
	if lastRequest.Add(time.Minute*time.Duration(config.CacheTime)).Unix() < time.Now().Unix() {
		// update the cache
		{
			writeMetricsResponse(w)
			go updateMetrics()
		}
	} else {
		// respond with the cache.
		writeMetricsResponse(w)
	}
}

func getMetricValues() []byte {
	res := sprintDirMetric(currentTimestamp)
	for _, value := range metricRegister {
		res += sprintDirMetric(value)
	}
	return []byte(res)
}

func writeMetricsResponse(w http.ResponseWriter) {
	_, _ = w.Write(cache)
}

func updateMetrics() {
	cacheLock.Lock()

	for _, dir := range config.Directories {
		if dir.Recursive {
			metricRegister[metricFilesInDir].metricValues[dir.Path] = metricValue{
				value:     int64(getFileCountInDirRecursively(dir.Path)),
				recursive: dir.Recursive,
				name:      dir.Name,
				labels: map[string]string{
					"dir":       dir.Name,
					"recursive": strconv.FormatBool(dir.Recursive),
				},
			}
			metricRegister[metricOldestFileTime].metricValues[dir.Path] = metricValue{
				value:     int64(getOldestAgeInDirRecursively(dir.Path)),
				recursive: dir.Recursive,
				name:      dir.Name,
				labels: map[string]string{
					"dir":       dir.Name,
					"recursive": strconv.FormatBool(dir.Recursive),
				},
			}
		} else {
			metricRegister[metricFilesInDir].metricValues[dir.Path] = metricValue{
				value:     int64(getFileCountInDir(dir.Path)),
				recursive: dir.Recursive,
				name:      dir.Name,
				labels: map[string]string{
					"dir":       dir.Name,
					"recursive": strconv.FormatBool(dir.Recursive),
				},
			}
			metricRegister[metricOldestFileTime].metricValues[dir.Path] = metricValue{
				value:     int64(getOldestAgeInDir(dir.Path)),
				recursive: dir.Recursive,
				name:      dir.Name,
				labels: map[string]string{
					"dir":       dir.Name,
					"recursive": strconv.FormatBool(dir.Recursive),
				},
			}
		}
	}
	currentTimestamp = metric{
		metricName:   metricCurrentTimestamp,
		metricHelp:   "the current timestamp in unix time.",
		metricType:   "gauge",
		metricValues: map[string]metricValue{"ts": {value: time.Now().Unix()}},
	}

	cache = getMetricValues()
	lastRequest = time.Now()

	cacheLock.Unlock()
}

// this should be replaced with one more generic generator.
//func sprintCurrentTimestamp(m metric) string {
//	str := ""
//	str += fmt.Sprintf("# HELP %s_%s %s\n", namespace, m.metricName, m.metricHelp)
//	str += fmt.Sprintf("# TYPE %s_%s %s\n", namespace, m.metricName, m.metricType)
//	for _, v := range m.metricValues {
//		str += fmt.Sprintf("%s_%s %v\n", namespace, m.metricName, v.value)
//	}
//	return str
//}

func sprintDirMetric(m metric) string {
	str := ""
	str += fmt.Sprintf("# HELP %s_%s %s\n", namespace, m.metricName, m.metricHelp)
	str += fmt.Sprintf("# TYPE %s_%s %s\n", namespace, m.metricName, m.metricType)
	for _, v := range m.metricValues {
		//str += fmt.Sprintf("%s_%s{dir=\"%s\",recursive=\"%t\"} %v\n", namespace, m.metricName, v.name, v.recursive, v.value)
		str += sprintMetric(namespace, m.metricName, v.value, v.labels)
	}
	return str
}

func sprintMetric(ns string, name string, value int64, labels map[string]string) string {
	strLbls := ""
	if labels != nil {
		var lblArr []string
		strLbls += "{"
		for k, v := range labels {
			lblArr = append(lblArr, fmt.Sprintf("%s=\"%s\"", k, v))
		}
		strLbls += strings.Join(lblArr, ",")
		strLbls += "}"
	}
	str := fmt.Sprintf("%s_%s%s %v\n", ns, name, strLbls, value)
	return str
}
