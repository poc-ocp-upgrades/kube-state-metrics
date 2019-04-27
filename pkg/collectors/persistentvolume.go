package collectors

import (
	"k8s.io/kube-state-metrics/pkg/metrics"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var (
	descPersistentVolumeLabelsName		= "kube_persistentvolume_labels"
	descPersistentVolumeLabelsHelp		= "Kubernetes labels converted to Prometheus labels."
	descPersistentVolumeLabelsDefaultLabels	= []string{"persistentvolume"}
	persistentVolumeMetricFamilies		= []metrics.FamilyGenerator{metrics.FamilyGenerator{Name: descPersistentVolumeLabelsName, Type: metrics.MetricTypeGauge, Help: descPersistentVolumeLabelsHelp, GenerateFunc: wrapPersistentVolumeFunc(func(p *v1.PersistentVolume) metrics.Family {
		labelKeys, labelValues := kubeLabelsToPrometheusLabels(p.Labels)
		return metrics.Family{&metrics.Metric{Name: descPersistentVolumeLabelsName, LabelKeys: labelKeys, LabelValues: labelValues, Value: 1}}
	})}, metrics.FamilyGenerator{Name: "kube_persistentvolume_status_phase", Type: metrics.MetricTypeGauge, Help: "The phase indicates if a volume is available, bound to a claim, or released by a claim.", GenerateFunc: wrapPersistentVolumeFunc(func(p *v1.PersistentVolume) metrics.Family {
		f := metrics.Family{}
		if p := p.Status.Phase; p != "" {
			f = append(f, &metrics.Metric{LabelValues: []string{string(v1.VolumePending)}, Value: boolFloat64(p == v1.VolumePending)}, &metrics.Metric{LabelValues: []string{string(v1.VolumeAvailable)}, Value: boolFloat64(p == v1.VolumeAvailable)}, &metrics.Metric{LabelValues: []string{string(v1.VolumeBound)}, Value: boolFloat64(p == v1.VolumeBound)}, &metrics.Metric{LabelValues: []string{string(v1.VolumeReleased)}, Value: boolFloat64(p == v1.VolumeReleased)}, &metrics.Metric{LabelValues: []string{string(v1.VolumeFailed)}, Value: boolFloat64(p == v1.VolumeFailed)})
		}
		for _, m := range f {
			m.Name = "kube_persistentvolume_status_phase"
			m.LabelKeys = []string{"phase"}
		}
		return f
	})}, metrics.FamilyGenerator{Name: "kube_persistentvolume_info", Type: metrics.MetricTypeGauge, Help: "Information about persistentvolume.", GenerateFunc: wrapPersistentVolumeFunc(func(p *v1.PersistentVolume) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_persistentvolume_info", LabelKeys: []string{"storageclass"}, LabelValues: []string{p.Spec.StorageClassName}, Value: 1}}
	})}}
)

func wrapPersistentVolumeFunc(f func(*v1.PersistentVolume) metrics.Family) func(interface{}) metrics.Family {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(obj interface{}) metrics.Family {
		persistentVolume := obj.(*v1.PersistentVolume)
		metricFamily := f(persistentVolume)
		for _, m := range metricFamily {
			m.LabelKeys = append(descPersistentVolumeLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{persistentVolume.Name}, m.LabelValues...)
		}
		return metricFamily
	}
}
func createPersistentVolumeListWatch(kubeClient clientset.Interface, ns string) cache.ListWatch {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cache.ListWatch{ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().PersistentVolumes().List(opts)
	}, WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
		return kubeClient.CoreV1().PersistentVolumes().Watch(opts)
	}}
}
