package collectors

import (
	"testing"
	"time"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLimitRangeollector(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	testMemory := "2.1G"
	testMemoryQuantity := resource.MustParse(testMemory)
	const metadata = `
	# HELP kube_limitrange_created Unix creation timestamp
	# TYPE kube_limitrange_created gauge
	# HELP kube_limitrange Information about limit range.
	# TYPE kube_limitrange gauge
	`
	cases := []generateMetricsTestCase{{Obj: &v1.LimitRange{ObjectMeta: metav1.ObjectMeta{Name: "quotaTest", CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)}, Namespace: "testNS"}, Spec: v1.LimitRangeSpec{Limits: []v1.LimitRangeItem{{Type: v1.LimitTypePod, Max: map[v1.ResourceName]resource.Quantity{v1.ResourceMemory: testMemoryQuantity}, Min: map[v1.ResourceName]resource.Quantity{v1.ResourceMemory: testMemoryQuantity}, Default: map[v1.ResourceName]resource.Quantity{v1.ResourceMemory: testMemoryQuantity}, DefaultRequest: map[v1.ResourceName]resource.Quantity{v1.ResourceMemory: testMemoryQuantity}, MaxLimitRequestRatio: map[v1.ResourceName]resource.Quantity{v1.ResourceMemory: testMemoryQuantity}}}}}, Want: `
        kube_limitrange_created{limitrange="quotaTest",namespace="testNS"} 1.5e+09
        kube_limitrange{constraint="default",limitrange="quotaTest",namespace="testNS",resource="memory",type="Pod"} 2.1e+09
        kube_limitrange{constraint="defaultRequest",limitrange="quotaTest",namespace="testNS",resource="memory",type="Pod"} 2.1e+09
        kube_limitrange{constraint="max",limitrange="quotaTest",namespace="testNS",resource="memory",type="Pod"} 2.1e+09
        kube_limitrange{constraint="maxLimitRequestRatio",limitrange="quotaTest",namespace="testNS",resource="memory",type="Pod"} 2.1e+09
        kube_limitrange{constraint="min",limitrange="quotaTest",namespace="testNS",resource="memory",type="Pod"} 2.1e+09

		`}}
	for i, c := range cases {
		c.Func = composeMetricGenFuncs(limitRangeMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
