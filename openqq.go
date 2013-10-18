// openqq project openqq.go
//腾讯开放平台(openqq) SDK for golang
package openqq

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//接口信息配置结构
type config struct {
	appid      string
	appkey     string
	serverName string
}

//初始化配置
func init() {
	currConfig.appid = "123456"
	currConfig.appkey = "228bf094169a40a3bd188ba37ebe8723"
	currConfig.serverName = "119.147.19.43"
}

var currConfig config

//API调用
//scriptName 为需要调用的api, 例如 /v3/user/get_info
//params 为需发送的参数
// protocol = https || http
func API(scriptName string, params map[string]string, protocol string) (map[string]interface{}, error) {
	delete(params, "sig")
	params["appid"] = currConfig.appid

	method := "POST"
	secuet := currConfig.appkey + "&"
	sig := MakeSig(method, scriptName, params, secuet)
	//fmt.Println("sig=", sig)
	params["sig"] = sig
	var buffer bytes.Buffer
	buffer.WriteString(protocol)
	buffer.WriteString("://")
	buffer.WriteString(currConfig.serverName)
	buffer.WriteString(scriptName)
	//fmt.Println("url=", buffer.String())
	return postRequest(buffer.String(), method, params, protocol)
}

//提交HTTP请求
func postRequest(callUrl string, method string, params map[string]string, protocol string) (map[string]interface{}, error) {
	var client *http.Client
	if strings.ToUpper(protocol) == "HTTPS" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = &http.Client{}
	}

	urlvs := url.Values{}
	for k, v := range params {
		urlvs.Add(k, v)
	}
	resp, err := client.PostForm(callUrl, urlvs)
	if err != nil {
		//fmt.Println("error")
		return nil, err
	}
	defer resp.Body.Close()
	body, tmpErr := ioutil.ReadAll(resp.Body)
	if tmpErr != nil {
		return nil, tmpErr
	}
	result := make(map[string]interface{})
	json.Unmarshal([]byte(body), &result)
	return result, nil

}
