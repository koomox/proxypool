package proxypool

import (
	"fmt"
	"sync"
)

func init() {
	LoadProxyList()
}

type urlItem struct {
	method   string
	protocol string
	addrs    []string
}

func LoadProxyList() {
	urls := []urlItem{
		urlItem{method: "table", protocol: "http", addrs: urlKuaiDaiLi()},
		urlItem{method: "text", protocol: "socks5", addrs: urlProxyListUS()},
		urlItem{method: "text", protocol: "socks5", addrs: urlProxyScrape()},
		urlItem{method: "text", protocol: "http", addrs: urlProxyListPlus()},
		urlItem{method: "table", protocol: "http", addrs: urlProxyXiCiDaiLi()},
		urlItem{method: "text", protocol: "http", addrs: urlProxySpysme()},
	}

	var wg sync.WaitGroup
	for _, uri := range urls {
		for _, addr := range uri.addrs {
			items := ParseAddrURL(uri.method, uri.protocol, addr)
			for _, item := range items {
				wg.Add(1)
				go func() {
					proxyAdd(item.protocol, item.addr, DefaultTestURL)
					wg.Done()
				}()
			}
		}
	}
	wg.Wait()
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
