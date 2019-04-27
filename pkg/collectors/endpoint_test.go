package collectors

import (
	"testing"
	"time"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEndpointCollector(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const metadata = `
		# HELP kube_endpoint_address_available Number of addresses available in endpoint.
		# TYPE kube_endpoint_address_available gauge
		# HELP kube_endpoint_address_not_ready Number of addresses not ready in endpoint
		# TYPE kube_endpoint_address_not_ready gauge
		# HELP kube_endpoint_created Unix creation timestamp
		# TYPE kube_endpoint_created gauge
		# HELP kube_endpoint_info Information about endpoint.
		# TYPE kube_endpoint_info gauge
		# HELP kube_endpoint_labels Kubernetes labels converted to Prometheus labels.
		# TYPE kube_endpoint_labels gauge
	`
	cases := []generateMetricsTestCase{{Obj: &v1.Endpoints{ObjectMeta: metav1.ObjectMeta{Name: "test-endpoint", CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)}, Namespace: "default", Labels: map[string]string{"app": "foobar"}}, Subsets: []v1.EndpointSubset{{Addresses: []v1.EndpointAddress{{IP: "127.0.0.1"}, {IP: "10.0.0.1"}}, Ports: []v1.EndpointPort{{Port: 8080}, {Port: 8081}}}, {Addresses: []v1.EndpointAddress{{IP: "172.22.23.202"}}, Ports: []v1.EndpointPort{{Port: 8443}, {Port: 9090}}}, {NotReadyAddresses: []v1.EndpointAddress{{IP: "192.168.1.1"}}, Ports: []v1.EndpointPort{{Port: 1234}, {Port: 5678}}}, {NotReadyAddresses: []v1.EndpointAddress{{IP: "192.168.1.3"}, {IP: "192.168.2.2"}}, Ports: []v1.EndpointPort{{Port: 1234}, {Port: 5678}}}}}, Want: `
				kube_endpoint_address_available{endpoint="test-endpoint",namespace="default"} 6
				kube_endpoint_address_not_ready{endpoint="test-endpoint",namespace="default"} 6
				kube_endpoint_created{endpoint="test-endpoint",namespace="default"} 1.5e+09
				kube_endpoint_info{endpoint="test-endpoint",namespace="default"} 1
				kube_endpoint_labels{endpoint="test-endpoint",label_app="foobar",namespace="default"} 1
			`}}
	for i, c := range cases {
		c.Func = composeMetricGenFuncs(endpointMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
