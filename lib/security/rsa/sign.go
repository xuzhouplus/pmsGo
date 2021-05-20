package rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

const (
	kPublicKeyPrefix = "-----BEGIN PUBLIC KEY-----"
	kPublicKeySuffix = "-----END PUBLIC KEY-----"

	kPKCS1Prefix = "-----BEGIN RSA PRIVATE KEY-----"
	KPKCS1Suffix = "-----END RSA PRIVATE KEY-----"

	kPKCS8Prefix = "-----BEGIN PRIVATE KEY-----"
	KPKCS8Suffix = "-----END PRIVATE KEY-----"
)

var (
	ErrPrivateKeyFailedToLoad = errors.New("private key failed to load")
	ErrPublicKeyFailedToLoad  = errors.New("public key failed to load")
)

func formatKey(raw, prefix, suffix string, lineCount int) []byte {
	if raw == "" {
		return nil
	}
	raw = strings.Replace(raw, prefix, "", 1)
	raw = strings.Replace(raw, suffix, "", 1)
	raw = strings.Replace(raw, " ", "", -1)
	raw = strings.Replace(raw, "\n", "", -1)
	raw = strings.Replace(raw, "\r", "", -1)
	raw = strings.Replace(raw, "\t", "", -1)

	var sl = len(raw)
	var c = sl / lineCount
	if sl%lineCount > 0 {
		c = c + 1
	}

	var buf bytes.Buffer
	buf.WriteString(prefix + "\n")
	for i := 0; i < c; i++ {
		var b = i * lineCount
		var e = b + lineCount
		if e > sl {
			buf.WriteString(raw[b:])
		} else {
			buf.WriteString(raw[b:e])
		}
		buf.WriteString("\n")
	}
	buf.WriteString(suffix)
	return buf.Bytes()
}

func FormatPublicKey(raw string) []byte {
	return formatKey(raw, kPublicKeyPrefix, kPublicKeySuffix, 64)
}

func FormatPKCS1PrivateKey(raw string) []byte {
	raw = strings.Replace(raw, kPKCS8Prefix, "", 1)
	raw = strings.Replace(raw, KPKCS8Suffix, "", 1)
	return formatKey(raw, kPKCS1Prefix, KPKCS1Suffix, 64)
}

func FormatPKCS8PrivateKey(raw string) []byte {
	raw = strings.Replace(raw, kPKCS1Prefix, "", 1)
	raw = strings.Replace(raw, KPKCS1Suffix, "", 1)
	return formatKey(raw, kPKCS8Prefix, KPKCS8Suffix, 64)
}

func ParsePKCS1PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrPrivateKeyFailedToLoad
	}

	key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, err
}

func ParsePKCS8PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrPrivateKeyFailedToLoad
	}

	rawKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := rawKey.(*rsa.PrivateKey)
	if ok == false {
		return nil, ErrPrivateKeyFailedToLoad
	}

	return key, err
}

func ParsePublicKey(data []byte) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrPublicKeyFailedToLoad
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, ErrPublicKeyFailedToLoad
	}

	return key, err
}

func packageData(originalData []byte, packageSize int) (r [][]byte) {
	var src = make([]byte, len(originalData))
	copy(src, originalData)

	r = make([][]byte, 0)
	if len(src) <= packageSize {
		return append(r, src)
	}
	for len(src) > 0 {
		var p = src[:packageSize]
		r = append(r, p)
		src = src[packageSize:]
		if len(src) <= packageSize {
			r = append(r, src)
			break
		}
	}
	return r
}

// RSAEncrypt 使用公钥 key 对数据 src 进行 RSA 加密
func RSAEncrypt(src, key []byte) ([]byte, error) {
	pub, err := ParsePublicKey(key)
	if err != nil {
		return nil, err
	}

	return RSAEncryptWithKey(src, pub)
}

// RSAEncryptWithKey 使用公钥 key 对数据 src 进行 RSA 加密
func RSAEncryptWithKey(src []byte, key *rsa.PublicKey) ([]byte, error) {
	var data = packageData(src, key.N.BitLen()/8-11)
	var cipher = make([]byte, 0, 0)

	for _, d := range data {
		var c, e = rsa.EncryptPKCS1v15(rand.Reader, key, d)
		if e != nil {
			return nil, e
		}
		cipher = append(cipher, c...)
	}

	return cipher, nil
}

// RSADecryptWithPKCS1 使用私钥 key 对数据 cipher 进行 RSA 解密，key 的格式为 pkcs1
func RSADecryptWithPKCS1(cipher, key []byte) ([]byte, error) {
	pri, err := ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}

	return RSADecryptWithKey(cipher, pri)
}

// RSADecryptWithPKCS1 使用私钥 key 对数据 cipher 进行 RSA 解密，key 的格式为 pkcs8
func RSADecryptWithPKCS8(cipher, key []byte) ([]byte, error) {
	pri, err := ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	return RSADecryptWithKey(cipher, pri)
}

// RSADecryptWithKey 使用私钥 key 对数据 cipher 进行 RSA 解密
func RSADecryptWithKey(cipher []byte, key *rsa.PrivateKey) ([]byte, error) {
	var data = packageData(cipher, key.PublicKey.N.BitLen()/8)
	var plainData = make([]byte, 0, 0)

	for _, d := range data {
		var p, e = rsa.DecryptPKCS1v15(rand.Reader, key, d)
		if e != nil {
			return nil, e
		}
		plainData = append(plainData, p...)
	}
	return plainData, nil
}

func RSASignWithPKCS1(src, key []byte, hash crypto.Hash) ([]byte, error) {
	pri, err := ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return RSASignWithKey(src, pri, hash)
}

func RSASignWithPKCS8(src, key []byte, hash crypto.Hash) ([]byte, error) {
	pri, err := ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return RSASignWithKey(src, pri, hash)
}

func RSASignWithKey(src []byte, key *rsa.PrivateKey, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}

func RSAVerify(src, sig, key []byte, hash crypto.Hash) error {
	pub, err := ParsePublicKey(key)
	if err != nil {
		return err
	}
	return RSAVerifyWithKey(src, sig, pub, hash)
}

func RSAVerifyWithKey(src, sig []byte, key *rsa.PublicKey, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.VerifyPKCS1v15(key, hash, hashed, sig)
}
