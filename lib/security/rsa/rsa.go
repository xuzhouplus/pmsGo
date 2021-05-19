package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
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

func DecryptByPrivateKey(cipherText string) (string, error) {
	pwd, error := os.Getwd()
	if error != nil {
		return "", error
	}
	priKey, error := readRsaPriKey(pwd + string(filepath.Separator) + config.Config.Web.Security["primaryKey"])
	if error != nil {
		return "", error
	}

	decrypted, error := base64.StdEncoding.DecodeString(cipherText)
	if error != nil {
		return "", error
	}

	//解密
	info, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, decrypted)
	if err != nil {
		return "", err
	}
	return string(info), nil
}
