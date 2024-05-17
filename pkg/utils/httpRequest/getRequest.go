package httpRequest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
)

func GetRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error in http get")
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body), nil
}

func GetRequestWithParams(url string, params map[string]string) (string, error) {
	parseURL, err := neturl.Parse(url)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	p := neturl.Values{}
	for k, v := range params {
		p.Set(k, v)
	}

	parseURL.RawQuery = p.Encode()
	urlWithParams := parseURL.String()

	return GetRequest(urlWithParams)
}

// GetRequestByObject
//
//	 @param url
//	 @param param_list 		:	params list (can be nil)
//	 @param response_target 	:	the **point** of response object
//									you can also get plain text through a string point
//	 @return int
//	 @return error
func GetRequestByObject(url string, param_list map[string]string, response_target interface{}) (int, error) {
	parseURL, err := neturl.Parse(url)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	p := neturl.Values{}
	if param_list != nil {
		for k, v := range param_list {
			p.Set(k, v)
		}
	}

	parseURL.RawQuery = p.Encode()
	urlWithParams := parseURL.String()
	response, err := http.Get(urlWithParams)

	if err != nil {
		fmt.Printf("GetRequestByObject: Get object failed, err msg: %s\n" + err.Error())
		return response.StatusCode, err
	}
	defer response.Body.Close()

	// 解析返回数据，存入response_target
	// 如果response_target为nil,即不需要返回数据只要状态码，则直接返回
	if response_target == nil {
		return response.StatusCode, nil
	}

	va, ok := response_target.(*string)
	if ok {
		// 如果response_target为*string,则按照文本方式读取
		text, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("GetRequestByObject: Decode response failed, err msg: %s\n", err.Error())
			return response.StatusCode, err
		}
		*va = string(text)
	} else {
		// 如果response_target为其他类型指针，则decode json为该类型结构体
		err = json.NewDecoder(response.Body).Decode(response_target)
		if err != nil {
			fmt.Printf("GetRequestByObject: Decode response failed, err msg: %s\n", err.Error())
			return response.StatusCode, err
		}
	}

	return response.StatusCode, nil
}
