package ext

import (
	"regexp"
)

const (
	portPattern = `(\d{1,5})`

	// 匹配 IP4
	ip4Pattern = `((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`

	ip4AddrPattern = `((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?):\d{1,5}`

	// 匹配 IP6，参考以下网页内容：
	// http://blog.csdn.net/jiangfeng08/article/details/7642018
	ip6Pattern = `(([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|` +
		`(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|` +
		`(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|` +
		`(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))`

	ip6AddrPattern = ip6Pattern + `:\d{1,5}`

	// 同时匹配 IP4 和 IP6
	ipPattern = "(" + ip4Pattern + ")|(" + ip6Pattern + ")"

	ipAddrPattern = "(" + ip4AddrPattern + ")|(" + ip6AddrPattern + ")"

	// 匹配域名
	domainPattern = `[a-zA-Z0-9][a-zA-Z0-9_-]{0,62}(\.[a-zA-Z0-9][a-zA-Z0-9_-]{0,62})*(\.[a-zA-Z][a-zA-Z0-9]{0,10}){1}`

	// 匹配 URL
	urlPattern = `((https|http|ftp|rtsp|mms)?://)?` + // 协议
		`(([0-9a-zA-Z]+:)?[0-9a-zA-Z_-]+@)?` + // pwd:user@
		"(" + ipPattern + "|(" + domainPattern + "))" + // IP 或域名
		`(:\d{1,5})?` + // 端口
		`(/+[a-zA-Z0-9][a-zA-Z0-9_.-]*)*/*` + // path
		`(\?([a-zA-Z0-9_-]+(=.*&?)*)*)*` // query
)

var (
	portExpMustCompile    = regexpMustCompile(portPattern)
	ip4ExpCompile         = regexpCompile(ip4Pattern)
	ip4ExpMustCompile     = regexpMustCompile(ip4Pattern)
	ip4AddrExpMustCompile = regexpMustCompile(ip4AddrPattern)
	ip6ExpCompile         = regexpCompile(ip6Pattern)
	ipExpCompile          = regexpCompile(ipPattern)
	ipExpMustCompile      = regexpMustCompile(ipPattern)
	ipAddrExpMustCompile  = regexpMustCompile(ipAddrPattern)
	regexpURLExpCompile   = regexpCompile(urlPattern)
	domainExpCompile      = regexpCompile(domainPattern)
	domainExpMustCompile  = regexpMustCompile(domainPattern)
)

func regexpCompile(str string) *regexp.Regexp {
	return regexp.MustCompile("^" + str + "$")
}

func regexpMustCompile(str string) *regexp.Regexp {
	return regexp.MustCompile(str)
}

// 判断val是否能正确匹配exp中的正则表达式。
// val可以是[]byte, []rune, string类型。
func isMatch(exp *regexp.Regexp, val interface{}) bool {
	switch v := val.(type) {
	case []rune:
		return exp.MatchString(string(v))
	case []byte:
		return exp.Match(v)
	case string:
		return exp.MatchString(v)
	default:
		return false
	}
}

func matchFindString(exp *regexp.Regexp, val interface{}) string {
	switch v := val.(type) {
	case []byte:
		return exp.FindString(string(v))
	case string:
		return exp.FindString(v)
	default:
		return ""
	}
}

// URL 验证一个值是否标准的URL格式。支持IP和域名等格式
func MatchURL(val interface{}) bool {
	return isMatch(regexpURLExpCompile, val)
}

// IP 验证一个值是否为IP，可验证IP4和IP6
func MatchIP(val interface{}) bool {
	return isMatch(ipExpCompile, val)
}

// IP6 验证一个值是否为IP6
func MatchIP6(val interface{}) bool {
	return isMatch(ip6ExpCompile, val)
}

// IP4 验证一个值是滞为IP4
func MatchIP4(val interface{}) bool {
	return isMatch(ip4ExpCompile, val)
}

func MatchDomain(val interface{}) bool {
	return isMatch(domainExpCompile, val)
}

func MatchFindIPv4(val interface{}) string {
	return matchFindString(ip4ExpMustCompile, val)
}

func MatchFindIP(val interface{}) string {
	return matchFindString(ipExpMustCompile, val)
}

func MatchFindAddr(val interface{}) string {
	return matchFindString(ipAddrExpMustCompile, val)
}

func MatchFindIPv4Addr(val interface{}) string {
	return matchFindString(ip4AddrExpMustCompile, val)
}

func MatchFindPort(val interface{}) string {
	return matchFindString(portExpMustCompile, val)
}

func MatchFindDomain(val interface{}) string {
	return matchFindString(domainExpMustCompile, val)
}

func DomainHasSuffix(s, suffix string) bool {
	if !MatchDomain(s) {
		return false
	}
	if len(s) < len(suffix) { // 字符串长度小于后缀长度
		return false
	}
	if s[len(s)-len(suffix):] != suffix { // 字符串后缀不等于后缀
		return false
	}
	if len(s) != len(suffix) { // 字符串长度不等于后缀长度
		if s[len(s)-len(suffix)-1:len(s)-len(suffix)] != "." {
			return false
		}
	}

	return true
}
