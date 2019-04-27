package options

import (
	"sort"
	"strings"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MetricSet map[string]struct{}

func (ms *MetricSet) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := *ms
	ss := s.asSlice()
	sort.Strings(ss)
	return strings.Join(ss, ",")
}
func (ms *MetricSet) Set(value string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := *ms
	metrics := strings.Split(value, ",")
	for _, metric := range metrics {
		metric = strings.TrimSpace(metric)
		if len(metric) != 0 {
			s[metric] = struct{}{}
		}
	}
	return nil
}
func (ms MetricSet) asSlice() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	metrics := []string{}
	for metric := range ms {
		metrics = append(metrics, metric)
	}
	return metrics
}
func (ms MetricSet) IsEmpty() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(ms.asSlice()) == 0
}
func (ms *MetricSet) Type() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "string"
}

type CollectorSet map[string]struct{}

func (c *CollectorSet) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := *c
	ss := s.AsSlice()
	sort.Strings(ss)
	return strings.Join(ss, ",")
}
func (c *CollectorSet) Set(value string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := *c
	cols := strings.Split(value, ",")
	for _, col := range cols {
		col = strings.TrimSpace(col)
		if len(col) != 0 {
			_, ok := DefaultCollectors[col]
			if !ok {
				return fmt.Errorf("collector \"%s\" does not exist", col)
			}
			s[col] = struct{}{}
		}
	}
	return nil
}
func (c CollectorSet) AsSlice() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cols := []string{}
	for col := range c {
		cols = append(cols, col)
	}
	return cols
}
func (c CollectorSet) isEmpty() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(c.AsSlice()) == 0
}
func (c *CollectorSet) Type() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "string"
}

type NamespaceList []string

func (n *NamespaceList) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strings.Join(*n, ",")
}
func (n *NamespaceList) IsAllNamespaces() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(*n) == 1 && (*n)[0] == metav1.NamespaceAll
}
func (n *NamespaceList) Set(value string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	splittedNamespaces := strings.Split(value, ",")
	for _, ns := range splittedNamespaces {
		ns = strings.TrimSpace(ns)
		if len(ns) != 0 {
			*n = append(*n, ns)
		}
	}
	return nil
}
func (n *NamespaceList) Type() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "string"
}
