package proxypool

import (
	"errors"
	"time"
)

var (
	proxyTimeout = 10 * time.Second
)

type proxyer interface {
	Test(string) bool
	Get(string) ([]byte, error)
	GetFile(string, string) error
}

func Proxy(protocol byte, proxyAddr string) (proxy proxyer) {
	switch protocol {
	case ProxyHTTP:
		proxy = &httpProxy{addr: proxyAddr}
	case ProxySOCKS5:
		proxy = &socks5Proxy{addr: proxyAddr}
	default:
		proxy = nil
	}

	return
}

func (this *proxyPool) ProxyHttpGet(reqAddr string) (ctx []byte, err error) {
	for _, proxy := range this.list {
		v := Proxy(proxy.protocol, proxy.addr)
		if v != nil {
			if ctx, err = v.Get(reqAddr); err == nil {
				return
			}
		}
	}

	return ctx, errors.New("Not Proxy")
}

func (this *proxyPool) ProxyHttpGetFile(reqAddr, dst string) (err error) {
	for _, proxy := range this.list {
		v := Proxy(proxy.protocol, proxy.addr)
		if v != nil {
			if err = v.GetFile(reqAddr, dst); err == nil {
				return
			}
		}
	}

	return errors.New("Not Proxy")
}
