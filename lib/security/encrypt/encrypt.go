package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/hkdf"
	"hash"
	"io"
)

func HashHmac(algo func() hash.Hash, data []byte, key []byte, rawHash bool) []byte {
	hash := hmac.New(algo, key)
	hash.Write(data)
	rawData := hash.Sum(nil)
	if rawHash {
		return rawData
	}
	hexData := hex.EncodeToString(rawData)
	return []byte(hexData)
}
func HashHkdf(algo func() hash.Hash, secret []byte, salt []byte, info []byte, length int) []byte {
	keyHkdf := hkdf.New(algo, secret, salt, info)
	key := make([]byte, length)
	io.ReadFull(keyHkdf, key)
	return key
}
func validateData(algo func() hash.Hash, data []byte, key []byte) []byte {
	hash := hmac.New(algo, []byte(""))
	hash.Write([]byte(""))
	hashLen := hash.BlockSize()
	dataLen := len(data)
	if dataLen >= hashLen {
		hashData := data[0:hashLen]
		pureData := data[hashLen:]
		calculatedHash := HashHmac(algo, pureData, key, false)
		if hmac.Equal(hashData, calculatedHash) {
			return pureData
		}
	}
	return nil
}
func randomBytes(length int) []byte {
	random := make([]byte, length)
	_, err := rand.Read(random)
	if err != nil {
		return nil
	}
	return random
}
func hashData(data []byte, key []byte, rawHash bool) []byte {
	hash := HashHmac(sha256.New, data, key, rawHash)
	return append(hash, data...)
}
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func PKCS7UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	if length-unpadding < 0 {
		return []byte("")
	}
	return src[:(length - unpadding)]
}
func opensslEncrypt(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                 // 获取秘钥块的长度
	origData := PKCS7Padding(data, blockSize)      // 补全码
	blockMode := cipher.NewCBCEncrypter(block, iv) // 加密模式
	encrypted := make([]byte, len(origData))       // 创建数组
	blockMode.CryptBlocks(encrypted, origData)     // 加密
	return encrypted, nil
}

func opensslDecrypt(data []byte, key []byte) ([]byte, error) {
	iv := data[0:16]
	encrypt := data[16:]
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cipher.NewCBCDecrypter(cipherBlock, iv).CryptBlocks(encrypt, encrypt)
	return encrypt, nil
}

func Encrypt(data []byte, salt []byte) ([]byte, error) {
	keySalt := randomBytes(16)
	key := HashHkdf(sha256.New, salt, keySalt, nil, 16)
	iv := randomBytes(16)
	encrypted, _ := opensslEncrypt(data, key, iv)
	authKey := HashHkdf(sha256.New, key, nil, []byte("AuthorizationKey"), 16)
	hashed := hashData(append(iv, encrypted...), authKey, false)
	d := append(keySalt, hashed...)
	return d, nil
}

func Decrypt(encrypted []byte, salt []byte) ([]byte, error) {
	keySalt := encrypted[0:16]
	key := HashHkdf(sha256.New, salt, keySalt, nil, 16)
	authKey := HashHkdf(sha256.New, key, nil, []byte("AuthorizationKey"), 16)
	data := validateData(sha256.New, encrypted[16:], authKey)
	if data == nil {
		return nil, errors.New("解析错误")
	}
	decrypt, err := opensslDecrypt(data, key)
	if err != nil {
		return nil, err
	}
	decrypt = PKCS7UnPadding(decrypt)
	return decrypt, nil
}
