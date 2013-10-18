//生成签名用
package openqq

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

//生成签名
func MakeSig(method string, url_path string, params map[string]string, secret string) string {
	hash := hmac.New(sha1.New, []byte(secret))
	ms := MakeSource(method, url_path, params)
	hash.Write([]byte(ms))
	val := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(val)
}

//生成签名所需源串
func MakeSource(method string, url_path string, params map[string]string) string {
	keys := []string{}
	for k, _ := range params {
		keys = append(keys, k)
	}
	//sort
	sort.Strings(keys)
	var allbuffer bytes.Buffer
	allbuffer.WriteString(strings.ToUpper(method))
	allbuffer.WriteString("&")
	allbuffer.WriteString(url.QueryEscape(url_path))
	allbuffer.WriteString("&")
	var buffer bytes.Buffer
	klen := len(keys)
	for i := 0; i < klen; i += 1 {
		buffer.WriteString(keys[i])
		buffer.WriteString("=")
		buffer.WriteString(params[keys[i]])
		if i != klen-1 {
			buffer.WriteString("&")
		}
	}
	allbuffer.WriteString(url.QueryEscape(buffer.String()))
	return allbuffer.String()
}

//验证签名
func VerifySig(method string, url_path string, params map[string]string, secret string, sig string) bool {
	// 确保不含sig
	delete(params, "sig")

	// 按照发货回调接口的编码规则对value编码
	CodePayVal(params)
	// 计算签名
	newsig := MakeSig(method, url_path, params, secret)
	return newsig == sig

}

//应用发货URL接口对腾讯回调传来的参数value值先进行一次编码方法，用于验签
//(编码规则为：除了 0~9 a~z A~Z !*() 之外其他字符按其ASCII码的十六进制加%进行表示，例如“-”编码为“%2D”)
//参考 <回调发货URL的协议说明_V3>
func CodePayVal(params map[string]string) {
	for k, v := range params {
		params[k] = EncodeValue(v)
	}
}

//应用发货URL接口的编码规则
func EncodeValue(s string) string {
	rexp := "[0-9a-zA-Z!*\\(\\)]"
	var buffer bytes.Buffer
	for i := 0; i < len(s); i += 1 {
		tempstr := string(s[i])
		b, err := regexp.MatchString(rexp, tempstr)
		if err == nil && b {
			buffer.WriteString(tempstr)
		} else {
			buffer.WriteString(hexString(tempstr))
		}

	}
	return buffer.String()
}

//应用发货URL　十六进制编码
func hexString(s string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(s); i++ {
		hex := hex.EncodeToString([]byte{s[i]})
		if len(hex) == 1 {
			buffer.WriteString("0")
		}
		buffer.WriteString("%")
		buffer.WriteString(hex)
	}
	return buffer.String()
}
