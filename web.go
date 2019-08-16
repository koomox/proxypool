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
	for _, proxy := range ProxyPool.pool {
		ctx += proxy.protocol + "://" + proxy.addr + "\n"
	}

	return
}
