package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

/**
generete string before to sign
*/
func BuildOrderedSignStr(path string, params map[string]interface{} , signStr *string) {
	var key []string
	//delete(params, "sign")
	for k := range params {
		key = append(key, k)
	}
	sort.Strings(key)

	builder := strings.Builder{}
	builder.WriteString(path+"?")
	for _, v := range key {
		builder.WriteString(v)
		builder.WriteString("=")
		switch params[v].(type) {
		case float64:
			builder.WriteString(fmt.Sprint(strconv.FormatFloat(params[v].(float64), 'f', -1, 64)))
			break
		default:
			builder.WriteString(fmt.Sprint(params[v]))
			break
		}
		builder.WriteString("&")

	}
	*signStr = builder.String()
	*signStr = (*signStr)[:len(*signStr)-1] //排序后去除尾部特殊字符
}

/**
get request
 */
func DoGet(serviceUrl string,params map[string]interface{},charset string, resStr *string) error{
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*2)    //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 30))    //设置发送接受数据超时
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	req, err := http.NewRequest(http.MethodGet, serviceUrl, nil)
	if err != nil {
		return err
	}
	// 添加请求头
	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset="+charset)
	req.Header.Add("APIGW-VERSION", "bg-go-v1")
	//加入get参数
	//values := url.Values{}
	q := req.URL.Query()
	for k, v := range params {
		switch v.(type) {
		case string:
			q.Add(k, v.(string))
		case int:
			q.Add(k, strconv.FormatInt(int64(v.(int)), 10))
		case int64:
			q.Add(k, strconv.FormatInt(v.(int64), 10))
		case float64:
			q.Add(k, strconv.FormatFloat(v.(float64), 'f', -1, 64))
		case float32:
			q.Add(k, strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64))
		}
	}

	req.URL.RawQuery = q.Encode()
	//fmt.Println("encode->", url.QueryEscape(q.Encode()))
	resp, derr := client.Do(req)
	if derr != nil {
		return derr
	}
	if resp.StatusCode != 200 {
		return errors.New("response status code is not valid. status code:"+string(resp.StatusCode))
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}

	*resStr = result.String()
	return nil
}

/**
post request
*/
func DoPost(serviceUrl string,params map[string]interface{},charset string, resStr *string) error{
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*1)    //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 30))    //设置发送接受数据超时
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 30,
		},

	}
	q := url.Values{}
	for k, v := range params {
		switch v.(type) {
		case string:
			q.Add(k, v.(string))
		case int:
			q.Add(k, strconv.FormatInt(int64(v.(int)), 10))
		case int64:
			q.Add(k, strconv.FormatInt(v.(int64), 10))
		case float64:
			q.Add(k, strconv.FormatFloat(v.(float64), 'f', -1, 64))
		case float32:
			q.Add(k, strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64))
		}
	}
	req, err := http.NewRequest(http.MethodPost, serviceUrl, strings.NewReader(q.Encode()))
	if err != nil {
		return err
	}
	// 添加请求头
	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset="+charset)
	req.Header.Add("APIGW-VERSION", "bg-go-v1")

	//加入get参数
	resp, derr := client.Do(req)
	if derr != nil {
		return derr
	}
	if resp.StatusCode != 200 {
		return errors.New("response status code is not valid. status code:"+string(resp.StatusCode))
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	*resStr = result.String()
	return nil
}

/*
build get url
 */
func BuildGetUrl(serviceUrl string, urlQueryParams map[string]interface{}, charset string) string {
	q := url.Values{}
	for k, v := range urlQueryParams {
		switch v.(type) {
		case string:
			q.Add(k, v.(string))
		case int:
			q.Add(k, strconv.FormatInt(int64(v.(int)), 10))
		case int64:
			q.Add(k, strconv.FormatInt(v.(int64), 10))
		case float64:
			q.Add(k, strconv.FormatFloat(v.(float64), 'f', -1, 64))
		case float32:
			q.Add(k, strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64))
		}
	}
	return serviceUrl+"?"+q.Encode()
}

/**
build form
 */
func BuildForm(url string, bodyParams map[string]interface{}) string {
	result :=""
	if bodyParams == nil || len(bodyParams) == 0 {
	} else {
		re3, _ := regexp.Compile("\"")
		fmt.Println(bodyParams)
		for k,v := range bodyParams {
			if v == "" || v==nil {

			} else {
				switch v.(type) {
				case string:
					result += "<input type=\"hidden\" name=\""+k+"\" value=\""+re3.ReplaceAllString(v.(string), "&quot;")+"\">\n"
					fmt.Println("result",result)
				case int:
					result += "<input type=\"hidden\" name=\""+k+"\" value=\""+strconv.FormatInt(int64(v.(int)), 10)+"\">\n"
				case int64:
					result += "<input type=\"hidden\" name=\""+k+"\" value=\""+strconv.FormatInt(v.(int64), 10)+"\">\n"
				case float64:
					result += "<input type=\"hidden\" name=\""+k+"\" value=\""+strconv.FormatFloat(v.(float64), 'f', -1, 64)+"\">\n"
				case float32:
					result += "<input type=\"hidden\" name=\""+k+"\" value=\""+strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64)+"\">\n"
				default:
					//result += "<input type=\"hidden\" name=\""+k+"\" value=\""+v.(string)+"\">\n"
				}
			}
			}
		}
	return "<form name=\"auto_submit_form\" method=\"post\" action=\""+url+"\">\n"+result+"<input type=\"submit\" value=\"立刻提交\" style=\"display:none\" >\n</form>\n<script>document.forms[0].submit();</script>"
}
