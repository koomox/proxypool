package proxypool

import (
	"fmt"
	"net/http"
)

func HttpHandleFunc(w http.ResponseWriter, r *http.Request) {
	ctx := handleFunc()
	if ctx == "" {
		ctx = "null"
	}
	fmt.Fprintln(w, ctx)
}

func handleFunc() (ctx string) {
	addrs := ProxyPool.Get()
	for _, addr := range addrs {
		ctx += addr + "\n"
	}

	return
}
