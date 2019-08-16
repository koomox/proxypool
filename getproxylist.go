package proxypool

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/koomox/ext"
	"net/http"
	"strings"
)

func getTextProxyList(reqAddr string) (ctx []string, err error) {
	var (
		content []byte
		addrs   []string
	)
	ctx = make([]string, 0)
	if content, err = ext.HttpGetRaw(reqAddr); err != nil {
		return
	}

	addrs = strings.Split(string(content), "\n")
	for _, addr := range addrs {
		addr = getIPAddrPort(addr)
		if addr == "" {
			continue
		}
		ctx = append(ctx, addr)
	}

	return
}

func getTableProxyList(reqAddr string) (ctx []string, err error) {
	ctx = make([]string, 0)
	res, err := http.Get(reqAddr)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return ctx, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}
	var (
		addr string
		port string
	)
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		addr = ""
		port = ""
		// For each item found, get the band and title
		tr.Find("td").Each(func(i int, td *goquery.Selection) {
			if port == "" {
				switch addr {
				case "":
					proxyAddr := getIPAddr(td.Text())
					if proxyAddr != "" {
						addr = proxyAddr
					}
				default:
					proxyPort := getPort(td.Text())
					if proxyPort != "" {
						port = proxyPort
					}
				}
			}
		})
		if addr != "" && port != "" {
			ctx = append(ctx, addr+":"+port)
		}
	})
	return
}

func LoadProxyURL(method, protocol, reqAddr string) {
	switch strings.ToUpper(method) {
	case "TEXT":
		if ctx, err := getTextProxyList(reqAddr); err == nil {
			LoadProxy(protocol, ctx)
		}
	case "TABLE":
		if ctx, err := getTableProxyList(reqAddr); err == nil {
			LoadProxy(protocol, ctx)
		}
	}
}
