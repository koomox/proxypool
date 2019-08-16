package ext

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	randomComplete = "`~^0OolI\"'/\\|"
	randomLetter   = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	randomNumber   = "123456789"
	randomSpecial  = "~!@#$%^&*_+:?`-=;,."
)

func generatorRandomStr(length int, complete, noNumber, noSpecial bool) (string, error) {
	var randomPool = randomLetter

	if complete {
		if noNumber || noSpecial {
			return "", fmt.Errorf("Cannot use `complete` flag with `no-number` and `no-special`.")
		}
		randomPool += randomNumber
		randomPool += randomSpecial
		randomPool += randomComplete
	} else {
		if !noNumber {
			randomPool += randomNumber
		}

		if !noSpecial {
			randomPool += randomSpecial
		}
	}

	randstr := make([]byte, length) // Random string to return
	charlen := big.NewInt(int64(len(randomPool)))
	for i := 0; i < length; i++ {
		b, err := rand.Int(rand.Reader, charlen)
		if err != nil {
			return "", fmt.Errorf("RandString Generator Err:%v", err.Error())
		}
		r := int(b.Int64())
		randstr[i] = randomPool[r]
	}
	return string(randstr), nil
}

func RandomString(length int) (result string, err error) {
	for i := 0; i < 3; i++ {
		if result, err = generatorRandomStr(length, false, false, true); err != nil {
			continue
		}
		return
	}
	return
}