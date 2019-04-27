package collectors

import (
	"testing"
	autoscaling "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	hpa1MinReplicas int32 = 2
)

func TestHPACollector(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	const metadata = `
		# HELP kube_hpa_metadata_generation The generation observed by the HorizontalPodAutoscaler controller.
		# TYPE kube_hpa_metadata_generation gauge
		# HELP kube_hpa_spec_max_replicas Upper limit for the number of pods that can be set by the autoscaler; cannot be smaller than MinReplicas.
		# TYPE kube_hpa_spec_max_replicas gauge
		# HELP kube_hpa_spec_min_replicas Lower limit for the number of pods that can be set by the autoscaler, default 1.
		# TYPE kube_hpa_spec_min_replicas gauge
		# HELP kube_hpa_status_current_replicas Current number of replicas of pods managed by this autoscaler.
		# TYPE kube_hpa_status_current_replicas gauge
		# HELP kube_hpa_status_desired_replicas Desired number of replicas of pods managed by this autoscaler.
		# TYPE kube_hpa_status_desired_replicas gauge
	`
	cases := []generateMetricsTestCase{{Obj: &autoscaling.HorizontalPodAutoscaler{ObjectMeta: metav1.ObjectMeta{Generation: 2, Name: "hpa1", Namespace: "ns1"}, Spec: autoscaling.HorizontalPodAutoscalerSpec{MaxReplicas: 4, MinReplicas: &hpa1MinReplicas, ScaleTargetRef: autoscaling.CrossVersionObjectReference{APIVersion: "extensions/v1beta1", Kind: "Deployment", Name: "deployment1"}}, Status: autoscaling.HorizontalPodAutoscalerStatus{CurrentReplicas: 2, DesiredReplicas: 2}}, Want: `
				kube_hpa_metadata_generation{hpa="hpa1",namespace="ns1"} 2
				kube_hpa_spec_max_replicas{hpa="hpa1",namespace="ns1"} 4
				kube_hpa_spec_min_replicas{hpa="hpa1",namespace="ns1"} 2
				kube_hpa_status_current_replicas{hpa="hpa1",namespace="ns1"} 2
				kube_hpa_status_desired_replicas{hpa="hpa1",namespace="ns1"} 2
			`, MetricNames: []string{"kube_hpa_metadata_generation", "kube_hpa_spec_max_replicas", "kube_hpa_spec_min_replicas", "kube_hpa_status_current_replicas", "kube_hpa_status_desired_replicas"}}}
	for i, c := range cases {
		c.Func = composeMetricGenFuncs(hpaMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
