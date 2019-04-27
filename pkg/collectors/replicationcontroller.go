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
	descReplicationControllerLabelsDefaultLabels	= []string{"namespace", "replicationcontroller"}
	replicationControllerMetricFamilies		= []metrics.FamilyGenerator{metrics.FamilyGenerator{Name: "kube_replicationcontroller_created", Type: metrics.MetricTypeGauge, Help: "Unix creation timestamp", GenerateFunc: wrapReplicationControllerFunc(func(r *v1.ReplicationController) metrics.Family {
		f := metrics.Family{}
		if !r.CreationTimestamp.IsZero() {
			f = append(f, &metrics.Metric{Name: "kube_replicationcontroller_created", Value: float64(r.CreationTimestamp.Unix())})
		}
		return f
	})}, metrics.FamilyGenerator{Name: "kube_replicationcontroller_status_replicas", Type: metrics.MetricTypeGauge, Help: "The number of replicas per ReplicationController.", GenerateFunc: wrapReplicationControllerFunc(func(r *v1.ReplicationController) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_replicationcontroller_status_replicas", Value: float64(r.Status.Replicas)}}
	})}, metrics.FamilyGenerator{Name: "kube_replicationcontroller_status_fully_labeled_replicas", Type: metrics.MetricTypeGauge, Help: "The number of fully labeled replicas per ReplicationController.", GenerateFunc: wrapReplicationControllerFunc(func(r *v1.ReplicationController) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_replicationcontroller_status_fully_labeled_replicas", Value: float64(r.Status.FullyLabeledReplicas)}}
	})}, metrics.FamilyGenerator{Name: "kube_replicationcontroller_status_ready_replicas", Type: metrics.MetricTypeGauge, Help: "The number of ready replicas per ReplicationController.", GenerateFunc: wrapReplicationControllerFunc(func(r *v1.ReplicationController) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_replicationcontroller_status_ready_replicas", Value: float64(r.Status.ReadyReplicas)}}
	})}, metrics.FamilyGenerator{Name: "kube_replicationcontroller_status_available_replicas", Type: metrics.MetricTypeGauge, Help: "The number of available replicas per ReplicationController.", GenerateFunc: wrapReplicationControllerFunc(func(r *v1.ReplicationController) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_replicationcontroller_status_available_replicas", Value: float64(r.Status.AvailableReplicas)}}
	})}, metrics.FamilyGenerator{Name: "kube_replicationcontroller_status_observed_generation", Type: metrics.MetricTypeGauge, Help: "The generation observed by the ReplicationController controller.", GenerateFunc: wrapReplicationControllerFunc(func(r *v1.ReplicationController) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_replicationcontroller_status_observed_generation", Value: float64(r.Status.ObservedGeneration)}}
	})}, metrics.FamilyGenerator{Name: "kube_replicationcontroller_spec_replicas", Type: metrics.MetricTypeGauge, Help: "Number of desired pods for a ReplicationController.", GenerateFunc: wrapReplicationControllerFunc(func(r *v1.ReplicationController) metrics.Family {
		f := metrics.Family{}
		if r.Spec.Replicas != nil {
			f = append(f, &metrics.Metric{Name: "kube_replicationcontroller_spec_replicas", Value: float64(*r.Spec.Replicas)})
		}
		return f
	})}, metrics.FamilyGenerator{Name: "kube_replicationcontroller_metadata_generation", Type: metrics.MetricTypeGauge, Help: "Sequence number representing a specific generation of the desired state.", GenerateFunc: wrapReplicationControllerFunc(func(r *v1.ReplicationController) metrics.Family {
		return metrics.Family{&metrics.Metric{Name: "kube_replicationcontroller_metadata_generation", Value: float64(r.ObjectMeta.Generation)}}
	})}}
)

func wrapReplicationControllerFunc(f func(*v1.ReplicationController) metrics.Family) func(interface{}) metrics.Family {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(obj interface{}) metrics.Family {
		replicationController := obj.(*v1.ReplicationController)
		metricFamily := f(replicationController)
		for _, m := range metricFamily {
			m.LabelKeys = append(descReplicationControllerLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{replicationController.Namespace, replicationController.Name}, m.LabelValues...)
		}
		return metricFamily
	}
}
func createReplicationControllerListWatch(kubeClient clientset.Interface, ns string) cache.ListWatch {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cache.ListWatch{ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().ReplicationControllers(ns).List(opts)
	}, WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
		return kubeClient.CoreV1().ReplicationControllers(ns).Watch(opts)
	}}
}
