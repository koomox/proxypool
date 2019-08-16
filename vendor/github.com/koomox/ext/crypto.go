package ext

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"io/ioutil"
)

const (
	saltSize                    = 12
	nonceSize                   = 12
	saltNonceSize               = saltSize + nonceSize
	keySize                     = 32
	versionCrypto        uint32 = 1
	versionSize                 = 4
	versionSaltSize             = versionSize + saltSize
	versionSaltNonceSize        = saltNonceSize + versionSize
	cryptoLineLength            = 1024
)

type cryptoGCM struct {
	salt   []byte
	nonce  []byte
	secret []byte
	key    []byte
	aesgcm cipher.AEAD
}

func NewNonce() []byte {
	nonce := make([]byte, nonceSize)
	rand.Read(nonce)
	return nonce
}

func NewKey(salt, password []byte) ([]byte, error) {
	key, err := scrypt.Key(password, salt, 16384, 8, 1, keySize)

	if err != nil {
		return nil, fmt.Errorf("Crypto New Key Failed! %v", err.Error())
	}

	return key, err
}

func NewGCM(salt, nonce, secret []byte) (*cryptoGCM, error) {
	key, err := NewKey(salt, secret)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	gcm := &cryptoGCM{
		salt:   salt,
		nonce:  nonce,
		secret: secret,
		key:    key,
		aesgcm: aesgcm,
	}

	return gcm, err
}

func (this *cryptoGCM) Decrypt(ciphertext []byte) ([]byte, error) {

	plaintext, err := this.aesgcm.Open(nil, this.nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (this *cryptoGCM) Encrypt(cleartext []byte) ([]byte, error) {

	var ciphertext []byte

	ciphertext = this.aesgcm.Seal(nil, this.nonce, cleartext, nil)

	return ciphertext, nil
}

func EncryptFile(src, dst string, secret []byte) (err error) {
	r, err := EncryptFileToRaw(src, secret)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(dst, r, 0775); err != nil {
		return fmt.Errorf("EncrytFile write DestFile Err:%v", err.Error())
	}
	return nil
}

func EncryptFileToRaw(src string, secret []byte) ([]byte, error) {
	var (
		ciphertext []byte
		clearData  []byte
		cipherData []byte
		err        error
	)
	salt := NewNonce()
	nonce := NewNonce()

	gcm, err := NewGCM(salt, nonce, secret)

	if err != nil {
		return nil, fmt.Errorf("EncrytFileToRaw NewGCM Err:%v", err.Error())
	}

	exist, err := PathExist(src)
	if !exist || err != nil {
		return nil, fmt.Errorf("EncrytFileToRaw Not Found input file Err:%v", err.Error())
	}

	if clearData, err = ioutil.ReadFile(src); err != nil {
		return nil, fmt.Errorf("EncrytFileToRaw ReadFile Err:%v", err.Error())
	}

	if ciphertext, err = gcm.Encrypt(clearData); err != nil {
		return nil, fmt.Errorf("EncrytFileToRaw encode... Err:%v", err.Error())
	}

	cipherData = versionedJoin(salt, nonce, ciphertext)
	return cipherData, nil
}

func DecryptFile(src, dst string, secret []byte) (err error) {
	r, err := DecryptFileToRaw(src, secret)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(dst, r, 0775); err != nil {
		return fmt.Errorf("DecryptFile write DestFile Err:%v", err.Error())
	}
	return nil
}

func DecryptFileToRaw(src string, secret []byte) ([]byte, error) {
	var (
		cipherData []byte
		ciphertext []byte
		cleartext  []byte
		salt       []byte
		nonce      []byte
		err        error
		version    uint32
	)

	exist, err := PathExist(src)
	if !exist || err != nil {
		return nil, fmt.Errorf("DecryptFileToRaw Not Found input file Err:%v", err.Error())
	}

	if cipherData, err = ioutil.ReadFile(src); err != nil {
		return nil, fmt.Errorf("DecryptFileToRaw ReadFile Err:%v", err.Error())
	}

	if version, salt, nonce, ciphertext, err = versionedSplit(cipherData); err != nil {
		return nil, err
	}

	if version != versionCrypto {
		return nil, errors.New("DecryptFileToRaw decode Version failed!")
	}

	gcm, err := NewGCM(salt, nonce, secret)
	if err != nil {
		return nil, fmt.Errorf("DecryptFileToRaw NewGCM Err:%v", err.Error())
	}

	if cleartext, err = gcm.Decrypt(ciphertext); err != nil {
		return nil, fmt.Errorf("DecryptFileToRaw decode... Err:%v", err.Error())
	}

	return cleartext, nil
}

func versionedJoin(in ...[]byte) []byte {
	out := make([]byte, versionSize)
	binary.LittleEndian.PutUint32(out, versionCrypto)
	for _, args := range in {
		out = append(out, args...)
	}

	return out
}

func versionedSplit(in []byte) (version uint32, salt, nonce, ciphertext []byte, err error) {
	if len(in) < versionSaltNonceSize {
		return 0, nil, nil, nil, errors.New("Invalid byte length.")
	}

	version = binary.LittleEndian.Uint32(in[:versionSize])
	salt = in[versionSize:versionSaltSize]
	nonce = in[versionSaltSize:versionSaltNonceSize]
	ciphertext = in[versionSaltNonceSize:]

	return version, salt, nonce, ciphertext, nil
}
