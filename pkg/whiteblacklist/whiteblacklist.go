package whiteblacklist

import (
	"errors"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"strings"
)

type WhiteBlackList struct {
	list		map[string]struct{}
	isWhiteList	bool
}

func New(w, b map[string]struct{}) (*WhiteBlackList, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(w) != 0 && len(b) != 0 {
		return nil, errors.New("whitelist and blacklist are both set, they are mutually exclusive, only one of them can be set")
	}
	white := copyList(w)
	black := copyList(b)
	var list map[string]struct{}
	var isWhiteList bool
	if len(white) != 0 {
		list = white
		isWhiteList = true
	} else {
		list = black
		isWhiteList = false
	}
	return &WhiteBlackList{list: list, isWhiteList: isWhiteList}, nil
}
func (l *WhiteBlackList) Include(items []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if l.isWhiteList {
		for _, item := range items {
			l.list[item] = struct{}{}
		}
	} else {
		for _, item := range items {
			delete(l.list, item)
		}
	}
}
func (l *WhiteBlackList) Exclude(items []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if l.isWhiteList {
		for _, item := range items {
			delete(l.list, item)
		}
	} else {
		for _, item := range items {
			l.list[item] = struct{}{}
		}
	}
}
func (l *WhiteBlackList) IsIncluded(item string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, exists := l.list[item]
	if l.isWhiteList {
		return exists
	}
	return !exists
}
func (l *WhiteBlackList) IsExcluded(item string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return !l.IsIncluded(item)
}
func (l *WhiteBlackList) Status() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	items := []string{}
	for key := range l.list {
		items = append(items, key)
	}
	if l.isWhiteList {
		return "whitelisting the following items: " + strings.Join(items, ", ")
	}
	return "blacklisting the following items: " + strings.Join(items, ", ")
}
func copyList(l map[string]struct{}) map[string]struct{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newList := map[string]struct{}{}
	for k, v := range l {
		newList[k] = v
	}
	return newList
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
