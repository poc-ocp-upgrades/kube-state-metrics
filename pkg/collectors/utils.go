package collectors

import (
	"regexp"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/api/core/v1"
	"k8s.io/kube-state-metrics/pkg/metrics"
)

var (
	resyncPeriod				= 5 * time.Minute
	ScrapeErrorTotalMetric		= prometheus.NewCounterVec(prometheus.CounterOpts{Name: "ksm_scrape_error_total", Help: "Total scrape errors encountered when scraping a resource"}, []string{"resource"})
	ResourcesPerScrapeMetric	= prometheus.NewSummaryVec(prometheus.SummaryOpts{Name: "ksm_resources_per_scrape", Help: "Number of resources returned per scrape"}, []string{"resource"})
	invalidLabelCharRE			= regexp.MustCompile(`[^a-zA-Z0-9_]`)
)

func boolFloat64(b bool) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if b {
		return 1
	}
	return 0
}
func addConditionMetrics(cs v1.ConditionStatus) []*metrics.Metric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []*metrics.Metric{&metrics.Metric{LabelValues: []string{"true"}, Value: boolFloat64(cs == v1.ConditionTrue)}, &metrics.Metric{LabelValues: []string{"false"}, Value: boolFloat64(cs == v1.ConditionFalse)}, &metrics.Metric{LabelValues: []string{"unknown"}, Value: boolFloat64(cs == v1.ConditionUnknown)}}
}
func kubeLabelsToPrometheusLabels(labels map[string]string) ([]string, []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	labelKeys := make([]string, len(labels))
	labelValues := make([]string, len(labels))
	i := 0
	for k, v := range labels {
		labelKeys[i] = "label_" + sanitizeLabelName(k)
		labelValues[i] = v
		i++
	}
	return labelKeys, labelValues
}
func kubeAnnotationsToPrometheusAnnotations(annotations map[string]string) ([]string, []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	annotationKeys := make([]string, len(annotations))
	annotationValues := make([]string, len(annotations))
	i := 0
	for k, v := range annotations {
		annotationKeys[i] = "annotation_" + sanitizeLabelName(k)
		annotationValues[i] = v
		i++
	}
	return annotationKeys, annotationValues
}
func sanitizeLabelName(s string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return invalidLabelCharRE.ReplaceAllString(s, "_")
}
