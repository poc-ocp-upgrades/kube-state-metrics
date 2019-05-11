package collectors

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"k8s.io/kube-state-metrics/pkg/metrics_store"
)

type generateMetricsTestCase struct {
	Obj			interface{}
	MetricNames	[]string
	Want		string
	Func		func(interface{}) []metricsstore.FamilyStringer
}

func (testCase *generateMetricsTestCase) run() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	metricFamilies := testCase.Func(testCase.Obj)
	metricFamilyStrings := []string{}
	for _, f := range metricFamilies {
		metricFamilyStrings = append(metricFamilyStrings, f.String())
	}
	metrics := strings.Split(strings.Join(metricFamilyStrings, ""), "\n")
	metrics = filterMetrics(metrics, testCase.MetricNames)
	out := strings.Join(metrics, "\n")
	if err := compareOutput(testCase.Want, out); err != nil {
		return fmt.Errorf("expected wanted output to equal output: %v", err.Error())
	}
	return nil
}
func compareOutput(a, b string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	entities := []string{a, b}
	for i := 0; i < len(entities); i++ {
		for _, f := range []func(string) string{removeUnusedWhitespace, sortLabels, sortByLine} {
			entities[i] = f(entities[i])
		}
	}
	if entities[0] != entities[1] {
		return fmt.Errorf("expected a to equal b but got:\n%v\nand:\n%v", entities[0], entities[1])
	}
	return nil
}
func sortLabels(s string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sorted := []string{}
	for _, line := range strings.Split(s, "\n") {
		split := strings.Split(line, "{")
		if len(split) != 2 {
			panic(fmt.Sprintf("failed to sort labels in \"%v\"", line))
		}
		name := split[0]
		split = strings.Split(split[1], "}")
		value := split[1]
		labels := strings.Split(split[0], ",")
		sort.Strings(labels)
		sorted = append(sorted, fmt.Sprintf("%v{%v}%v", name, strings.Join(labels, ","), value))
	}
	return strings.Join(sorted, "\n")
}
func sortByLine(s string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	split := strings.Split(s, "\n")
	sort.Strings(split)
	return strings.Join(split, "\n")
}
func filterMetrics(ms []string, names []string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if names == nil {
		return ms
	}
	filtered := []string{}
	regexps := []*regexp.Regexp{}
	for _, n := range names {
		regexps = append(regexps, regexp.MustCompile(fmt.Sprintf("^%v", n)))
	}
	for _, m := range ms {
		drop := true
		for _, r := range regexps {
			if r.MatchString(m) {
				drop = false
				break
			}
		}
		if !drop {
			filtered = append(filtered, m)
		}
	}
	return filtered
}
func removeUnusedWhitespace(s string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		trimmedLine		string
		trimmedLines	[]string
		lines			= strings.Split(s, "\n")
	)
	for _, l := range lines {
		trimmedLine = strings.TrimSpace(l)
		if len(trimmedLine) > 0 {
			trimmedLines = append(trimmedLines, trimmedLine)
		}
	}
	return strings.Join(trimmedLines, "\n")
}