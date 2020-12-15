package main

import (
	"crypto"
	"errors"
)

func Sign(strToSign string, signType string, privateKey string, charset string , signStr *string) error {
	if SIGN_TYPE_RSA == signType {
		*signStr = RsaSign(strToSign, privateKey, crypto.SHA1)
	}else if SIGN_TYPE_RSA2 == signType {
		*signStr = RsaSign(strToSign, privateKey, crypto.SHA256)
	}else{
		return errors.New("Only support RSA signature!")
	}
	return nil
}

func Verify()  {

}
