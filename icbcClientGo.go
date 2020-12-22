/**
工行api请求
*/
package icbc_go

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"net/url"
	"strings"
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
func (icbc *IcbcClient) prepareParams(request *map[string]interface{}, msgId string, appAuthToken string) (map[string]interface{},error){
	//params to return
	params := map[string]interface{}{}
	//biz to json string
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	eerr := jsonEncoder.Encode((*request)["biz_content"])
	if eerr!=nil {
		return nil,eerr
	}

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
	path, gerr := url.Parse((*request)["serviceUrl"].(string))
	if gerr != nil {
		return params,gerr
	}
	//if encrypt
	if (*request)["isNeedEncrypt"].(bool) == true && bizContentStr!=""{
		if icbc.encryptType != "AES" {
			return params,errors.New("only support aes")
		}
		params[ENCRYPT_TYPE] = icbc.encryptType
		//var aeserr error
		aesKey, aeErr := base64.StdEncoding.DecodeString(icbc.encryptKey)
		if aeErr != nil {
			return params,aeErr
		}
		var aeserr error
		params[BIZ_CONTENT_KEY], aeserr = EncryptByAes([]byte(bizContentStr), aesKey)
		if aeserr!=nil {
			return params,aeserr
		}
		//params[BIZ_CONTENT_KEY] = base64.StdEncoding.EncodeToString(tmp)

	} else {
		params[BIZ_CONTENT_KEY] = bizContentStr
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
func (icbc *IcbcClient) Execute(request *map[string]interface{}, msgId string, auToken string) (string,error){
	params, perr := icbc.prepareParams(request, msgId ,auToken)
	if perr!=nil {
		return "",perr
	}
	//发送请求
	//接收响应
	var respStr string
	if (*request)["method"] == "GET" {
		error := DoGet((*request)["serviceUrl"].(string),params,icbc.charset , &respStr)
		if error != nil {
			return "", error
		}
	} else if (*request)["method"] == "POST" {
		error := DoPost((*request)["serviceUrl"].(string),params,icbc.charset, &respStr)
		if error != nil {
			return "", error
		}
	}else{
		return "",errors.New("Only support GET or POST http method!")
	}
	//解析json
	var jsonRes = gjson.GetMany(respStr, "response_biz_content","sign")
	//var b error
	b := error(nil)
	if SIGN_TYPE_RSA == icbc.signType {
		b =RsaVerifySign(jsonRes[0].String(),icbc.icbcPulicKey, crypto.SHA1, jsonRes[1].String())
	}else if SIGN_TYPE_RSA2 == icbc.signType {
		b =RsaVerifySign(jsonRes[0].String(),icbc.icbcPulicKey, crypto.SHA256, jsonRes[1].String())
	}else{
		b = errors.New("Only support RSA signature!in respose")
	}
	if (*request)["isNeedEncrypt"].(bool) == true {
		if icbc.encryptType != "AES" {
			return respStr,errors.New("only support aes;reponse")
		}
		//params[ENCRYPT_TYPE] = icbc.encryptType
		//var aeserr error
		aesKey, aeErr := base64.StdEncoding.DecodeString(icbc.encryptKey)
		if aeErr != nil {
			return respStr,aeErr
		}
		bzByte, bzErr := base64.StdEncoding.DecodeString(jsonRes[0].String())
		if bzErr != nil {
			return respStr,bzErr
		}
		tmp, aeserr := AesDecrypt(bzByte, aesKey)
		if aeserr!=nil {
			return respStr,aeserr
		}
		respStr = string(tmp)

	}
	return respStr,b
}

