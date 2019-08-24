package proxypool

import (
	"errors"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
)

type socks5Proxy struct {
	addr string
}

func (this *socks5Proxy) createSocksClient() *http.Client {
	dt := net.Dialer{
		Timeout:   proxyTimeout,
		KeepAlive: proxyTimeout,
	}
	dialer, err := proxy.SOCKS5("tcp", this.addr, nil, &dt)
	if err != nil {
		return nil
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	client := &http.Client{
		Timeout:   proxyTimeout,
		Transport: transport,
	}

	return client
}

func (this *socks5Proxy) Test(reqAddr string) bool {
	var (
		reqURL   *url.URL
		response *http.Response
		err      error
	)
	if reqAddr == "" {
		reqAddr = DefaultTestURL
	}

	if reqURL, err = url.Parse(reqAddr); err != nil {
		return false
	}

	client := this.createSocksClient()
	if client == nil {
		return false
	}

	if response, err = client.Get(reqURL.String()); err != nil {
		return false
	}
	response.Body.Close()

	return true
}

func (this *socks5Proxy) Get(reqAddr string) (ctx []byte, err error) {
	var (
		reqURL   *url.URL
		response *http.Response
	)
	if reqURL, err = url.Parse(reqAddr); err != nil {
		return
	}

	client := this.createSocksClient()
	if client == nil {
		return nil, errors.New("createSocksClient failed!")
	}

	if response, err = client.Get(reqURL.String()); err != nil {
		return
	}
	defer response.Body.Close()

	ctx, err = ioutil.ReadAll(response.Body)
	return
}

func (this *socks5Proxy) GetFile(reqAddr, dst string) (err error) {
	var (
		reqURL   *url.URL
		response *http.Response
	)
	if reqURL, err = url.Parse(reqAddr); err != nil {
		return
	}

	if dst == "" {
		dst = path.Base(reqAddr)
	}

	client := this.createSocksClient()
	if client == nil {
		return errors.New("createSocksClient failed!")
	}

	if response, err = client.Get(reqURL.String()); err != nil {
		return
	}
	defer response.Body.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)

	return err
}
