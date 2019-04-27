package collectors

import (
	"k8s.io/kube-state-metrics/pkg/metrics"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var (
	descDaemonSetLabelsName			= "kube_daemonset_labels"
	descDaemonSetLabelsHelp			= "Kubernetes labels converted to Prometheus labels."
	descDaemonSetLabelsDefaultLabels	= []string{"namespace", "daemonset"}
	daemonSetMetricFamilies			= []metrics.FamilyGenerator{metrics.FamilyGenerator{Name: "kube_daemonset_created", Type: metrics.MetricTypeGauge, Help: "Unix creation timestamp", GenerateFunc: wrapDaemonSetFunc(func(d *v1beta1.DaemonSet) metrics.Family {
		f := metrics.Family{}
		if !d.CreationTimestamp.IsZero() {
			f = append(f, &metrics.Metric{Name: "kube_daemonset_created", LabelKeys: []string{}, LabelValues: []string{}, Value: float64(d.CreationTimestamp.Unix())})
		}
		return f
	})}, metrics.FamilyGenerator{Name: "kube_daemonset_status_current_number_scheduled", Type: metrics.MetricTypeGauge, Help: "The number of nodes running at least one daemon pod and are supposed to.", GenerateFunc: wrapDaemonSetFunc(func(d *v1beta1.DaemonSet) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_daemonset_status_current_number_scheduled", LabelKeys: []string{}, LabelValues: []string{}, Value: float64(d.Status.CurrentNumberScheduled)}}
	})}, metrics.FamilyGenerator{Name: "kube_daemonset_status_desired_number_scheduled", Type: metrics.MetricTypeGauge, Help: "The number of nodes that should be running the daemon pod.", GenerateFunc: wrapDaemonSetFunc(func(d *v1beta1.DaemonSet) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_daemonset_status_desired_number_scheduled", LabelKeys: []string{}, LabelValues: []string{}, Value: float64(d.Status.DesiredNumberScheduled)}}
	})}, metrics.FamilyGenerator{Name: "kube_daemonset_status_number_available", Type: metrics.MetricTypeGauge, Help: "The number of nodes that should be running the daemon pod and have one or more of the daemon pod running and available", GenerateFunc: wrapDaemonSetFunc(func(d *v1beta1.DaemonSet) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_daemonset_status_number_available", LabelKeys: []string{}, LabelValues: []string{}, Value: float64(d.Status.NumberAvailable)}}
	})}, metrics.FamilyGenerator{Name: "kube_daemonset_status_number_misscheduled", Type: metrics.MetricTypeGauge, Help: "The number of nodes running a daemon pod but are not supposed to.", GenerateFunc: wrapDaemonSetFunc(func(d *v1beta1.DaemonSet) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_daemonset_status_number_misscheduled", LabelKeys: []string{}, LabelValues: []string{}, Value: float64(d.Status.NumberMisscheduled)}}
	})}, metrics.FamilyGenerator{Name: "kube_daemonset_status_number_ready", Type: metrics.MetricTypeGauge, Help: "The number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready.", GenerateFunc: wrapDaemonSetFunc(func(d *v1beta1.DaemonSet) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_daemonset_status_number_ready", LabelKeys: []string{}, LabelValues: []string{}, Value: float64(d.Status.NumberReady)}}
	})}, metrics.FamilyGenerator{Name: "kube_daemonset_status_number_unavailable", Type: metrics.MetricTypeGauge, Help: "The number of nodes that should be running the daemon pod and have none of the daemon pod running and available", GenerateFunc: wrapDaemonSetFunc(func(d *v1beta1.DaemonSet) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_daemonset_status_number_unavailable", LabelKeys: []string{}, LabelValues: []string{}, Value: float64(d.Status.NumberUnavailable)}}
	})}, metrics.FamilyGenerator{Name: "kube_daemonset_updated_number_scheduled", Type: metrics.MetricTypeGauge, Help: "The total number of nodes that are running updated daemon pod", GenerateFunc: wrapDaemonSetFunc(func(d *v1beta1.DaemonSet) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_daemonset_updated_number_scheduled", Value: float64(d.Status.UpdatedNumberScheduled)}}
	})}, metrics.FamilyGenerator{Name: "kube_daemonset_metadata_generation", Type: metrics.MetricTypeGauge, Help: "Sequence number representing a specific generation of the desired state.", GenerateFunc: wrapDaemonSetFunc(func(d *v1beta1.DaemonSet) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_daemonset_metadata_generation", LabelKeys: []string{}, LabelValues: []string{}, Value: float64(d.ObjectMeta.Generation)}}
	})}, metrics.FamilyGenerator{Name: descDaemonSetLabelsName, Type: metrics.MetricTypeGauge, Help: descDaemonSetLabelsHelp, GenerateFunc: wrapDaemonSetFunc(func(d *v1beta1.DaemonSet) metrics.Family {
		labelKeys, labelValues := kubeLabelsToPrometheusLabels(d.ObjectMeta.Labels)
		return metrics.Family{&metrics.Metric{Name: descDaemonSetLabelsName, LabelKeys: labelKeys, LabelValues: labelValues, Value: 1}}
	})}}
)

func wrapDaemonSetFunc(f func(*v1beta1.DaemonSet) metrics.Family) func(interface{}) metrics.Family {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(obj interface{}) metrics.Family {
		daemonSet := obj.(*v1beta1.DaemonSet)
		metricFamily := f(daemonSet)
		for _, m := range metricFamily {
			m.LabelKeys = append(descDaemonSetLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{daemonSet.Namespace, daemonSet.Name}, m.LabelValues...)
		}
		return metricFamily
	}
}
func createDaemonSetListWatch(kubeClient clientset.Interface, ns string) cache.ListWatch {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cache.ListWatch{ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.ExtensionsV1beta1().DaemonSets(ns).List(opts)
	}, WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
		return kubeClient.ExtensionsV1beta1().DaemonSets(ns).Watch(opts)
	}}
}
