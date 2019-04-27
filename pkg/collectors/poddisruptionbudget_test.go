package collectors

import (
	"testing"
	"time"
	"k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodDisruptionBudgetCollector(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const metadata = `
	# HELP kube_poddisruptionbudget_created Unix creation timestamp
	# TYPE kube_poddisruptionbudget_created gauge
	# HELP kube_poddisruptionbudget_status_current_healthy Current number of healthy pods
	# TYPE kube_poddisruptionbudget_status_current_healthy gauge
	# HELP kube_poddisruptionbudget_status_desired_healthy Minimum desired number of healthy pods
	# TYPE kube_poddisruptionbudget_status_desired_healthy gauge
	# HELP kube_poddisruptionbudget_status_pod_disruptions_allowed Number of pod disruptions that are currently allowed
	# TYPE kube_poddisruptionbudget_status_pod_disruptions_allowed gauge
	# HELP kube_poddisruptionbudget_status_expected_pods Total number of pods counted by this disruption budget
	# TYPE kube_poddisruptionbudget_status_expected_pods gauge
	# HELP kube_poddisruptionbudget_status_observed_generation Most recent generation observed when updating this PDB status
	# TYPE kube_poddisruptionbudget_status_observed_generation gauge
	`
	cases := []generateMetricsTestCase{{Obj: &v1beta1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: "pdb1", CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)}, Namespace: "ns1", Generation: 21}, Status: v1beta1.PodDisruptionBudgetStatus{CurrentHealthy: 12, DesiredHealthy: 10, PodDisruptionsAllowed: 2, ExpectedPods: 15, ObservedGeneration: 111}}, Want: `
			kube_poddisruptionbudget_created{namespace="ns1",poddisruptionbudget="pdb1"} 1.5e+09
			kube_poddisruptionbudget_status_current_healthy{namespace="ns1",poddisruptionbudget="pdb1"} 12
			kube_poddisruptionbudget_status_desired_healthy{namespace="ns1",poddisruptionbudget="pdb1"} 10
			kube_poddisruptionbudget_status_pod_disruptions_allowed{namespace="ns1",poddisruptionbudget="pdb1"} 2
			kube_poddisruptionbudget_status_expected_pods{namespace="ns1",poddisruptionbudget="pdb1"} 15
			kube_poddisruptionbudget_status_observed_generation{namespace="ns1",poddisruptionbudget="pdb1"} 111
			`}, {Obj: &v1beta1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: "pdb2", Namespace: "ns2", Generation: 14}, Status: v1beta1.PodDisruptionBudgetStatus{CurrentHealthy: 8, DesiredHealthy: 9, PodDisruptionsAllowed: 0, ExpectedPods: 10, ObservedGeneration: 1111}}, Want: `
				kube_poddisruptionbudget_status_current_healthy{namespace="ns2",poddisruptionbudget="pdb2"} 8
				kube_poddisruptionbudget_status_desired_healthy{namespace="ns2",poddisruptionbudget="pdb2"} 9
				kube_poddisruptionbudget_status_pod_disruptions_allowed{namespace="ns2",poddisruptionbudget="pdb2"} 0
				kube_poddisruptionbudget_status_expected_pods{namespace="ns2",poddisruptionbudget="pdb2"} 10
				kube_poddisruptionbudget_status_observed_generation{namespace="ns2",poddisruptionbudget="pdb2"} 1111
			`}}
	for i, c := range cases {
		c.Func = composeMetricGenFuncs(podDisruptionBudgetMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
