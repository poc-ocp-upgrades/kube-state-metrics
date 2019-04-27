package collectors

import (
	"testing"
)

func TestSortLabels(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	in := `kube_pod_container_info{container_id="docker://cd456",image="k8s.gcr.io/hyperkube2",container="container2",image_id="docker://sha256:bbb",namespace="ns2",pod="pod2"} 1
kube_pod_container_info{namespace="ns2",container="container3",container_id="docker://ef789",image="k8s.gcr.io/hyperkube3",image_id="docker://sha256:ccc",pod="pod2"} 1`
	want := `kube_pod_container_info{container="container2",container_id="docker://cd456",image="k8s.gcr.io/hyperkube2",image_id="docker://sha256:bbb",namespace="ns2",pod="pod2"} 1
kube_pod_container_info{container="container3",container_id="docker://ef789",image="k8s.gcr.io/hyperkube3",image_id="docker://sha256:ccc",namespace="ns2",pod="pod2"} 1`
	out := sortLabels(in)
	if want != out {
		t.Fatalf("expected:\n%v\nbut got:\n%v", want, out)
	}
}
