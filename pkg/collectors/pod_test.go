package collectors

import (
	"testing"
	"time"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/util/node"
)

func TestPodCollector(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var test = true
	startTime := 1501569018
	metav1StartTime := metav1.Unix(int64(startTime), 0)
	const metadata = ""
	cases := []generateMetricsTestCase{{Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container1", Image: "k8s.gcr.io/hyperkube1", ImageID: "docker://sha256:aaa", ContainerID: "docker://ab123"}}}}, Want: `kube_pod_container_info{container="container1",container_id="docker://ab123",image="k8s.gcr.io/hyperkube1",image_id="docker://sha256:aaa",namespace="ns1",pod="pod1"} 1`, MetricNames: []string{"kube_pod_container_info"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "ns2"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container2", Image: "k8s.gcr.io/hyperkube2", ImageID: "docker://sha256:bbb", ContainerID: "docker://cd456"}, v1.ContainerStatus{Name: "container3", Image: "k8s.gcr.io/hyperkube3", ImageID: "docker://sha256:ccc", ContainerID: "docker://ef789"}}}}, Want: `kube_pod_container_info{container="container2",container_id="docker://cd456",image="k8s.gcr.io/hyperkube2",image_id="docker://sha256:bbb",namespace="ns2",pod="pod2"} 1
		kube_pod_container_info{container="container3",container_id="docker://ef789",image="k8s.gcr.io/hyperkube3",image_id="docker://sha256:ccc",namespace="ns2",pod="pod2"} 1`, MetricNames: []string{"kube_pod_container_info"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container1", Ready: true}}}}, Want: `kube_pod_container_status_ready{container="container1",namespace="ns1",pod="pod1"} 1`, MetricNames: []string{"kube_pod_container_status_ready"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "ns2"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container2", Ready: true}, v1.ContainerStatus{Name: "container3", Ready: false}}}}, Want: metadata + `
				kube_pod_container_status_ready{container="container2",namespace="ns2",pod="pod2"} 1
				kube_pod_container_status_ready{container="container3",namespace="ns2",pod="pod2"} 0
				`, MetricNames: []string{"kube_pod_container_status_ready"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container1", RestartCount: 0}}}}, Want: `kube_pod_container_status_restarts_total{container="container1",namespace="ns1",pod="pod1"} 0`, MetricNames: []string{"kube_pod_container_status_restarts_total"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "ns2"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container2", RestartCount: 0}, v1.ContainerStatus{Name: "container3", RestartCount: 1}}}}, Want: metadata + `
				kube_pod_container_status_restarts_total{container="container2",namespace="ns2",pod="pod2"} 0
				kube_pod_container_status_restarts_total{container="container3",namespace="ns2",pod="pod2"} 1
				`, MetricNames: []string{"kube_pod_container_status_restarts_total"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container1", State: v1.ContainerState{Running: &v1.ContainerStateRunning{}}}}}}, Want: `
				kube_pod_container_status_running{container="container1",namespace="ns1",pod="pod1"} 1
				kube_pod_container_status_terminated_reason{container="container1",namespace="ns1",pod="pod1",reason="Completed"} 0
                kube_pod_container_status_terminated_reason{container="container1",namespace="ns1",pod="pod1",reason="ContainerCannotRun"} 0
				kube_pod_container_status_terminated_reason{container="container1",namespace="ns1",pod="pod1",reason="Error"} 0
				kube_pod_container_status_terminated_reason{container="container1",namespace="ns1",pod="pod1",reason="OOMKilled"} 0
                kube_pod_container_status_terminated{container="container1",namespace="ns1",pod="pod1"} 0
				kube_pod_container_status_waiting{container="container1",namespace="ns1",pod="pod1"} 0
				kube_pod_container_status_waiting_reason{container="container1",namespace="ns1",pod="pod1",reason="ContainerCreating"} 0
				kube_pod_container_status_waiting_reason{container="container1",namespace="ns1",pod="pod1",reason="ImagePullBackOff"} 0
				kube_pod_container_status_waiting_reason{container="container1",namespace="ns1",pod="pod1",reason="CrashLoopBackOff"} 0
				kube_pod_container_status_waiting_reason{container="container1",namespace="ns1",pod="pod1",reason="ErrImagePull"} 0
				kube_pod_container_status_waiting_reason{container="container1",namespace="ns1",pod="pod1",reason="CreateContainerConfigError"} 0

`, MetricNames: []string{"kube_pod_container_status_running", "kube_pod_container_status_waiting", "kube_pod_container_status_waiting_reason", "kube_pod_container_status_terminated", "kube_pod_container_status_terminated_reason"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "ns2"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container2", State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{Reason: "OOMKilled"}}}, v1.ContainerStatus{Name: "container3", State: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: "ContainerCreating"}}}}}}, Want: `
				kube_pod_container_status_running{container="container2",namespace="ns2",pod="pod2"} 0
                kube_pod_container_status_running{container="container3",namespace="ns2",pod="pod2"} 0
				kube_pod_container_status_terminated{container="container2",namespace="ns2",pod="pod2"} 1
				kube_pod_container_status_terminated_reason{container="container2",namespace="ns2",pod="pod2",reason="Completed"} 0
				kube_pod_container_status_terminated_reason{container="container2",namespace="ns2",pod="pod2",reason="ContainerCannotRun"} 0
				kube_pod_container_status_terminated_reason{container="container2",namespace="ns2",pod="pod2",reason="Error"} 0
				kube_pod_container_status_terminated_reason{container="container2",namespace="ns2",pod="pod2",reason="OOMKilled"} 1
                kube_pod_container_status_terminated_reason{container="container3",namespace="ns2",pod="pod2",reason="Completed"} 0
                kube_pod_container_status_terminated_reason{container="container3",namespace="ns2",pod="pod2",reason="ContainerCannotRun"} 0
                kube_pod_container_status_terminated_reason{container="container3",namespace="ns2",pod="pod2",reason="Error"} 0
                kube_pod_container_status_terminated_reason{container="container3",namespace="ns2",pod="pod2",reason="OOMKilled"} 0
				kube_pod_container_status_waiting{container="container2",namespace="ns2",pod="pod2"} 0
                kube_pod_container_status_waiting{container="container3",namespace="ns2",pod="pod2"} 1
                kube_pod_container_status_terminated{container="container3",namespace="ns2",pod="pod2"} 0
				kube_pod_container_status_waiting_reason{container="container2",namespace="ns2",pod="pod2",reason="ContainerCreating"} 0
				kube_pod_container_status_waiting_reason{container="container2",namespace="ns2",pod="pod2",reason="ImagePullBackOff"} 0
				kube_pod_container_status_waiting_reason{container="container2",namespace="ns2",pod="pod2",reason="CrashLoopBackOff"} 0
				kube_pod_container_status_waiting_reason{container="container2",namespace="ns2",pod="pod2",reason="ErrImagePull"} 0
				kube_pod_container_status_waiting_reason{container="container2",namespace="ns2",pod="pod2",reason="CreateContainerConfigError"} 0
                kube_pod_container_status_waiting_reason{container="container3",namespace="ns2",pod="pod2",reason="ContainerCreating"} 1
                kube_pod_container_status_waiting_reason{container="container3",namespace="ns2",pod="pod2",reason="CrashLoopBackOff"} 0
                kube_pod_container_status_waiting_reason{container="container3",namespace="ns2",pod="pod2",reason="ErrImagePull"} 0
				kube_pod_container_status_waiting_reason{container="container3",namespace="ns2",pod="pod2",reason="ImagePullBackOff"} 0
				kube_pod_container_status_waiting_reason{container="container3",namespace="ns2",pod="pod2",reason="CreateContainerConfigError"} 0

`, MetricNames: []string{"kube_pod_container_status_running", "kube_pod_container_status_waiting", "kube_pod_container_status_waiting_reason", "kube_pod_container_status_terminated", "kube_pod_container_status_terminated_reason"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod3", Namespace: "ns3"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container4", State: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: "CrashLoopBackOff"}}}}}}, Want: `
				kube_pod_container_status_running{container="container4",namespace="ns3",pod="pod3"} 0
				kube_pod_container_status_terminated{container="container4",namespace="ns3",pod="pod3"} 0
kube_pod_container_status_terminated_reason{container="container4",namespace="ns3",pod="pod3",reason="Completed"} 0
				kube_pod_container_status_terminated_reason{container="container4",namespace="ns3",pod="pod3",reason="ContainerCannotRun"} 0
				kube_pod_container_status_terminated_reason{container="container4",namespace="ns3",pod="pod3",reason="Error"} 0
				kube_pod_container_status_terminated_reason{container="container4",namespace="ns3",pod="pod3",reason="OOMKilled"} 0
				kube_pod_container_status_waiting{container="container4",namespace="ns3",pod="pod3"} 1
kube_pod_container_status_waiting_reason{container="container4",namespace="ns3",pod="pod3",reason="ContainerCreating"} 0
				kube_pod_container_status_waiting_reason{container="container4",namespace="ns3",pod="pod3",reason="ImagePullBackOff"} 0
				kube_pod_container_status_waiting_reason{container="container4",namespace="ns3",pod="pod3",reason="CrashLoopBackOff"} 1
				kube_pod_container_status_waiting_reason{container="container4",namespace="ns3",pod="pod3",reason="ErrImagePull"} 0
				kube_pod_container_status_waiting_reason{container="container4",namespace="ns3",pod="pod3",reason="CreateContainerConfigError"} 0
kube_pod_container_status_last_terminated_reason{container="container4",namespace="ns3",pod="pod3",reason="Completed"} 0
				kube_pod_container_status_last_terminated_reason{container="container4",namespace="ns3",pod="pod3",reason="ContainerCannotRun"} 0
				kube_pod_container_status_last_terminated_reason{container="container4",namespace="ns3",pod="pod3",reason="Error"} 0
				kube_pod_container_status_last_terminated_reason{container="container4",namespace="ns3",pod="pod3",reason="OOMKilled"} 0
`, MetricNames: []string{"kube_pod_container_status_running", "kube_pod_container_status_terminated", "kube_pod_container_status_terminated_reason", "kube_pod_container_status_terminated_reason", "kube_pod_container_status_terminated_reason", "kube_pod_container_status_terminated_reason", "kube_pod_container_status_waiting", "kube_pod_container_status_waiting_reason", "kube_pod_container_status_waiting_reason", "kube_pod_container_status_waiting_reason", "kube_pod_container_status_waiting_reason", "kube_pod_container_status_waiting_reason", "kube_pod_container_status_last_terminated_reason", "kube_pod_container_status_last_terminated_reason", "kube_pod_container_status_last_terminated_reason", "kube_pod_container_status_last_terminated_reason"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod6", Namespace: "ns6"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container7", State: v1.ContainerState{Running: &v1.ContainerStateRunning{}}, LastTerminationState: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{Reason: "OOMKilled"}}}}}}, Want: `
				kube_pod_container_status_running{container="container7",namespace="ns6",pod="pod6"} 1
				kube_pod_container_status_terminated{container="container7",namespace="ns6",pod="pod6"} 0
kube_pod_container_status_terminated_reason{container="container7",namespace="ns6",pod="pod6",reason="Completed"} 0
				kube_pod_container_status_terminated_reason{container="container7",namespace="ns6",pod="pod6",reason="ContainerCannotRun"} 0
				kube_pod_container_status_terminated_reason{container="container7",namespace="ns6",pod="pod6",reason="Error"} 0
				kube_pod_container_status_terminated_reason{container="container7",namespace="ns6",pod="pod6",reason="OOMKilled"} 0
				kube_pod_container_status_waiting{container="container7",namespace="ns6",pod="pod6"} 0
kube_pod_container_status_waiting_reason{container="container7",namespace="ns6",pod="pod6",reason="ContainerCreating"} 0
				kube_pod_container_status_waiting_reason{container="container7",namespace="ns6",pod="pod6",reason="ImagePullBackOff"} 0
				kube_pod_container_status_waiting_reason{container="container7",namespace="ns6",pod="pod6",reason="CrashLoopBackOff"} 0
				kube_pod_container_status_waiting_reason{container="container7",namespace="ns6",pod="pod6",reason="ErrImagePull"} 0
				kube_pod_container_status_waiting_reason{container="container7",namespace="ns6",pod="pod6",reason="CreateContainerConfigError"} 0
kube_pod_container_status_last_terminated_reason{container="container7",namespace="ns6",pod="pod6",reason="Completed"} 0
				kube_pod_container_status_last_terminated_reason{container="container7",namespace="ns6",pod="pod6",reason="ContainerCannotRun"} 0
				kube_pod_container_status_last_terminated_reason{container="container7",namespace="ns6",pod="pod6",reason="Error"} 0
				kube_pod_container_status_last_terminated_reason{container="container7",namespace="ns6",pod="pod6",reason="OOMKilled"} 1
`, MetricNames: []string{"kube_pod_container_status_last_terminated_reason", "kube_pod_container_status_running", "kube_pod_container_status_terminated", "kube_pod_container_status_terminated_reason", "kube_pod_container_status_waiting", "kube_pod_container_status_waiting_reason"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod4", Namespace: "ns4"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container5", State: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: "ImagePullBackOff"}}}}}}, Want: `
				kube_pod_container_status_running{container="container5",namespace="ns4",pod="pod4"} 0
				kube_pod_container_status_terminated{container="container5",namespace="ns4",pod="pod4"} 0
				kube_pod_container_status_terminated_reason{container="container5",namespace="ns4",pod="pod4",reason="Completed"} 0
				kube_pod_container_status_terminated_reason{container="container5",namespace="ns4",pod="pod4",reason="ContainerCannotRun"} 0
				kube_pod_container_status_terminated_reason{container="container5",namespace="ns4",pod="pod4",reason="Error"} 0
				kube_pod_container_status_terminated_reason{container="container5",namespace="ns4",pod="pod4",reason="OOMKilled"} 0
				kube_pod_container_status_waiting{container="container5",namespace="ns4",pod="pod4"} 1
				kube_pod_container_status_waiting_reason{container="container5",namespace="ns4",pod="pod4",reason="ContainerCreating"} 0
				kube_pod_container_status_waiting_reason{container="container5",namespace="ns4",pod="pod4",reason="ImagePullBackOff"} 1
				kube_pod_container_status_waiting_reason{container="container5",namespace="ns4",pod="pod4",reason="CrashLoopBackOff"} 0
				kube_pod_container_status_waiting_reason{container="container5",namespace="ns4",pod="pod4",reason="ErrImagePull"} 0
				kube_pod_container_status_waiting_reason{container="container5",namespace="ns4",pod="pod4",reason="CreateContainerConfigError"} 0
`, MetricNames: []string{"kube_pod_container_status_running", "kube_pod_container_status_waiting", "kube_pod_container_status_waiting_reason", "kube_pod_container_status_terminated", "kube_pod_container_status_terminated_reason"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod5", Namespace: "ns5"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container6", State: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: "ErrImagePull"}}}}}}, Want: `
				kube_pod_container_status_running{container="container6",namespace="ns5",pod="pod5"} 0
				kube_pod_container_status_terminated{container="container6",namespace="ns5",pod="pod5"} 0
				kube_pod_container_status_terminated_reason{container="container6",namespace="ns5",pod="pod5",reason="Completed"} 0
				kube_pod_container_status_terminated_reason{container="container6",namespace="ns5",pod="pod5",reason="ContainerCannotRun"} 0
				kube_pod_container_status_terminated_reason{container="container6",namespace="ns5",pod="pod5",reason="Error"} 0
				kube_pod_container_status_terminated_reason{container="container6",namespace="ns5",pod="pod5",reason="OOMKilled"} 0
				kube_pod_container_status_waiting{container="container6",namespace="ns5",pod="pod5"} 1
				kube_pod_container_status_waiting_reason{container="container6",namespace="ns5",pod="pod5",reason="ContainerCreating"} 0
				kube_pod_container_status_waiting_reason{container="container6",namespace="ns5",pod="pod5",reason="ImagePullBackOff"} 0
				kube_pod_container_status_waiting_reason{container="container6",namespace="ns5",pod="pod5",reason="CrashLoopBackOff"} 0
				kube_pod_container_status_waiting_reason{container="container6",namespace="ns5",pod="pod5",reason="ErrImagePull"} 1
				kube_pod_container_status_waiting_reason{container="container6",namespace="ns5",pod="pod5",reason="CreateContainerConfigError"} 0
`, MetricNames: []string{"kube_pod_container_status_running", "kube_pod_container_status_waiting", "kube_pod_container_status_waiting_reason", "kube_pod_container_status_terminated", "kube_pod_container_status_terminated_reason"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod7", Namespace: "ns7"}, Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container8", State: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: "CreateContainerConfigError"}}}}}}, Want: `
					kube_pod_container_status_running{container="container8",namespace="ns7",pod="pod7"} 0
					kube_pod_container_status_terminated{container="container8",namespace="ns7",pod="pod7"} 0
					kube_pod_container_status_terminated_reason{container="container8",namespace="ns7",pod="pod7",reason="Completed"} 0
					kube_pod_container_status_terminated_reason{container="container8",namespace="ns7",pod="pod7",reason="ContainerCannotRun"} 0
					kube_pod_container_status_terminated_reason{container="container8",namespace="ns7",pod="pod7",reason="Error"} 0
					kube_pod_container_status_terminated_reason{container="container8",namespace="ns7",pod="pod7",reason="OOMKilled"} 0
					kube_pod_container_status_waiting{container="container8",namespace="ns7",pod="pod7"} 1
					kube_pod_container_status_waiting_reason{container="container8",namespace="ns7",pod="pod7",reason="ContainerCreating"} 0
					kube_pod_container_status_waiting_reason{container="container8",namespace="ns7",pod="pod7",reason="ImagePullBackOff"} 0
					kube_pod_container_status_waiting_reason{container="container8",namespace="ns7",pod="pod7",reason="CrashLoopBackOff"} 0
					kube_pod_container_status_waiting_reason{container="container8",namespace="ns7",pod="pod7",reason="ErrImagePull"} 0
					kube_pod_container_status_waiting_reason{container="container8",namespace="ns7",pod="pod7",reason="CreateContainerConfigError"} 1
`, MetricNames: []string{"kube_pod_container_status_running", "kube_pod_container_status_terminated", "kube_pod_container_status_terminated_reason", "kube_pod_container_status_waiting", "kube_pod_container_status_waiting_reason"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)}, Namespace: "ns1", UID: "abc-123-xxx"}, Spec: v1.PodSpec{NodeName: "node1"}, Status: v1.PodStatus{HostIP: "1.1.1.1", PodIP: "1.2.3.4", StartTime: &metav1StartTime}}, Want: `
				kube_pod_created{namespace="ns1",pod="pod1"} 1.5e+09
				kube_pod_info{created_by_kind="<none>",created_by_name="<none>",host_ip="1.1.1.1",namespace="ns1",node="node1",pod="pod1",pod_ip="1.2.3.4",uid="abc-123-xxx"} 1
				kube_pod_start_time{namespace="ns1",pod="pod1"} 1.501569018e+09
				kube_pod_owner{namespace="ns1",owner_is_controller="<none>",owner_kind="<none>",owner_name="<none>",pod="pod1"} 1
`, MetricNames: []string{"kube_pod_created", "kube_pod_info", "kube_pod_start_time", "kube_pod_completion_time", "kube_pod_owner"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "ns2", UID: "abc-456-xxx", OwnerReferences: []metav1.OwnerReference{{Kind: "ReplicaSet", Name: "rs-name", Controller: &test}}}, Spec: v1.PodSpec{NodeName: "node2"}, Status: v1.PodStatus{HostIP: "1.1.1.1", PodIP: "2.3.4.5", ContainerStatuses: []v1.ContainerStatus{v1.ContainerStatus{Name: "container2_1", Image: "k8s.gcr.io/hyperkube2", ImageID: "docker://sha256:bbb", ContainerID: "docker://cd456", State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{FinishedAt: metav1.Time{Time: time.Unix(1501777018, 0)}}}}, v1.ContainerStatus{Name: "container2_2", Image: "k8s.gcr.io/hyperkube2", ImageID: "docker://sha256:bbb", ContainerID: "docker://cd456", State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{FinishedAt: metav1.Time{Time: time.Unix(1501888018, 0)}}}}, v1.ContainerStatus{Name: "container2_3", Image: "k8s.gcr.io/hyperkube2", ImageID: "docker://sha256:bbb", ContainerID: "docker://cd456", State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{FinishedAt: metav1.Time{Time: time.Unix(1501666018, 0)}}}}}}}, Want: metadata + `
				kube_pod_info{created_by_kind="ReplicaSet",created_by_name="rs-name",host_ip="1.1.1.1",namespace="ns2",node="node2",pod="pod2",pod_ip="2.3.4.5",uid="abc-456-xxx"} 1
				kube_pod_completion_time{namespace="ns2",pod="pod2"} 1.501888018e+09
				kube_pod_owner{namespace="ns2",owner_is_controller="true",owner_kind="ReplicaSet",owner_name="rs-name",pod="pod2"} 1
				`, MetricNames: []string{"kube_pod_created", "kube_pod_info", "kube_pod_start_time", "kube_pod_completion_time", "kube_pod_owner"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"}, Status: v1.PodStatus{Phase: v1.PodRunning}}, Want: `
				kube_pod_status_phase{namespace="ns1",phase="Failed",pod="pod1"} 0
				kube_pod_status_phase{namespace="ns1",phase="Pending",pod="pod1"} 0
				kube_pod_status_phase{namespace="ns1",phase="Running",pod="pod1"} 1
				kube_pod_status_phase{namespace="ns1",phase="Succeeded",pod="pod1"} 0
				kube_pod_status_phase{namespace="ns1",phase="Unknown",pod="pod1"} 0
`, MetricNames: []string{"kube_pod_status_phase"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "ns2"}, Status: v1.PodStatus{Phase: v1.PodPending}}, Want: `
				kube_pod_status_phase{namespace="ns2",phase="Failed",pod="pod2"} 0
				kube_pod_status_phase{namespace="ns2",phase="Pending",pod="pod2"} 1
				kube_pod_status_phase{namespace="ns2",phase="Running",pod="pod2"} 0
				kube_pod_status_phase{namespace="ns2",phase="Succeeded",pod="pod2"} 0
				kube_pod_status_phase{namespace="ns2",phase="Unknown",pod="pod2"} 0
`, MetricNames: []string{"kube_pod_status_phase"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod3", Namespace: "ns3"}, Status: v1.PodStatus{Phase: v1.PodUnknown}}, Want: `
				kube_pod_status_phase{namespace="ns3",phase="Failed",pod="pod3"} 0
				kube_pod_status_phase{namespace="ns3",phase="Pending",pod="pod3"} 0
				kube_pod_status_phase{namespace="ns3",phase="Running",pod="pod3"} 0
				kube_pod_status_phase{namespace="ns3",phase="Succeeded",pod="pod3"} 0
				kube_pod_status_phase{namespace="ns3",phase="Unknown",pod="pod3"} 1
`, MetricNames: []string{"kube_pod_status_phase"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod4", Namespace: "ns4", DeletionTimestamp: &metav1.Time{}}, Status: v1.PodStatus{Phase: v1.PodRunning, Reason: node.NodeUnreachablePodReason}}, Want: `
				kube_pod_status_phase{namespace="ns4",phase="Failed",pod="pod4"} 0
				kube_pod_status_phase{namespace="ns4",phase="Pending",pod="pod4"} 0
				kube_pod_status_phase{namespace="ns4",phase="Running",pod="pod4"} 0
				kube_pod_status_phase{namespace="ns4",phase="Succeeded",pod="pod4"} 0
				kube_pod_status_phase{namespace="ns4",phase="Unknown",pod="pod4"} 1
`, MetricNames: []string{"kube_pod_status_phase"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"}, Status: v1.PodStatus{Conditions: []v1.PodCondition{v1.PodCondition{Type: v1.PodReady, Status: v1.ConditionTrue}}}}, Want: metadata + `
				kube_pod_status_ready{condition="false",namespace="ns1",pod="pod1"} 0
				kube_pod_status_ready{condition="true",namespace="ns1",pod="pod1"} 1
				kube_pod_status_ready{condition="unknown",namespace="ns1",pod="pod1"} 0
			`, MetricNames: []string{"kube_pod_status_ready"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "ns2"}, Status: v1.PodStatus{Conditions: []v1.PodCondition{v1.PodCondition{Type: v1.PodReady, Status: v1.ConditionFalse}}}}, Want: metadata + `
				kube_pod_status_ready{condition="false",namespace="ns2",pod="pod2"} 1
				kube_pod_status_ready{condition="true",namespace="ns2",pod="pod2"} 0
				kube_pod_status_ready{condition="unknown",namespace="ns2",pod="pod2"} 0
			`, MetricNames: []string{"kube_pod_status_ready"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"}, Status: v1.PodStatus{Conditions: []v1.PodCondition{v1.PodCondition{Type: v1.PodScheduled, Status: v1.ConditionTrue, LastTransitionTime: metav1.Time{Time: time.Unix(1501666018, 0)}}}}}, Want: metadata + `
				kube_pod_status_scheduled_time{namespace="ns1",pod="pod1"} 1.501666018e+09
				kube_pod_status_scheduled{condition="false",namespace="ns1",pod="pod1"} 0
				kube_pod_status_scheduled{condition="true",namespace="ns1",pod="pod1"} 1
				kube_pod_status_scheduled{condition="unknown",namespace="ns1",pod="pod1"} 0
			`, MetricNames: []string{"kube_pod_status_scheduled", "kube_pod_status_scheduled_time"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "ns2"}, Status: v1.PodStatus{Conditions: []v1.PodCondition{v1.PodCondition{Type: v1.PodScheduled, Status: v1.ConditionFalse}}}}, Want: metadata + `
				kube_pod_status_scheduled{condition="false",namespace="ns2",pod="pod2"} 1
				kube_pod_status_scheduled{condition="true",namespace="ns2",pod="pod2"} 0
				kube_pod_status_scheduled{condition="unknown",namespace="ns2",pod="pod2"} 0
			`, MetricNames: []string{"kube_pod_status_scheduled", "kube_pod_status_scheduled_time"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"}, Spec: v1.PodSpec{NodeName: "node1", Containers: []v1.Container{v1.Container{Name: "pod1_con1", Resources: v1.ResourceRequirements{Requests: map[v1.ResourceName]resource.Quantity{v1.ResourceCPU: resource.MustParse("200m"), v1.ResourceMemory: resource.MustParse("100M"), v1.ResourceEphemeralStorage: resource.MustParse("300M"), v1.ResourceStorage: resource.MustParse("400M"), v1.ResourceName("nvidia.com/gpu"): resource.MustParse("1")}, Limits: map[v1.ResourceName]resource.Quantity{v1.ResourceCPU: resource.MustParse("200m"), v1.ResourceMemory: resource.MustParse("100M"), v1.ResourceEphemeralStorage: resource.MustParse("300M"), v1.ResourceStorage: resource.MustParse("400M"), v1.ResourceName("nvidia.com/gpu"): resource.MustParse("1")}}}, v1.Container{Name: "pod1_con2", Resources: v1.ResourceRequirements{Requests: map[v1.ResourceName]resource.Quantity{v1.ResourceCPU: resource.MustParse("300m"), v1.ResourceMemory: resource.MustParse("200M")}, Limits: map[v1.ResourceName]resource.Quantity{v1.ResourceCPU: resource.MustParse("300m"), v1.ResourceMemory: resource.MustParse("200M")}}}}}}, Want: metadata + `
				kube_pod_container_resource_requests_cpu_cores{container="pod1_con1",namespace="ns1",node="node1",pod="pod1"} 0.2
				kube_pod_container_resource_requests_cpu_cores{container="pod1_con2",namespace="ns1",node="node1",pod="pod1"} 0.3
				kube_pod_container_resource_requests_memory_bytes{container="pod1_con1",namespace="ns1",node="node1",pod="pod1"} 1e+08
				kube_pod_container_resource_requests_memory_bytes{container="pod1_con2",namespace="ns1",node="node1",pod="pod1"} 2e+08
				kube_pod_container_resource_limits_cpu_cores{container="pod1_con1",namespace="ns1",node="node1",pod="pod1"} 0.2
				kube_pod_container_resource_limits_cpu_cores{container="pod1_con2",namespace="ns1",node="node1",pod="pod1"} 0.3
				kube_pod_container_resource_limits_memory_bytes{container="pod1_con1",namespace="ns1",node="node1",pod="pod1"} 1e+08
				kube_pod_container_resource_limits_memory_bytes{container="pod1_con2",namespace="ns1",node="node1",pod="pod1"} 2e+08
				kube_pod_container_resource_requests{container="pod1_con1",namespace="ns1",node="node1",pod="pod1",resource="cpu",unit="core"} 0.2
				kube_pod_container_resource_requests{container="pod1_con2",namespace="ns1",node="node1",pod="pod1",resource="cpu",unit="core"} 0.3
				kube_pod_container_resource_requests{container="pod1_con1",namespace="ns1",node="node1",pod="pod1",resource="nvidia_com_gpu",unit="integer"} 1
				kube_pod_container_resource_requests{container="pod1_con1",namespace="ns1",node="node1",pod="pod1",resource="memory",unit="byte"} 1e+08
				kube_pod_container_resource_requests{container="pod1_con2",namespace="ns1",node="node1",pod="pod1",resource="memory",unit="byte"} 2e+08
				kube_pod_container_resource_requests{container="pod1_con1",namespace="ns1",node="node1",pod="pod1",resource="storage",unit="byte"} 4e+08
				kube_pod_container_resource_requests{container="pod1_con1",namespace="ns1",node="node1",pod="pod1",resource="ephemeral_storage",unit="byte"} 3e+08
				kube_pod_container_resource_limits{container="pod1_con1",namespace="ns1",node="node1",pod="pod1",resource="cpu",unit="core"} 0.2
				kube_pod_container_resource_limits{container="pod1_con1",namespace="ns1",node="node1",pod="pod1",resource="nvidia_com_gpu",unit="integer"} 1
				kube_pod_container_resource_limits{container="pod1_con2",namespace="ns1",node="node1",pod="pod1",resource="cpu",unit="core"} 0.3
				kube_pod_container_resource_limits{container="pod1_con1",namespace="ns1",node="node1",pod="pod1",resource="memory",unit="byte"} 1e+08
				kube_pod_container_resource_limits{container="pod1_con2",namespace="ns1",node="node1",pod="pod1",resource="memory",unit="byte"} 2e+08
				kube_pod_container_resource_limits{container="pod1_con1",namespace="ns1",node="node1",pod="pod1",resource="storage",unit="byte"} 4e+08
				kube_pod_container_resource_limits{container="pod1_con1",namespace="ns1",node="node1",pod="pod1",resource="ephemeral_storage",unit="byte"} 3e+08
		`, MetricNames: []string{"kube_pod_container_resource_requests_cpu_cores", "kube_pod_container_resource_requests_memory_bytes", "kube_pod_container_resource_limits_cpu_cores", "kube_pod_container_resource_limits_memory_bytes", "kube_pod_container_resource_requests", "kube_pod_container_resource_limits"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod2", Namespace: "ns2"}, Spec: v1.PodSpec{NodeName: "node2", Containers: []v1.Container{v1.Container{Name: "pod2_con1", Resources: v1.ResourceRequirements{Requests: map[v1.ResourceName]resource.Quantity{v1.ResourceCPU: resource.MustParse("400m"), v1.ResourceMemory: resource.MustParse("300M")}, Limits: map[v1.ResourceName]resource.Quantity{v1.ResourceCPU: resource.MustParse("400m"), v1.ResourceMemory: resource.MustParse("300M")}}}, v1.Container{Name: "pod2_con2", Resources: v1.ResourceRequirements{Requests: map[v1.ResourceName]resource.Quantity{v1.ResourceCPU: resource.MustParse("500m"), v1.ResourceMemory: resource.MustParse("400M")}, Limits: map[v1.ResourceName]resource.Quantity{v1.ResourceCPU: resource.MustParse("500m"), v1.ResourceMemory: resource.MustParse("400M")}}}, v1.Container{Name: "pod2_con3"}}}}, Want: metadata + `
				kube_pod_container_resource_requests_cpu_cores{container="pod2_con1",namespace="ns2",node="node2",pod="pod2"} 0.4
				kube_pod_container_resource_requests_cpu_cores{container="pod2_con2",namespace="ns2",node="node2",pod="pod2"} 0.5
				kube_pod_container_resource_requests_memory_bytes{container="pod2_con1",namespace="ns2",node="node2",pod="pod2"} 3e+08
				kube_pod_container_resource_requests_memory_bytes{container="pod2_con2",namespace="ns2",node="node2",pod="pod2"} 4e+08
				kube_pod_container_resource_limits_cpu_cores{container="pod2_con1",namespace="ns2",node="node2",pod="pod2"} 0.4
				kube_pod_container_resource_limits_cpu_cores{container="pod2_con2",namespace="ns2",node="node2",pod="pod2"} 0.5
				kube_pod_container_resource_limits_memory_bytes{container="pod2_con1",namespace="ns2",node="node2",pod="pod2"} 3e+08
				kube_pod_container_resource_limits_memory_bytes{container="pod2_con2",namespace="ns2",node="node2",pod="pod2"} 4e+08
				kube_pod_container_resource_requests{container="pod2_con1",namespace="ns2",node="node2",pod="pod2",resource="cpu",unit="core"} 0.4
				kube_pod_container_resource_requests{container="pod2_con2",namespace="ns2",node="node2",pod="pod2",resource="cpu",unit="core"} 0.5
				kube_pod_container_resource_requests{container="pod2_con1",namespace="ns2",node="node2",pod="pod2",resource="memory",unit="byte"} 3e+08
				kube_pod_container_resource_requests{container="pod2_con2",namespace="ns2",node="node2",pod="pod2",resource="memory",unit="byte"} 4e+08
				kube_pod_container_resource_limits{container="pod2_con1",namespace="ns2",node="node2",pod="pod2",resource="cpu",unit="core"} 0.4
				kube_pod_container_resource_limits{container="pod2_con2",namespace="ns2",node="node2",pod="pod2",resource="cpu",unit="core"} 0.5
				kube_pod_container_resource_limits{container="pod2_con1",namespace="ns2",node="node2",pod="pod2",resource="memory",unit="byte"} 3e+08
				kube_pod_container_resource_limits{container="pod2_con2",namespace="ns2",node="node2",pod="pod2",resource="memory",unit="byte"} 4e+08
		`, MetricNames: []string{"kube_pod_container_resource_requests_cpu_cores", "kube_pod_container_resource_requests_cpu_cores", "kube_pod_container_resource_requests_memory_bytes", "kube_pod_container_resource_requests_memory_bytes", "kube_pod_container_resource_limits_cpu_cores", "kube_pod_container_resource_limits_cpu_cores", "kube_pod_container_resource_limits_memory_bytes", "kube_pod_container_resource_limits_memory_bytes", "kube_pod_container_resource_requests", "kube_pod_container_resource_requests", "kube_pod_container_resource_requests", "kube_pod_container_resource_requests", "kube_pod_container_resource_limits", "kube_pod_container_resource_limits", "kube_pod_container_resource_limits", "kube_pod_container_resource_limits"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1", Labels: map[string]string{"app": "example"}}, Spec: v1.PodSpec{}}, Want: metadata + `
				kube_pod_labels{label_app="example",namespace="ns1",pod="pod1"} 1
		`, MetricNames: []string{"kube_pod_labels"}}, {Obj: &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1", Labels: map[string]string{"app": "example"}}, Spec: v1.PodSpec{Volumes: []v1.Volume{v1.Volume{Name: "myvol", VolumeSource: v1.VolumeSource{PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{ClaimName: "claim1", ReadOnly: false}}}, v1.Volume{Name: "my-readonly-vol", VolumeSource: v1.VolumeSource{PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{ClaimName: "claim2", ReadOnly: true}}}, v1.Volume{Name: "not-pvc-vol", VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{Medium: "memory"}}}}}}, Want: metadata + `
				kube_pod_spec_volumes_persistentvolumeclaims_info{namespace="ns1",persistentvolumeclaim="claim1",pod="pod1",volume="myvol"} 1
				kube_pod_spec_volumes_persistentvolumeclaims_info{namespace="ns1",persistentvolumeclaim="claim2",pod="pod1",volume="my-readonly-vol"} 1
				kube_pod_spec_volumes_persistentvolumeclaims_readonly{namespace="ns1",persistentvolumeclaim="claim1",pod="pod1",volume="myvol"} 0
				kube_pod_spec_volumes_persistentvolumeclaims_readonly{namespace="ns1",persistentvolumeclaim="claim2",pod="pod1",volume="my-readonly-vol"} 1

		`, MetricNames: []string{"kube_pod_spec_volumes_persistentvolumeclaims_info", "kube_pod_spec_volumes_persistentvolumeclaims_readonly"}}}
	for i, c := range cases {
		c.Func = composeMetricGenFuncs(podMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
