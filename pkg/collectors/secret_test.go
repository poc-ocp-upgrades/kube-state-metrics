package collectors

import (
	"testing"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSecretCollector(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	startTime := 1501569018
	metav1StartTime := metav1.Unix(int64(startTime), 0)
	const metadata = `
        # HELP kube_secret_labels Kubernetes labels converted to Prometheus labels.
		# TYPE kube_secret_labels gauge
        # HELP kube_secret_info Information about secret.
		# TYPE kube_secret_info gauge
		# HELP kube_secret_type Type about secret.
		# TYPE kube_secret_type gauge
		# HELP kube_secret_created Unix creation timestamp
		# TYPE kube_secret_created gauge
		# HELP kube_secret_metadata_resource_version Resource version representing a specific version of secret.
		# TYPE kube_secret_metadata_resource_version gauge
	`
	cases := []generateMetricsTestCase{{Obj: &v1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "secret1", Namespace: "ns1", ResourceVersion: "000000"}, Type: v1.SecretTypeOpaque}, Want: `
				kube_secret_info{namespace="ns1",secret="secret1"} 1
				kube_secret_type{namespace="ns1",secret="secret1",type="Opaque"} 1
				kube_secret_metadata_resource_version{namespace="ns1",resource_version="000000",secret="secret1"} 1
				kube_secret_labels{namespace="ns1",secret="secret1"} 1
`, MetricNames: []string{"kube_secret_info", "kube_secret_metadata_resource_version", "kube_secret_created", "kube_secret_labels", "kube_secret_type"}}, {Obj: &v1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "secret2", Namespace: "ns2", CreationTimestamp: metav1StartTime, ResourceVersion: "123456"}, Type: v1.SecretTypeServiceAccountToken}, Want: `
				kube_secret_info{namespace="ns2",secret="secret2"} 1
				kube_secret_type{namespace="ns2",secret="secret2",type="kubernetes.io/service-account-token"} 1
				kube_secret_created{namespace="ns2",secret="secret2"} 1.501569018e+09
				kube_secret_metadata_resource_version{namespace="ns2",resource_version="123456",secret="secret2"} 1
				kube_secret_labels{namespace="ns2",secret="secret2"} 1
				`, MetricNames: []string{"kube_secret_info", "kube_secret_metadata_resource_version", "kube_secret_created", "kube_secret_labels", "kube_secret_type"}}, {Obj: &v1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "secret3", Namespace: "ns3", CreationTimestamp: metav1StartTime, Labels: map[string]string{"test-3": "test-3"}, ResourceVersion: "abcdef"}, Type: v1.SecretTypeDockercfg}, Want: `
				kube_secret_info{namespace="ns3",secret="secret3"} 1
				kube_secret_type{namespace="ns3",secret="secret3",type="kubernetes.io/dockercfg"} 1
				kube_secret_created{namespace="ns3",secret="secret3"} 1.501569018e+09
				kube_secret_metadata_resource_version{namespace="ns3",resource_version="abcdef",secret="secret3"} 1
				kube_secret_labels{label_test_3="test-3",namespace="ns3",secret="secret3"} 1
`, MetricNames: []string{"kube_secret_info", "kube_secret_metadata_resource_version", "kube_secret_created", "kube_secret_labels", "kube_secret_type"}}}
	for i, c := range cases {
		c.Func = composeMetricGenFuncs(secretMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
