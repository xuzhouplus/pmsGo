package security

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

func MD5(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func ReadRsaPubKey(filename string) (*rsa.PublicKey, error) {
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
	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, err
	}

	//断言失败
	pubkey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("pubkey no ok")
	}
	return pubkey, nil
}

//读取私钥文件
func ReadRsaPriKey(filename string) (*rsa.PrivateKey, error) {
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
	pri, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, err
	}
	//断言失败
	prikey, ok := pri.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("prikey no ok")
	}
	return prikey, nil
}

//加密
func RsaEncryptData(filename string, src []byte) ([]byte, error) {
	//获取公钥
	pubKey, err := ReadRsaPubKey(filename)
	if err != nil {
		return nil, err
	}

	//加密
	encryptInfo, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, src)
	if err != nil {
		return nil, err
	}
	return encryptInfo, nil
}

//解密
func RsaDecryptData(filename string, src []byte) ([]byte, error) {
	//获取私钥
	priKey, err := ReadRsaPriKey(filename)
	if err != nil {
		return nil, err
	}

	//解密
	info, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, src)
	if err != nil {
		return nil, err
	}
	return info, nil
}
