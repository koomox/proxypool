package ext

// Get Public IP address
import (
	"fmt"
	"regexp"
)

var (
	akamaiCheckIPURI = "http://whatismyip.akamai.com/"
	amazonCheckIPURI = "https://checkip.amazonaws.com/"
)

func HttpGetPublicIPAddr(uri string) (string, error) {
	// 读取在线通用配置文件, 如果失败，使用代理尝试
	buf, err := HttpGetRaw(uri)
	if err != nil {
		return "", err
	}

	exp := regexp.MustCompile(`((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`)
	addr := exp.FindString(string(buf))
	if addr == "" {
		return addr, fmt.Errorf("Get Pulib IP Addr failed, URL: %v", uri)
	}

	return addr, nil
}

func GetPublicIPAddr() (addr string, err error) {
	if addr, err = HttpGetPublicIPAddr(akamaiCheckIPURI); err == nil {
		return addr, err
	}
	return HttpGetPublicIPAddr(amazonCheckIPURI)
}
