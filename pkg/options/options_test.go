package options

import (
	"os"
	"sync"
	"testing"
	"github.com/spf13/pflag"
)

func TestOptionsParse(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		Desc		string
		Args		[]string
		RecoverInvoked	bool
	}{{Desc: "collectors command line argument", Args: []string{"./kube-state-metrics", "--collectors=configmaps,pods"}, RecoverInvoked: false}, {Desc: "namespace command line argument", Args: []string{"./kube-state-metrics", "--namespace=default,kube-system"}, RecoverInvoked: false}}
	for _, test := range tests {
		var wg sync.WaitGroup
		opts := NewOptions()
		opts.AddFlags()
		flags := pflag.NewFlagSet("options_test", pflag.PanicOnError)
		flags.AddFlagSet(opts.flags)
		opts.flags = flags
		os.Args = test.Args
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					test.RecoverInvoked = true
				}
			}()
			opts.Parse()
		}()
		wg.Wait()
		if test.RecoverInvoked {
			t.Errorf("Test error for Desc: %s. Test panic", test.Desc)
		}
	}
}
