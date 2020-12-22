/**
rsa 处理程序
 */
package icbc_go

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)
//加签
func RsaSign(signContent string, privateKey string, hash crypto.Hash) (string,error) {
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)
	priKey, err := ParsePrivateKey(privateKey)
	if err != nil {
		return "",err
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, hash, hashed)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(signature),nil
}
//验签
func RsaVerifySign(signContent string, publicKey string, hash crypto.Hash, sign string) error {
	sign1, err1 := base64.StdEncoding.DecodeString(sign)
	if err1 != nil {
		return err1
	}
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)
	pubKey, err := ParsePublicKey(publicKey)
	if err != nil {
		return err
	}
	signature := rsa.VerifyPKCS1v15( pubKey, hash, hashed, sign1)
	if signature != nil {
		return signature
	}
	return nil
}
//私钥转换
func ParsePrivateKey(privateKey string)(*rsa.PrivateKey, error) {
	privateKey = FormatPrivateKey(privateKey)
	// 2、解码私钥字节，生成加密对象
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("私钥信息错误！")
	}
	// 3、解析DER编码的私钥，生成私钥对象
	//priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	priKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priKey.(*rsa.PrivateKey), nil
}
//私钥格式化
func FormatPrivateKey(privateKey string) string  {
	if !strings.HasPrefix(privateKey, PEM_BEGIN) {
		privateKey = PEM_BEGIN + privateKey
	}
	if !strings.HasSuffix(privateKey, PEM_END) {
		privateKey = privateKey + PEM_END
	}
	return privateKey
}
//公钥转换
func ParsePublicKey(pulicKey string)(*rsa.PublicKey, error) {
	pulicKey = FormatPublicKey(pulicKey)
	// 2、解码私钥字节，生成加密对象
	block, _ := pem.Decode([]byte(pulicKey))
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	// 3、解析DER编码的私钥，生成私钥对象
	//priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}
//公钥格式化
func FormatPublicKey(pulicKey string) string  {
	if !strings.HasPrefix(pulicKey, PPEM_BEGIN) {
		pulicKey = PPEM_BEGIN + pulicKey
	}
	if !strings.HasSuffix(pulicKey, PPEM_END) {
		pulicKey = pulicKey + PPEM_END
	}
	return pulicKey
}