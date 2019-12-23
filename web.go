package proxypool

import (
	"fmt"
	"net/http"
)

func (pool *proxyPool) HttpHandleFunc(w http.ResponseWriter, r *http.Request) {
	ctx := handleFunc(pool)
	if ctx == "" {
		ctx = "null"
	}
	fmt.Fprintln(w, ctx)
}

func handleFunc(pool *proxyPool) (ctx string) {
	addrs := pool.Get()
	for _, addr := range addrs {
		ctx += addr + "\n"
	}

	return
}
