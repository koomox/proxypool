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
proxypool.go [source](https://github.com/koomox/go-example/blob/master/source/proxypool.go)          
Use          
```sh
export GO111MODULE=on
export GOPROXY=https://goproxy.io
wget -O proxypool.go https://raw.githubusercontent.com/koomox/go-example/master/source/proxypool.go
go mod init .
go build -o proxypool proxypool.go
```
```bat
SET GO111MODULE=on
SET GOPROXY=https://goproxy.io
go mod init .
go build -o proxypool.exe proxypool.go

SET GOPATH=""
go install github.com/koomox/ext
go install github.com/koomox/proxypool
go build -o proxypool.exe proxypool.go
```