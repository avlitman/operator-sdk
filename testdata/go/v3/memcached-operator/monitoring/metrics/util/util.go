/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type MetricType string

const (
	Counter MetricType = "counter"
	Gauge   MetricType = "gauge"
)

type Metric struct {
	Name string
	Help string
	Type MetricType
}

var collectors map[string]prometheus.Collector

func RegisterMetrics(allMetrics [][]Metric) {
	collectors = map[string]prometheus.Collector{}

	for _, m := range flatMetrics(allMetrics) {
		v := createMetric(m)
		metrics.Registry.MustRegister(v)
		collectors[m.Name] = v
	}
}

func ListMetrics(allMetrics [][]Metric) []Metric {
	return flatMetrics(allMetrics)
}

func GetCounterMetric(metricName string) prometheus.Counter {
	return castCounter(metricName, getMetric(metricName))
}

func GetGaugeMetric(metricName string) prometheus.Gauge {
	return castGauge(metricName, getMetric(metricName))
}

func flatMetrics(allMetrics [][]Metric) []Metric {
	var result []Metric
	for _, metricFile := range allMetrics {
		result = append(result, metricFile...)
	}
	return result
}

func createMetric(metric Metric) prometheus.Collector {
	opts := prometheus.Opts{
		Name: metric.Name,
		Help: metric.Help,
	}

	switch metric.Type {
	case Counter:
		return prometheus.NewCounter(prometheus.CounterOpts(opts))
	case Gauge:
		return prometheus.NewCounter(prometheus.CounterOpts(opts))
	}

	panic(unknownMetricTypeError(metric.Name, string(metric.Type)))
}

func getMetric(metricName string) prometheus.Collector {
	metric, found := collectors[metricName]
	if !found {
		panic(unknownMetricNameError(metricName))
	}
	return metric
}

func castCounter(metricName string, metric prometheus.Collector) prometheus.Counter {
	v, ok := metric.(prometheus.Counter)
	if !ok {
		panic(unknownMetricTypeError(metricName, "Counter"))
	}
	return v
}

func castGauge(metricName string, metric prometheus.Collector) prometheus.Gauge {
	v, ok := metric.(prometheus.Gauge)
	if !ok {
		panic(unknownMetricTypeError(metricName, "Gauge"))
	}
	return v
}

func unknownMetricNameError(metricName string) error {
	return fmt.Errorf("unknown metric name %s", metricName)
}

func unknownMetricTypeError(metricName string, requestedType string) error {
	return fmt.Errorf("%s is not requested %s metric type", metricName, requestedType)
}
