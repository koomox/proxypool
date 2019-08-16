package ext

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func GenPassword(pwd string) string {
	return GetMD5(GetSHA1(pwd))
}

func GetMD5(ctx string) string {
	h := md5.New()
	io.WriteString(h, ctx)
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func GetSHA1(s string) string {
	h := sha1.New()
	io.WriteString(h, s)
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func GetFileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("GetFileMD5 OpenFile Err:%v", err.Error())
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("GetFileMD5 io.Copy Err:%v", err.Error())
	}

	sum := h.Sum(nil)
	return hex.EncodeToString(sum), nil
}

func GetFileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("GetFileSHA256 Open File Err:%v", err.Error())
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("GetFileSHA256 io.Copy Err:%v", err.Error())
	}

	sum := h.Sum(nil)
	return hex.EncodeToString(sum), nil
}
