package proxypool

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

type httpProxy struct {
	addr string
}

func (this *httpProxy) createHttpClient() *http.Client {
	proxyURL, err := url.Parse(this.addr)
	if err != nil {
		return nil
	}
	transport := &http.Transport{
		//MaxIdleConns:        10,
		//IdleConnTimeout:     30 * time.Second,
		//TLSHandshakeTimeout:   proxyTimeout,
		DisableKeepAlives: true,
		Proxy:             http.ProxyURL(proxyURL),
		//TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // skip cert verify
		ResponseHeaderTimeout: proxyTimeout, // Default while {}, Read HTML Header timeout
	}

	client := &http.Client{
		Timeout:   proxyTimeout,
		Transport: transport,
	}

	return client
}

func (this *httpProxy) Test(reqAddr string) bool {
	if reqAddr == "" {
		reqAddr = DefaultTestURL
	}
	if _, err := this.Get(reqAddr); err != nil {
		return false
	}
	return true
}

func (this *httpProxy) Get(reqAddr string) (ctx []byte, err error) {
	var (
		reqURL   *url.URL
		response *http.Response
	)

	if reqURL, err = url.Parse(reqAddr); err != nil {
		return
	}

	client := this.createHttpClient()
	if client == nil {
		return nil, errors.New("createHttpClient Failed!")
	}

	if response, err = client.Get(reqURL.String()); err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HttpProxyGet URL(%v) bad status %v", reqAddr, response.Status)
	}

	ctx, err = ioutil.ReadAll(response.Body)
	return
}

func (this *httpProxy) GetFile(reqAddr, dst string) (err error) {
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

	client := this.createHttpClient()
	if client == nil {
		return errors.New("createHttpClient Failed!")
	}

	if response, err = client.Get(reqURL.String()); err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HttpProxyGetFile URL(%v) bad status %v", reqAddr, response.Status)
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)

	return err
}
