package collectors

import (
	"testing"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConfigMapCollector(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	startTime := 1501569018
	metav1StartTime := metav1.Unix(int64(startTime), 0)
	const metadata = `
        # HELP kube_configmap_info Information about configmap.
		# TYPE kube_configmap_info gauge
		# HELP kube_configmap_created Unix creation timestamp
		# TYPE kube_configmap_created gauge
		# HELP kube_configmap_metadata_resource_version Resource version representing a specific version of the configmap.
		# TYPE kube_configmap_metadata_resource_version gauge
	`
	cases := []generateMetricsTestCase{{Obj: &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "configmap1", Namespace: "ns1", ResourceVersion: "123456"}}, Want: `
				kube_configmap_info{configmap="configmap1",namespace="ns1"} 1
				kube_configmap_metadata_resource_version{configmap="configmap1",namespace="ns1",resource_version="123456"} 1
`, MetricNames: []string{"kube_configmap_info", "kube_configmap_metadata_resource_version"}}, {Obj: &v1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "configmap2", Namespace: "ns2", CreationTimestamp: metav1StartTime, ResourceVersion: "abcdef"}}, Want: `
				kube_configmap_info{configmap="configmap2",namespace="ns2"} 1
				kube_configmap_created{configmap="configmap2",namespace="ns2"} 1.501569018e+09
				kube_configmap_metadata_resource_version{configmap="configmap2",namespace="ns2",resource_version="abcdef"} 1
				`, MetricNames: []string{"kube_configmap_info", "kube_configmap_created", "kube_configmap_metadata_resource_version"}}}
	for i, c := range cases {
		c.Func = composeMetricGenFuncs(configMapMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
