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
	descNamespaceLabelsName			= "kube_namespace_labels"
	descNamespaceLabelsHelp			= "Kubernetes labels converted to Prometheus labels."
	descNamespaceLabelsDefaultLabels	= []string{"namespace"}
	descNamespaceAnnotationsName		= "kube_namespace_annotations"
	descNamespaceAnnotationsHelp		= "Kubernetes annotations converted to Prometheus labels."
	descNamespaceAnnotationsDefaultLabels	= []string{"namespace"}
	namespaceMetricFamilies			= []metrics.FamilyGenerator{metrics.FamilyGenerator{Name: "kube_namespace_created", Type: metrics.MetricTypeGauge, Help: "Unix creation timestamp", GenerateFunc: wrapNamespaceFunc(func(n *v1.Namespace) metrics.Family {
		f := metrics.Family{}
		if !n.CreationTimestamp.IsZero() {
			f = append(f, &metrics.Metric{Name: "kube_namespace_created", Value: float64(n.CreationTimestamp.Unix())})
		}
		return f
	})}, metrics.FamilyGenerator{Name: descNamespaceLabelsName, Type: metrics.MetricTypeGauge, Help: descNamespaceLabelsHelp, GenerateFunc: wrapNamespaceFunc(func(n *v1.Namespace) metrics.Family {
		labelKeys, labelValues := kubeLabelsToPrometheusLabels(n.Labels)
		return metrics.Family{&metrics.Metric{Name: descNamespaceLabelsName, LabelKeys: labelKeys, LabelValues: labelValues, Value: 1}}
	})}, metrics.FamilyGenerator{Name: descNamespaceAnnotationsName, Type: metrics.MetricTypeGauge, Help: descNamespaceAnnotationsHelp, GenerateFunc: wrapNamespaceFunc(func(n *v1.Namespace) metrics.Family {
		annotationKeys, annotationValues := kubeAnnotationsToPrometheusAnnotations(n.Annotations)
		return metrics.Family{&metrics.Metric{Name: descNamespaceAnnotationsName, LabelKeys: annotationKeys, LabelValues: annotationValues, Value: 1}}
	})}, metrics.FamilyGenerator{Name: "kube_namespace_status_phase", Type: metrics.MetricTypeGauge, Help: "kubernetes namespace status phase.", GenerateFunc: wrapNamespaceFunc(func(n *v1.Namespace) metrics.Family {
		families := metrics.Family{&metrics.Metric{LabelValues: []string{string(v1.NamespaceActive)}, Value: boolFloat64(n.Status.Phase == v1.NamespaceActive)}, &metrics.Metric{LabelValues: []string{string(v1.NamespaceTerminating)}, Value: boolFloat64(n.Status.Phase == v1.NamespaceTerminating)}}
		for _, f := range families {
			f.Name = "kube_namespace_status_phase"
			f.LabelKeys = []string{"phase"}
		}
		return families
	})}}
)

func wrapNamespaceFunc(f func(*v1.Namespace) metrics.Family) func(interface{}) metrics.Family {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(obj interface{}) metrics.Family {
		namespace := obj.(*v1.Namespace)
		metricFamily := f(namespace)
		for _, m := range metricFamily {
			m.LabelKeys = append(descNamespaceLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{namespace.Name}, m.LabelValues...)
		}
		return metricFamily
	}
}
func createNamespaceListWatch(kubeClient clientset.Interface, ns string) cache.ListWatch {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cache.ListWatch{ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().Namespaces().List(opts)
	}, WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
		return kubeClient.CoreV1().Namespaces().Watch(opts)
	}}
}
