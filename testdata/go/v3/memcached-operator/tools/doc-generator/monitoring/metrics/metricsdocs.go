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

package main

import (
	"fmt"
	"sort"

	"github.com/example/memcached-operator/monitoring/metrics"
	"github.com/example/memcached-operator/monitoring/metrics/util"
)

// constant parts of the file
const (
	title      = "# Memcached Operator Metrics\n"
	background = "This document aims to help users that are not familiar with metrics exposed by the memcached-operator.\n" +
		"All metrics documented here are auto-generated by the utility tool `tools/doc-generator/monitoring/metrics/metricsdocs` and reflects exactly what is being exposed.\n\n"

	KVSpecificMetrics = "## Memcached Operator Metrics List\n"

	opening = title +
		background +
		KVSpecificMetrics

	// footer
	footerHeading = "## Developing new metrics\n"
	footerContent = "After developing new metrics or changing old ones, please run `make generate-metricsdocs` to regenerate this document.\n\n" +
		"If you feel that the new metric doesn't follow these rules, please change `tools/doc-generator/monitoring/metrics` with your needs.\n"

	footer = footerHeading + footerContent
)

// MemcachedMetricList contains the name, description, and type for each metric.
func main() {
	MemcachedMetricList := metricDescriptionListToMetricList(metrics.ListMetrics())
	sort.Sort(MemcachedMetricList)
	writeToFile(MemcachedMetricList)
}

// writeToFile receives a list of metrics and writes them to a file.
func writeToFile(metricsList metricList) {
	fmt.Print(opening)
	metricsList.writeOut()
	fmt.Print(footer)
}

// metric is an exported struct that defines the metric
// name, description, and type as a new type named metric.
type metric struct {
	name        string
	description string
	metricType  util.MetricType
}

// metricDescriptionToMetric receives a metric of type Metric defined
// in ./monitoring/metrics/util/util.go, and returns as a metric type.
func metricDescriptionToMetric(md util.Metric) metric {
	return metric{
		name:        md.Name,
		description: md.Help,
		metricType:  md.Type,
	}
}

// writeOut receives a metric of type metric and prints
// the metric name, description, and type.
func (m metric) writeOut() {
	fmt.Println("###", m.name)
	fmt.Println(m.description+".", "Type: "+m.metricType+".")
}

// metricList is an array that contain metrics from type metric,
// as a new type named metricList.
type metricList []metric

// metricDescriptionListToMetricList collects the metrics exposed by the
// memcached-operator, and inserts them into the metricList array.
func metricDescriptionListToMetricList(mdl []util.Metric) metricList {
	res := make([]metric, len(mdl))
	for i, md := range mdl {
		res[i] = metricDescriptionToMetric(md)
	}

	return res
}

// Len implements sort.Interface.Len
func (m metricList) Len() int {
	return len(m)
}

// Less implements sort.Interface.Less
func (m metricList) Less(i, j int) bool {
	return m[i].name < m[j].name
}

// Swap implements sort.Interface.Swap
func (m metricList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m metricList) writeOut() {
	for _, met := range m {
		met.writeOut()
	}
}
