package proxypool

import (
	"github.com/koomox/ext"
	"strings"
)

var (
	DefaultTestURL = "https://raw.githubusercontent.com/koomox/Test/master/test"
)

type proxyPool struct {
	pool map[string]*proxyInfo
}

type proxyInfo struct {
	protocol string
	addr     string
}

var ProxyPool = &proxyPool{pool: make(map[string]*proxyInfo)}

func (this *proxyPool) Get() (addrs []string) {
	for _, proxy := range this.pool {
		addrs = append(addrs, proxy.Encode())
	}

	return
}

func (proxy *proxyInfo) Encode() string {
	return proxy.protocol + "://" + proxy.addr
}

func (proxy *proxyInfo) setProtocol(protocol ...string) {
	sp := strings.Split(proxy.protocol, ",")
	for _, p := range protocol {
		if !proxy.findProtocol(p) {
			sp = append(sp, p)
		}
	}
	proxy.protocol = strings.Join(sp, ",")
}

func (proxy *proxyInfo) findProtocol(protocol string) bool {
	sp := strings.Split(proxy.protocol, ",")
	for _, p := range sp {
		if p == protocol {
			return true
		}
	}
	return false
}

func (this *proxyPool) add(protocol, addr string) {
	proxy := &proxyInfo{protocol: protocol, addr: addr}
	md5sum := ext.GetMD5(addr)
	if v, ok := this.pool[md5sum]; !ok {
		this.pool[md5sum] = proxy
	} else {
		if v.protocol != proxy.protocol {
			v.setProtocol(proxy.protocol)
		}
	}
}

func (this *proxyPool) delete(addr string) {
	md5sum := ext.GetMD5(addr)
	if _, ok := this.pool[md5sum]; ok {
		delete(this.pool, md5sum)
	}
}

func proxyAdd(protocol, proxyAddr, reqAddr string) {
	if Proxy(protocol, proxyAddr).Test(reqAddr) {
		ProxyPool.add(protocol, proxyAddr)
	}
}
