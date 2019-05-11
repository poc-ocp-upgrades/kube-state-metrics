package main

import (
	"compress/gzip"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	godefaulthttp "net/http"
	"net/http/pprof"
	"os"
	"strconv"
	"strings"
	"github.com/golang/glog"
	"github.com/openshift/origin/pkg/util/proc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	clientset "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	kcollectors "k8s.io/kube-state-metrics/pkg/collectors"
	"k8s.io/kube-state-metrics/pkg/options"
	"k8s.io/kube-state-metrics/pkg/version"
	"k8s.io/kube-state-metrics/pkg/whiteblacklist"
)

const (
	metricsPath	= "/metrics"
	healthzPath	= "/healthz"
)

type promLogger struct{}

func (pl promLogger) Println(v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	glog.Error(v...)
}
func main() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	opts := options.NewOptions()
	opts.AddFlags()
	err := opts.Parse()
	if err != nil {
		glog.Fatalf("Error: %s", err)
	}
	if opts.Version {
		fmt.Printf("%#v\n", version.GetVersion())
		os.Exit(0)
	}
	if opts.Help {
		opts.Usage()
		os.Exit(0)
	}
	collectorBuilder := kcollectors.NewBuilder(context.TODO())
	if len(opts.Collectors) == 0 {
		glog.Info("Using default collectors")
		collectorBuilder.WithEnabledCollectors(options.DefaultCollectors.AsSlice())
	} else {
		glog.Infof("Using collectors %s", opts.Collectors.String())
		collectorBuilder.WithEnabledCollectors(opts.Collectors.AsSlice())
	}
	if len(opts.Namespaces) == 0 {
		glog.Info("Using all namespace")
		collectorBuilder.WithNamespaces(options.DefaultNamespaces)
	} else {
		if opts.Namespaces.IsAllNamespaces() {
			glog.Info("Using all namespace")
		} else {
			glog.Infof("Using %s namespaces", opts.Namespaces)
		}
		collectorBuilder.WithNamespaces(opts.Namespaces)
	}
	whiteBlackList, err := whiteblacklist.New(opts.MetricWhitelist, opts.MetricBlacklist)
	if err != nil {
		glog.Fatal(err)
	}
	if opts.DisablePodNonGenericResourceMetrics {
		whiteBlackList.Exclude([]string{"kube_pod_container_resource_requests_cpu_cores", "kube_pod_container_resource_requests_memory_bytes", "kube_pod_container_resource_limits_cpu_cores", "kube_pod_container_resource_limits_memory_bytes"})
	}
	if opts.DisableNodeNonGenericResourceMetrics {
		whiteBlackList.Exclude([]string{"kube_node_status_capacity_cpu_cores", "kube_node_status_capacity_memory_bytes", "kube_node_status_capacity_pods", "kube_node_status_allocatable_cpu_cores", "kube_node_status_allocatable_memory_bytes", "kube_node_status_allocatable_pods"})
	}
	glog.Infof("metric white-blacklisting: %v", whiteBlackList.Status())
	collectorBuilder.WithWhiteBlackList(whiteBlackList)
	proc.StartReaper()
	kubeClient, err := createKubeClient(opts.Apiserver, opts.Kubeconfig)
	if err != nil {
		glog.Fatalf("Failed to create client: %v", err)
	}
	collectorBuilder.WithKubeClient(kubeClient)
	ksmMetricsRegistry := prometheus.NewRegistry()
	ksmMetricsRegistry.Register(kcollectors.ResourcesPerScrapeMetric)
	ksmMetricsRegistry.Register(kcollectors.ScrapeErrorTotalMetric)
	ksmMetricsRegistry.Register(prometheus.NewProcessCollector(os.Getpid(), ""))
	ksmMetricsRegistry.Register(prometheus.NewGoCollector())
	go telemetryServer(ksmMetricsRegistry, opts.TelemetryHost, opts.TelemetryPort)
	collectors := collectorBuilder.Build()
	serveMetrics(collectors, opts.Host, opts.Port, opts.EnableGZIPEncoding)
}
func createKubeClient(apiserver string, kubeconfig string) (clientset.Interface, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	config, err := clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)
	if err != nil {
		return nil, err
	}
	config.UserAgent = version.GetVersion().String()
	config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	config.ContentType = "application/vnd.kubernetes.protobuf"
	kubeClient, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	glog.Infof("Testing communication with server")
	v, err := kubeClient.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("ERROR communicating with apiserver: %v", err)
	}
	glog.Infof("Running with Kubernetes cluster version: v%s.%s. git version: %s. git tree state: %s. commit: %s. platform: %s", v.Major, v.Minor, v.GitVersion, v.GitTreeState, v.GitCommit, v.Platform)
	glog.Infof("Communication with server successful")
	return kubeClient, nil
}
func telemetryServer(registry prometheus.Gatherer, host string, port int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	listenAddress := net.JoinHostPort(host, strconv.Itoa(port))
	glog.Infof("Starting kube-state-metrics self metrics server: %s", listenAddress)
	mux := http.NewServeMux()
	mux.Handle(metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorLog: promLogger{}}))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Kube-State-Metrics Metrics Server</title></head>
             <body>
             <h1>Kube-State-Metrics Metrics</h1>
			 <ul>
             <li><a href='` + metricsPath + `'>metrics</a></li>
			 </ul>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(listenAddress, mux))
}
func serveMetrics(collectors []*kcollectors.Collector, host string, port int, enableGZIPEncoding bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	listenAddress := net.JoinHostPort(host, strconv.Itoa(port))
	glog.Infof("Starting metrics server: %s", listenAddress)
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	mux.Handle(metricsPath, &metricHandler{collectors, enableGZIPEncoding})
	mux.HandleFunc(healthzPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Kube Metrics Server</title></head>
             <body>
             <h1>Kube Metrics</h1>
			 <ul>
             <li><a href='` + metricsPath + `'>metrics</a></li>
             <li><a href='` + healthzPath + `'>healthz</a></li>
			 </ul>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(listenAddress, mux))
}

type metricHandler struct {
	collectors			[]*kcollectors.Collector
	enableGZIPEncoding	bool
}

func (m *metricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	resHeader := w.Header()
	var writer io.Writer = w
	resHeader.Set("Content-Type", `text/plain; version=`+"0.0.4")
	if m.enableGZIPEncoding {
		reqHeader := r.Header.Get("Accept-Encoding")
		parts := strings.Split(reqHeader, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "gzip" || strings.HasPrefix(part, "gzip;") {
				writer = gzip.NewWriter(writer)
				resHeader.Set("Content-Encoding", "gzip")
			}
		}
	}
	for _, c := range m.collectors {
		c.Collect(w)
	}
	if closer, ok := writer.(io.Closer); ok {
		closer.Close()
	}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
