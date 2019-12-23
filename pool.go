package proxypool

import (
	"github.com/koomox/ext"
	"strings"
	"sync"
)

var (
	DefaultTestURL = "https://raw.githubusercontent.com/koomox/Test/master/test"
)

const (
	ProxyHTTP byte = 1 << iota
	ProxyHTTPS
	ProxySOCKS4
	ProxySOCKS5
)

type proxyPool struct {
	list map[string]*proxyItem
}

type proxyItem struct {
	protocol byte
	addr     string
}

func New(uri string) *proxyPool {
	pool := &proxyPool{list: make(map[string]*proxyItem)}
	addrs := loadProxyList()
	if addrs == nil {
		return pool
	}

	reqURI := DefaultTestURL
	if uri == "" {
		reqURI = uri
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(addrs))

	for _, item := range addrs {
		go proxyAdd(item.protocol, item.addr, reqURI, pool, wg)
	}

	wg.Wait()

	return pool
}

func (this *proxyPool) Get() (addrs []string) {
	for _, proxy := range this.list {
		if proxy.protocol&ProxySOCKS5 == ProxySOCKS5 {
			addrs = append(addrs, "socks5://"+proxy.addr)
		}
		if proxy.protocol&ProxyHTTP == ProxyHTTP {
			addrs = append(addrs, "http://"+proxy.addr)
		}
	}

	return
}

func (this *proxyPool) add(protocol byte, addr string) {
	proxy := &proxyItem{protocol: protocol, addr: addr}
	md5sum := ext.GetMD5(addr)
	if v, ok := this.list[md5sum]; !ok {
		this.list[md5sum] = proxy
	} else {
		v.protocol |= protocol
	}
}

func (this *proxyPool) delete(addr string) {
	md5sum := ext.GetMD5(addr)
	if _, ok := this.list[md5sum]; ok {
		delete(this.list, md5sum)
	}
}

func proxyAdd(protocol byte, addr, uri string, pool *proxyPool, wg *sync.WaitGroup) {
	v := Proxy(protocol, addr)
	if v != nil {
		if v.Test(uri) {
			pool.add(protocol, addr)
		}
	}
	wg.Done()
}

func protocolCode(protocol string) byte {
	switch strings.ToUpper(protocol) {
	case "HTTP":
		return ProxyHTTP
	case "HTTPS":
		return ProxyHTTPS
	case "SOCKS4":
		return ProxySOCKS4
	case "SOCKS5":
		return ProxySOCKS5
	default:
		return 0
	}
}
