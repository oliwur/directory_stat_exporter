package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestMetricsWriter(t *testing.T) {
	t.Run("each metric must end with a line feed (\\n)", func(t *testing.T) {
		namespace := "namespace"
		name := "name"
		value := 1

		txt := sprintMetric(namespace, name, int64(value), nil)

		//expected := fmt.Sprintf("%v_%v %v\n", namespace, name, value)
		if !strings.HasSuffix(txt, "\n") {
			t.Fail()
			t.Errorf("the generated code must end in a line feed (\\n)")
		}
	})

	t.Run("should write single metric without label", func(t *testing.T) {
		namespace := "dirstat"
		name := "name"
		value := 1

		txt := sprintMetric(namespace, name, int64(value), nil)

		expected := fmt.Sprintf("%v_%v %v\n", namespace, name, value)
		if txt != expected {
			t.Fail()
			t.Errorf("the expected text was not retured:\nexpected: %v\nreturned: %v\n", expected, txt)
		}
	})

	t.Run("should write single metric with a single label", func(t *testing.T) {
		namespace := "dirstat"
		name := "name"
		value := 1

		lblKey := "lbl"
		lblValue := "lblValue"

		lbls := make(map[string]string, 1)
		lbls[lblKey] = lblValue

		txt := sprintMetric(namespace, name, int64(value), lbls)

		expected := fmt.Sprintf("%s_%s{%s=\"%s\"} %v\n", namespace, name, lblKey, lblValue, value)
		if txt != expected {
			t.Fail()
			t.Errorf("the expected text was not retured:\nexpected: %v\nreturned: %v\n", expected, txt)
		}
	})

	t.Run("should write single metric with a multiple labels", func(t *testing.T) {
		namespace := "dirstat"
		name := "name"
		value := 1

		lbls := make(map[string]string, 2)
		lbls["lblKey1"] = "lblValue1"
		lbls["lblKey2"] = "lblValue2"

		var keys []string
		for k := range lbls {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		lblSlce := make([]string, 0)
		for _, k := range keys {
			lblSlce = append(lblSlce, fmt.Sprintf("%s=\"%s\"", k, lbls[k]))
		}
		lblTxt := fmt.Sprintf("{%s}", strings.Join(lblSlce, ","))

		txt := sprintMetric(namespace, name, int64(value), lbls)

		expected := fmt.Sprintf("%s_%s%s %v\n", namespace, name, lblTxt, value)
		if txt != expected {
			t.Fail()
			t.Errorf("the expected text was not retured:\nexpected: %v\nreturned: %v\n", expected, txt)
		}
	})
}

func TestDirMetric(t *testing.T) {
	// setup some metrics to test.
	m := metric{
		metricName: "name",
		metricType: "type",
		metricValues: map[string]metricValue{
			"name": metricValue{
				value: 1,
				labels: map[string]string{
					"dir":       "dir",
					"recursive": strconv.FormatBool(false),
				},
			},
		},
	}

	returned := sprintDirMetric(m)

	t.Run("given a not empty metric when response is generated then the result must not be empty", func(t *testing.T) {
		if len(returned) == 0 {
			t.Fail()
			t.Errorf("the result metric string is empty")
		}
	})

	t.Run("given a not empty metric when response is generated then the result must start with expected value", func(t *testing.T) {
		expectedStart := fmt.Sprintf("# HELP %[1]s_%[2]s\n# TYPE %[1]s_%[2]s type\n%[1]s_%[2]s", namespace, m.metricName)

		if strings.HasPrefix(returned, expectedStart) {
			t.Fail()
			t.Errorf("it must be in the correct format.\n\texpected: %s\n\treturned: %s\n", expectedStart, returned)
		}
	})

	t.Run("given a not empty metric with labels when response is generated then the labels must be in the result string", func(t *testing.T) {
		r := regexp.MustCompile("{([^}]+)}")

		matches := r.FindAllStringSubmatch(returned, -1) // i only want the first one

		if len(matches) == 0 {
			t.Fail()
			t.Errorf("I'd expect some labels.\n%s", returned)
		} else {
			// it has labels, now parse and test if they are correct$
			labels := strings.Split(matches[0][1], ",")
			for _, label := range labels {
				fmt.Println("testing", label)
				keyValue := strings.Split(label, "=")
				key := keyValue[0]
				value := strings.Replace(keyValue[1], "\"", "", -1)

				fmt.Println("key, value", key, value)

				if value != m.metricValues["name"].labels[key] {
					t.Fail()
					t.Errorf("the label does not exist or the label does not contain the correct value.\n%s\n", returned)
					t.Error(m.metricValues["name"].labels)
				}
			}
		}
	})
}
