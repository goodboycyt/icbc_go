## icbc go lib
封装的功能包括，自动加签，验签，加密，解密
aesutil.go:aes加解密文件
icbc_client_go.go:api请求文件
icbc_uiclient_go.go:ui请求文件
icbcSign.go :签名验签文件
myconst.go:常量文件
rsa.go:rsa相关文件
webUtils.go:请求相关


### use mod config
`require github.com/goodboycyt/icbc_go v1.4.4`
or
`go get -u github.com/goodboycyt/icbc_client_go`


### useage
```go
    var icbc IcbcClient
	icbc.New("1211111", "ddsad+sadsa+RMWK3Ci+sad+YaeH/Qm/r/Topq3lABw==","RSA","UTF-8","json","MIGfMA0GCSqGSIb3DQEBwIDAQAB","","")
	request_b := map[string]interface{}{"serviceUrl":"https://url","method":"POST","isNeedEncrypt":false,"extraParams":""}
	request_b["biz_content"] = map[string]interface{}{"corp_no":"123213","trx_acc_date":"2020-12-14"}
	resP,err :=icbc.execute(&request_b, "202012241521929252" , "")
	if err!=nil {
        fmt.printLn(err)
    }   


    var icbc IcbcClientUi
    icbc.New("11", "=","RSA","UTF-8","json","MIB","","")
    request_b := map[string]interface{}{"serviceUrl":"https://1.1.com.cn/ui/1/ui/1/1/1/V1","method":"POST","isNeedEncrypt":false,"extraParams":""}
    request_b["biz_content"] = map[string]interface{}{"121":"12121"}
    resP,err :=icbc.BuildPostForm(&request_b, "202012241521929252" , "")
    if err!=nil {
            fmt.printLn(err)
    } 
```
