package main

import (
	"crypto"
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"net/url"
	"time"
)
type IcbcClient struct {
	appid string //appid
	privateKey string //我方私钥
	signType string //签名方式
	charset string //字符集，仅支持UTF-8,可填空‘’
	format string //请求参数格式，仅支持json，可填空‘’
	icbcPulicKey string //工行公钥
	encryptKey string //加密key
	encryptType string //加密方式
}
/**
初始化icbc对象
*/
func (icbc *IcbcClient) New(appid string, privateKey string, signType string, charset string, format string, icbcPulicKey string,encryptKey string, encryptType string) error{//设置基础信息
	if appid=="" || privateKey=="" || signType=="" || icbcPulicKey=="" {
		return errors.New("some params can not be empty")
	}
	icbc.appid = appid
	icbc.privateKey = privateKey
	icbc.signType = signType
	icbc.charset = charset
	icbc.format = format
	icbc.icbcPulicKey = icbcPulicKey
	icbc.encryptKey = encryptKey
	icbc.encryptType = encryptType
	return nil
}

/**
请求参数预处理
 */
func (icbc *IcbcClient) prepareParams(request map[string]interface{}, msgId string, appAuthToken string) (map[string]interface{},error){
	//params to return
	params := map[string]interface{}{}
	//biz to json string
	bizContentStr, jmErr := json.Marshal(request["biz_content"])
	if jmErr!= nil {
		return params,jmErr
	}
	//prepare public params
	params[APP_ID] = icbc.appid
	params[SIGN_TYPE] = icbc.signType
	params[CHARSET] = icbc.charset
	params[FORMAT] = icbc.format
	params[MSG_ID] = msgId
	params[TIMESTAMP] = time.Now().Format("2006-01-02 15:04:05")
	params[BIZ_CONTENT_KEY] = string(bizContentStr)

	//get path
	path, gerr := url.Parse(request["serviceUrl"].(string))
	if gerr != nil {
		return params,gerr
	}
	//build sign string
	var signStr string
	BuildOrderedSignStr(path.Path, params , &signStr)
	//
	var signStrHad string
	sErr := Sign(signStr, icbc.signType, icbc.privateKey, icbc.charset , &signStrHad)
	if sErr!=nil {
		return nil,sErr
	}
	params[SIGN] = signStrHad
	return params,nil
}

/**
请求执行程序
 */
func (icbc *IcbcClient) execute(request map[string]interface{}, msgId string, auToken string) (string,error){
	params, perr := icbc.prepareParams(request, msgId ,auToken)
	if perr!=nil {
		return "",nil
	}
	//发送请求
	//接收响应
	var respStr string
	if request["method"] == "GET" {
		error := DoGet(request["serviceUrl"].(string),params,icbc.charset , &respStr)
		if error != nil {
			return EMPTY, error
		}
	} else if request["method"] == "POST" {
		error := DoPost(request["serviceUrl"].(string),params,icbc.charset, &respStr)
		if error != nil {
			return EMPTY, error
		}
	}else{
		return EMPTY,errors.New("Only support GET or POST http method!")
	}
	//解析json
	var jsonRes = gjson.GetMany(respStr, "response_biz_content","sign")
	var b error
	if SIGN_TYPE_RSA == icbc.signType {
		b =RsaVerifySign(jsonRes[0].String(),icbc.icbcPulicKey, crypto.SHA1, jsonRes[1].String())
	}else if SIGN_TYPE_RSA2 == icbc.signType {
		b =RsaVerifySign(jsonRes[0].String(),icbc.icbcPulicKey, crypto.SHA256, jsonRes[1].String())
	}else{
		return EMPTY,errors.New("Only support RSA signature!in respose")
	}
	//b :=RsaVerifySign(jsonRes[0].String(),icbc.icbcPulicKey, crypto.SHA1, jsonRes[1].String())
	if b!=nil {
		return respStr,b
	}
	return respStr,nil
}

