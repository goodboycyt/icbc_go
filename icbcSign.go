package main

import (
	"crypto"
	"errors"
)

func Sign(strToSign string, signType string, privateKey string, charset string , signStr *string) error {
	var err error
	err = nil
	if SIGN_TYPE_RSA == signType {
		*signStr,err = RsaSign(strToSign, privateKey, crypto.SHA1)
	}else if SIGN_TYPE_RSA2 == signType {
		*signStr,err = RsaSign(strToSign, privateKey, crypto.SHA256)
	}else{
		err = errors.New("Only support RSA signature!")
	}
	return err
}

func verify()  {//no use

}
