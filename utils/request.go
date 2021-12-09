package utils

import (
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//Post 请求第三方接口
func Post(url string, data io.Reader) (body []byte, err error) {
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	return
}

// Request 发送请求
func Request(url, method string, postData string, headers ...map[string]string) ([]byte, error) {
	if method == "" {
		method = "GET"
	}

	var body io.Reader
	if postData != "" {
		body = strings.NewReader(postData)
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	request, _ := http.NewRequest(method, url, body)
	if len(headers) > 0 && headers[0] != nil {
		for key, val := range headers[0] {
			request.Header.Set(key, val)
		}
	}
	response, err := client.Do(request)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	return responseBody, err
}

// Data2json 转为json字符串
func Data2json(m interface{}) string {
	str, _ := jsoniter.MarshalToString(m)
	return str
}
