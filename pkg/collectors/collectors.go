package collectors

import (
	"io"
)

type Store interface{ WriteAll(io.Writer) }
type Collector struct{ Store Store }

func NewCollector(s Store) *Collector {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Collector{s}
}
func (c *Collector) Collect(w io.Writer) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.Store.WriteAll(w)
}
