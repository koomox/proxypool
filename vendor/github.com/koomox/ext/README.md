# ext
Extensions Libs

```sh
go get -u -d golang.org/x/crypto
go get -u -d github.com/koomox/ext
```

### Use         
```go
package main

import (
	"github.com/koomox/ext"
	"fmt"
)

func main() {
	var (
		r string
		fs []string
		err error
	)
	if r, err = ext.RandomString(32); err != nil {
		fmt.Printf("Err:%v", err.Error())
		return
	}
	fmt.Printf("RandomString(\"%v\")\n", r)

	if r, err = ext.RandomSecurePassword(32); err != nil {
		fmt.Printf("Err:%v", err.Error())
		return
	}
	fmt.Printf("RandomSecurePassword(\"%v\")\n", r)
	
	fmt.Printf("str=\"%v\" MD5=\"%v\"\n", r, ext.GetMD5(r))

	if fs, err = ext.GetCustomDirectoryAllFile(""); err != nil {
		fmt.Printf("Err:%v", err.Error())
		return
	}

	for _, f := range fs {
		fmt.Println(f)
	}

	t, _ := ext.TimeNowCST()
	fmt.Println(t)
}
```