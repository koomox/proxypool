package proxypool

import (
	"crypto/md5"
	"io"
	"regexp"
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

func (this *proxyPool) Get() (ctx []string) {
	ctx = make([]string, 0)
	for _, proxy := range this.pool {
		ctx = append(ctx, proxy.Encode())
	}

	return ctx
}

func (proxy *proxyInfo) getMD5() (ctx string, err error) {
	h := md5.New()
	if _, err = io.WriteString(h, proxy.addr); err != nil {
		return
	}
	return string(h.Sum(nil)), nil
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
	var (
		md5sum string
		err    error
	)
	proxy := &proxyInfo{protocol: protocol, addr: addr}
	if md5sum, err = proxy.getMD5(); err != nil {
		return
	}
	if v, ok := this.pool[md5sum]; !ok {
		this.pool[md5sum] = proxy
	} else {
		if v.protocol != proxy.protocol {
			v.setProtocol(proxy.protocol)
		}
	}
}

func (this *proxyPool) delete(addr string) {
	var (
		md5sum string
		err    error
	)
	proxy := &proxyInfo{protocol: "", addr: addr}
	if md5sum, err = proxy.getMD5(); err != nil {
		return
	}
	if _, ok := this.pool[md5sum]; ok {
		delete(this.pool, md5sum)
	}
}

func proxyAdd(protocol, proxyAddr, reqAddr string) {
	if _, err := Proxy(protocol, proxyAddr).Get(reqAddr); err == nil {
		ProxyPool.add(protocol, proxyAddr)
	}
}

// return ip:port string or ""
func getIPAddrPort(addr string) string {
	exp := regexp.MustCompile(`((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?):\d{1,5}`)
	return exp.FindString(string(addr))
}

func getIPAddr(str string) string {
	exp := regexp.MustCompile(`((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`)
	return exp.FindString(string(str))
}

func getPort(str string) string {
	exp := regexp.MustCompile(`(\d{1,5})`)
	return exp.FindString(string(str))
}
