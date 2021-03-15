package security

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"github.com/go-basic/uuid"
	"io/ioutil"
	"os"
	"strings"
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

func RsaDecryptByPrivateKey(cipherText string) (string, error) {
	pwd, error := os.Getwd()
	if error != nil {
		return "", error
	}
	priKey, error := readRsaPriKey(pwd + "\\config\\rsa_1024_priv.pem")
	if error != nil {
		return "", error
	}

	decrypted, error := base64.StdEncoding.DecodeString(cipherText)
	if error != nil {
		return "", error
	}

	//解密
	info, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, []byte(decrypted))
	if err != nil {
		return "", err
	}
	return string(info), nil
}

func MD5(text string, salt string) string {
	ctx := md5.New()
	ctx.Write([]byte(salt + text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func Uuid(separator bool) string {
	uuid:= uuid.New()
	if separator{
		return uuid
	}
	return strings.ReplaceAll(uuid, "-", "")
}
