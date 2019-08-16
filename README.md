# proxypool        

### Install     
Dependent      
```sh
go get -u -d github.com/koomox/proxypool

go get -u -d github.com/koomox/ext
go get -u -d golang.org/x/sys
go get -u -d golang.org/x/net
go get -u -d github.com/PuerkitoBio/goquery
go get -u -d github.com/qiniu/log
```

Use          
```go
package main

import (
	"github.com/koomox/proxypool"
	"github.com/qiniu/log"
	"net/http"
	"strings"
)

var (
	reqURL = "https://raw.githubusercontent.com/torvalds/linux/master/README"
	msURL  = "https://raw.githubusercontent.com/microsoft/vscode/master/LICENSE.txt"
	cmdQ   = make(chan string)
)

func main() {
	go httpServ("127.0.0.1:8000")
	go proxyHttpGet(reqURL)
	go proxyHttpGet(msURL)

	log.Info("Wait...")
	select {
	case cmd := <-cmdQ: // 收到控制指令
		if strings.EqualFold(cmd, "quit") {
			log.Error("quit")
			break
		}
	}
}

func httpServ(addr string) {
	http.HandleFunc("/", proxypool.HttpHandleFunc)
	http.ListenAndServe(addr, nil)
}

func proxyHttpGet(reqAddr string) {
	if ctx, err := proxypool.ProxyHttpGet(reqAddr); err == nil {
		log.Infof("Http Content:\n%v", string(ctx))
	} else {
		log.Errorf("Err: %v", err.Error())
	}
}
```