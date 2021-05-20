package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"pmsGo/lib/config"
)

func readRsaPubKey(filename string) (*rsa.PublicKey, error) {
	//1. 读取公钥文件
	info, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	//2. 解码，得到block
	block, _ := pem.Decode(info)
	//3. 得到der
	der := block.Bytes
	//4. 得到公钥
	pub, err := x509.ParsePKCS1PublicKey(der)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

func readRsaPriKey(filename string) (*rsa.PrivateKey, error) {
	//1. 读取私钥文件
	info, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	//2. 解码，得到block
	block, _ := pem.Decode(info)
	//3. 得到der
	der := block.Bytes
	//4. 得到私钥
	pri, err := x509.ParsePKCS1PrivateKey(der)
	if err != nil {
		return nil, err
	}
	return pri, nil
}

func decodeCipher(cipherText string, primaryKey *rsa.PrivateKey) ([]byte, error) {
	decrypted, error := base64.StdEncoding.DecodeString(cipherText)
	if error != nil {
		return nil, error
	}
	//解密
	info, err := rsa.DecryptPKCS1v15(rand.Reader, primaryKey, decrypted)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func DecryptByPrivateKey(cipherText interface{}) (string, error) {
	pwd, error := os.Getwd()
	if error != nil {
		return "", error
	}
	priKey, error := readRsaPriKey(pwd + string(filepath.Separator) + config.Config.Web.Security["primaryKey"])
	if error != nil {
		return "", error
	}
	switch cipherText.(type) {
	case []string:
		plainText := make([]byte, 0)
		cipherTextSections := cipherText.([]string)
		for _, section := range cipherTextSections {
			cipher, err := decodeCipher(section, priKey)
			if err != nil {
				return "", err
			}
			plainText = append(plainText, cipher...)
		}
		return string(plainText), nil
	case string:
		cipher, err := decodeCipher(cipherText.(string), priKey)
		if err != nil {
			return "", err
		}
		return string(cipher), nil
	default:
		return "", errors.New("数据只能是[]string或string")
	}
}
