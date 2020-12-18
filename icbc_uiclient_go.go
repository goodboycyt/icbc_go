/**
工行ui请求
 */
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"
)
type IcbcClientUi struct {
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
func (icbc *IcbcClientUi) New(appid string, privateKey string, signType string, charset string, format string, icbcPulicKey string,encryptKey string, encryptType string) error{//设置基础信息
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
url query params build
 */
func (icbc *IcbcClientUi) buildUrlQueryParams(params map[string]interface{} , urlQueryParams *map[string]interface{}, urlBodyParams *map[string]interface{}) {
	apiParamNames := make(map[string]bool)
	apiParamNames[SIGN] = true
	apiParamNames[APP_ID] = true
	apiParamNames[SIGN_TYPE] = true
	apiParamNames[CHARSET] = true
	apiParamNames[FORMAT] = true
	apiParamNames[ENCRYPT_TYPE] = true
	apiParamNames[TIMESTAMP] = true
	apiParamNames[MSG_ID] = true
	for k,v := range params {
		if _,ok := apiParamNames[k];ok {
			(*urlQueryParams)[k] = v
		} else {
			(*urlBodyParams)[k] = v
		}
	}
}

/**
build url
 */
func (icbc *IcbcClientUi) BuildPostForm(request map[string]interface{}, msgId string, appAuthToken string) (string,error) {
	params, perr := icbc.prepareParams(request, msgId ,appAuthToken)
	if perr!=nil {
		return "",nil
	}
	urlQueryParams := map[string]interface{}{}
	urlBodyParams :=  map[string]interface{}{}
	icbc.buildUrlQueryParams(params, &urlQueryParams, &urlBodyParams)
	url := BuildGetUrl(request["serviceUrl"].(string),urlQueryParams,icbc.charset)
	return BuildForm(url,urlBodyParams),nil
}
/**
请求参数预处理
*/
func (icbc *IcbcClientUi) prepareParams(request map[string]interface{}, msgId string, appAuthToken string) (map[string]interface{},error){
	//params to return
	params := map[string]interface{}{}
	//biz to json string
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(request["biz_content"])

	//bizContentStr := bf.String()
	bizContentStr := strings.Replace(bf.String(), "/", "\\/", -1)
	bizContentStr = strings.TrimRight(bizContentStr, "\n")
	//prepare public params
	params[APP_ID] = icbc.appid
	params[SIGN_TYPE] = icbc.signType
	params[CHARSET] = icbc.charset
	params[FORMAT] = icbc.format
	params[MSG_ID] = msgId
	params[TIMESTAMP] = time.Now().Format("2006-01-02 15:04:05")
	params[BIZ_CONTENT_KEY] = bizContentStr
	//get path
	path, gerr := url.Parse(request["serviceUrl"].(string))
	if gerr != nil {
		return params,gerr
	}
	//build sign string

	var signStr string
	BuildOrderedSignStr(path.Path, params , &signStr)
	var signStrHad string

	sErr := Sign(signStr, icbc.signType, icbc.privateKey, icbc.charset , &signStrHad)

	if sErr!=nil {
		return nil,sErr
	}
	params[SIGN] = signStrHad
	return params,nil
}