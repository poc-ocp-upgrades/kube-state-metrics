package collectors

import (
	"k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kube-state-metrics/pkg/metrics"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientset "k8s.io/client-go/kubernetes"
)

var (
	descPodDisruptionBudgetLabelsDefaultLabels	= []string{"namespace", "poddisruptionbudget"}
	podDisruptionBudgetMetricFamilies			= []metrics.FamilyGenerator{metrics.FamilyGenerator{Name: "kube_poddisruptionbudget_created", Type: metrics.MetricTypeGauge, Help: "Unix creation timestamp", GenerateFunc: wrapPodDisruptionBudgetFunc(func(p *v1beta1.PodDisruptionBudget) metrics.Family {
		f := metrics.Family{}
		if !p.CreationTimestamp.IsZero() {
			f = append(f, &metrics.Metric{Name: "kube_poddisruptionbudget_created", Value: float64(p.CreationTimestamp.Unix())})
		}
		return f
	})}, metrics.FamilyGenerator{Name: "kube_poddisruptionbudget_status_current_healthy", Type: metrics.MetricTypeGauge, Help: "Current number of healthy pods", GenerateFunc: wrapPodDisruptionBudgetFunc(func(p *v1beta1.PodDisruptionBudget) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_poddisruptionbudget_status_current_healthy", Value: float64(p.Status.CurrentHealthy)}}
	})}, metrics.FamilyGenerator{Name: "kube_poddisruptionbudget_status_desired_healthy", Type: metrics.MetricTypeGauge, Help: "Minimum desired number of healthy pods", GenerateFunc: wrapPodDisruptionBudgetFunc(func(p *v1beta1.PodDisruptionBudget) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_poddisruptionbudget_status_desired_healthy", Value: float64(p.Status.DesiredHealthy)}}
	})}, metrics.FamilyGenerator{Name: "kube_poddisruptionbudget_status_pod_disruptions_allowed", Type: metrics.MetricTypeGauge, Help: "Number of pod disruptions that are currently allowed", GenerateFunc: wrapPodDisruptionBudgetFunc(func(p *v1beta1.PodDisruptionBudget) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_poddisruptionbudget_status_pod_disruptions_allowed", Value: float64(p.Status.PodDisruptionsAllowed)}}
	})}, metrics.FamilyGenerator{Name: "kube_poddisruptionbudget_status_expected_pods", Type: metrics.MetricTypeGauge, Help: "Total number of pods counted by this disruption budget", GenerateFunc: wrapPodDisruptionBudgetFunc(func(p *v1beta1.PodDisruptionBudget) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_poddisruptionbudget_status_expected_pods", Value: float64(p.Status.ExpectedPods)}}
	})}, metrics.FamilyGenerator{Name: "kube_poddisruptionbudget_status_observed_generation", Type: metrics.MetricTypeGauge, Help: "Most recent generation observed when updating this PDB status", GenerateFunc: wrapPodDisruptionBudgetFunc(func(p *v1beta1.PodDisruptionBudget) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_poddisruptionbudget_status_observed_generation", Value: float64(p.Status.ObservedGeneration)}}
	})}}
)

func wrapPodDisruptionBudgetFunc(f func(*v1beta1.PodDisruptionBudget) metrics.Family) func(interface{}) metrics.Family {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(obj interface{}) metrics.Family {
		podDisruptionBudget := obj.(*v1beta1.PodDisruptionBudget)
		metricFamily := f(podDisruptionBudget)
		for _, m := range metricFamily {
			m.LabelKeys = append(descPodDisruptionBudgetLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{podDisruptionBudget.Namespace, podDisruptionBudget.Name}, m.LabelValues...)
		}
		return metricFamily
	}
}
func createPodDisruptionBudgetListWatch(kubeClient clientset.Interface, ns string) cache.ListWatch {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cache.ListWatch{ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.PolicyV1beta1().PodDisruptionBudgets(ns).List(opts)
	}, WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
		return kubeClient.PolicyV1beta1().PodDisruptionBudgets(ns).Watch(opts)
	}}
}
