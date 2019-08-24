package proxypool

import (
	"fmt"
)

func init() {
	go LoadProxyList()
}

type urlAddrs struct {
	method   string
	protocol string
	addrs    []string
}

type urlItem struct {
	method   string
	protocol string
	addr     string
}

func LoadProxyList() {
	urls := []urlAddrs{
		urlAddrs{method: "table", protocol: "http", addrs: urlKuaiDaiLi()},
		urlAddrs{method: "text", protocol: "socks5", addrs: urlProxyListUS()},
		urlAddrs{method: "text", protocol: "socks5", addrs: urlProxyScrape()},
		urlAddrs{method: "text", protocol: "http", addrs: urlProxyListPlus()},
		urlAddrs{method: "table", protocol: "http", addrs: urlProxyXiCiDaiLi()},
		urlAddrs{method: "text", protocol: "http", addrs: urlProxySpysme()},
	}

	var items []urlItem
	for _, uri := range urls {
		for _, addr := range uri.addrs {
			items = append(items, urlItem{method: uri.method, protocol: uri.protocol, addr: addr})
		}
	}

	for _, item := range items {
		go func(method, protocol, addr string) {
			proxys := ParseAddrURL(method, protocol, addr)
			for _, proxy := range proxys {
				go func(protocol, addr, uri string) {
					proxyAdd(protocol, addr, uri)
				}(proxy.protocol, proxy.addr, DefaultTestURL)
			}
		}(item.method, item.protocol, item.addr)
	}
}

func urlKuaiDaiLi() (addrs []string) {
	for i := 1; i < 3; i++ {
		proxyAddr := fmt.Sprintf("https://www.kuaidaili.com/free/intr/%d/", i)
		addrs = append(addrs, proxyAddr)
	}
	return
}

func urlProxyListUS() (addrs []string) {
	country := []string{"US", "HK", "JP", "SG", "SK", "RU"}
	for _, c := range country {
		proxyAddr := fmt.Sprintf("https://www.proxy-list.download/api/v1/get?type=socks5&anon=elite&country=%s", c)
		addrs = append(addrs, proxyAddr)
	}
	return
}

func urlProxyScrape() (addrs []string) {
	country := []string{"US", "HK", "JP", "SG", "SK", "RU"}
	for _, c := range country {
		proxyAddr := fmt.Sprintf("https://api.proxyscrape.com/?request=getproxies&proxytype=socks5&timeout=5000&country=%s", c)
		addrs = append(addrs, proxyAddr)
	}
	return
}

func urlProxyListPlus() (addrs []string) {
	proxyAddr := "https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-1"
	addrs = append(addrs, proxyAddr)
	return
}

func urlProxyXiCiDaiLi() (addrs []string) {
	proxyAddr := "https://www.xicidaili.com/wt/"
	addrs = append(addrs, proxyAddr)
	return
}

func urlProxySpysme() (addrs []string) {
	proxyAddr := "http://spys.me/proxy.txt"
	addrs = append(addrs, proxyAddr)
	return
}
