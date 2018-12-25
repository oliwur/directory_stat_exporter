package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestMetricsWriter(t *testing.T) {
	t.Run("should write single metric without label", func(t *testing.T) {
		namespace := "dirstat"
		name := "name"
		value := 1

		txt := sprintMetric(namespace, name, int64(value), nil)

		expected := fmt.Sprintf("%v_%v %v", namespace, name, value)
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

		expected := fmt.Sprintf("%s_%s{%s=\"%s\"} %v", namespace, name, lblKey, lblValue, value)
		if txt != expected {
			t.Fail()
			t.Errorf("the expected text was not retured:\nexpected: %v\nreturned: %v\n", expected, txt)
		}
	})

	t.Run("should write single metric with a multiple labels", func(t *testing.T) {
		namespace := "dirstat"
		name := "name"
		value := 1

		lbls := make(map[string]string, 0)
		lbls["lblKey1"] = "lblValue1"
		lbls["lblKey2"] = "lblValue2"

		keys := make([]string, 0)
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

		expected := fmt.Sprintf("%s_%s%s %v", namespace, name, lblTxt, value)
		if txt != expected {
			t.Fail()
			t.Errorf("the expected text was not retured:\nexpected: %v\nreturned: %v\n", expected, txt)
		}
	})
}
