package options

import (
	"reflect"
	"testing"
)

func TestCollectorSetSet(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		Desc		string
		Value		string
		Wanted		CollectorSet
		WantedError	bool
	}{{Desc: "empty collectors", Value: "", Wanted: CollectorSet{}, WantedError: false}, {Desc: "normal collectors", Value: "configmaps,cronjobs,daemonsets,deployments", Wanted: CollectorSet(map[string]struct{}{"configmaps": {}, "cronjobs": {}, "daemonsets": {}, "deployments": {}}), WantedError: false}, {Desc: "none exist collectors", Value: "none-exists", Wanted: CollectorSet{}, WantedError: true}}
	for _, test := range tests {
		cs := &CollectorSet{}
		gotError := cs.Set(test.Value)
		if !(((gotError == nil && !test.WantedError) || (gotError != nil && test.WantedError)) && reflect.DeepEqual(*cs, test.Wanted)) {
			t.Errorf("Test error for Desc: %s. Want: %+v. Got: %+v. Wanted Error: %v, Got Error: %v", test.Desc, test.Wanted, *cs, test.WantedError, gotError)
		}
	}
}
