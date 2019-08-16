package ext

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

var (
	httpTimeout = 10 * time.Second
)

var httpClient = createHttpClient()

func createHttpClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns: 10,
		Dial: func(network, addr string) (net.Conn, error) {
			deadline := time.Now().Add(httpTimeout)
			c, err := net.DialTimeout(network, addr, httpTimeout)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)
			return c, nil
		},
		ResponseHeaderTimeout: httpTimeout,
	}

	client := &http.Client{
		Transport: transport,
	}

	return client
}

func HttpGetRaw(reqAddr string) (ctx []byte, err error) {
	var (
		response *http.Response
	)
	if _, err = url.Parse(reqAddr); err != nil {
		return
	}

	if response, err = httpClient.Get(reqAddr); err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HttpGetRaw URL(%v) bad status %v", reqAddr, response.Status)
	}

	ctx, err = ioutil.ReadAll(response.Body)
	return
}

func HttpGetFile(reqAddr, dst string) (err error) {
	var (
		response *http.Response
	)
	if dst == "" {
		dst = path.Base(reqAddr)
	}

	if _, err = url.Parse(reqAddr); err != nil {
		return
	}

	if response, err = httpClient.Get(reqAddr); err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HttpGetFile URL(%v) bad status %v", reqAddr, response.Status)
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)

	return err
}
