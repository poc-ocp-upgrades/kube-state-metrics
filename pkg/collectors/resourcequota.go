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
	descResourceQuotaLabelsDefaultLabels	= []string{"namespace", "resourcequota"}
	resourceQuotaMetricFamilies				= []metrics.FamilyGenerator{metrics.FamilyGenerator{Name: "kube_resourcequota_created", Type: metrics.MetricTypeGauge, Help: "Unix creation timestamp", GenerateFunc: wrapResourceQuotaFunc(func(r *v1.ResourceQuota) metrics.Family {
		f := metrics.Family{}
		if !r.CreationTimestamp.IsZero() {
			f = append(f, &metrics.Metric{Name: "kube_resourcequota_created", Value: float64(r.CreationTimestamp.Unix())})
		}
		return f
	})}, metrics.FamilyGenerator{Name: "kube_resourcequota", Type: metrics.MetricTypeGauge, Help: "Information about resource quota.", GenerateFunc: wrapResourceQuotaFunc(func(r *v1.ResourceQuota) metrics.Family {
		f := metrics.Family{}
		for res, qty := range r.Status.Hard {
			f = append(f, &metrics.Metric{LabelValues: []string{string(res), "hard"}, Value: float64(qty.MilliValue()) / 1000})
		}
		for res, qty := range r.Status.Used {
			f = append(f, &metrics.Metric{LabelValues: []string{string(res), "used"}, Value: float64(qty.MilliValue()) / 1000})
		}
		for _, m := range f {
			m.Name = "kube_resourcequota"
			m.LabelKeys = []string{"resource", "type"}
		}
		return f
	})}}
)

func wrapResourceQuotaFunc(f func(*v1.ResourceQuota) metrics.Family) func(interface{}) metrics.Family {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(obj interface{}) metrics.Family {
		resourceQuota := obj.(*v1.ResourceQuota)
		metricFamily := f(resourceQuota)
		for _, m := range metricFamily {
			m.LabelKeys = append(descResourceQuotaLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{resourceQuota.Namespace, resourceQuota.Name}, m.LabelValues...)
		}
		return metricFamily
	}
}
func createResourceQuotaListWatch(kubeClient clientset.Interface, ns string) cache.ListWatch {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cache.ListWatch{ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().ResourceQuotas(ns).List(opts)
	}, WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
		return kubeClient.CoreV1().ResourceQuotas(ns).Watch(opts)
	}}
}
