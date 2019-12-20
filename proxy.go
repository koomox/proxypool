package proxypool

import (
	"errors"
	"strings"
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

func Proxy(protocol, proxyAddr string) proxyer {
	switch strings.ToUpper(protocol) {
	case "HTTP":
		return &httpProxy{addr: proxyAddr}
	case "SOCKS5":
		return &socks5Proxy{addr: proxyAddr}
	default:
		return nil
	}
}

func ProxyHttpGet(reqAddr string) (ctx []byte, err error) {
	for {
		if len(ProxyPool.pool) > 0 {
			break
		}
		ticker := time.NewTicker(300 * time.Millisecond)
		select {
		case <-ticker.C:
		}
	}
	for _, proxy := range ProxyPool.pool {
		v := Proxy(proxy.protocol, proxy.addr)
		if v != nil {
			if ctx, err = v.Get(reqAddr); err == nil {
				return
			}
		}
	}

	return ctx, errors.New("Not Proxy")
}

func ProxyHttpGetFile(reqAddr, dst string) (err error) {
	for {
		if len(ProxyPool.pool) > 0 {
			break
		}
		ticker := time.NewTicker(300 * time.Millisecond)
		select {
		case <-ticker.C:
		}
	}
	for _, proxy := range ProxyPool.pool {
		v := Proxy(proxy.protocol, proxy.addr)
		if v != nil {
			if err = v.GetFile(reqAddr, dst); err == nil {
				return
			}
		}
	}

	return errors.New("Not Proxy")
}
