package proxypool

import "fmt"

func init() {
	LoadProxyList()
}

func LoadProxyList() {
	loadProxyKuaiDaiLi()
	loadProxyScrape()
	loadProxyListUS()
	loadProxyListPlus()
	loadProxyXiCiDaiLi()
	loadProxySpysme()
}

func loadProxyKuaiDaiLi() {
	for i := 1; i < 3; i++ {
		proxyAddr := fmt.Sprintf("https://www.kuaidaili.com/free/intr/%d/", i)
		go LoadProxyURL("table", "http", proxyAddr)
	}
}

func loadProxyScrape() {
	country := []string{"US", "HK", "JP", "SG", "SK", "RU"}
	for _, c := range country {
		proxyAddr := fmt.Sprintf("https://api.proxyscrape.com/?request=getproxies&proxytype=socks5&timeout=5000&country=%s", c)
		go LoadProxyURL("text", "socks5", proxyAddr)
	}
}

func loadProxyListUS() {
	country := []string{"US", "HK", "JP", "SG", "SK", "RU"}
	for _, c := range country {
		proxyAddr := fmt.Sprintf("https://www.proxy-list.download/api/v1/get?type=socks5&anon=elite&country=%s", c)
		go LoadProxyURL("text", "socks5", proxyAddr)
	}
}

func loadProxyListPlus() {
	proxyAddr := "https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-1"
	go LoadProxyURL("text", "http", proxyAddr)
}

func loadProxyXiCiDaiLi() {
	proxyAddr := "https://www.xicidaili.com/wt/"
	go LoadProxyURL("table", "http", proxyAddr)
}

func loadProxySpysme() {
	proxyAddr := "http://spys.me/proxy.txt"
	go LoadProxyURL("text", "http", proxyAddr)
}
