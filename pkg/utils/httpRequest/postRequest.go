package httpRequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func PostRequest(url string, jsonData []byte) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("error in post request")
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", body)

	return string(body), nil
}

// PostRequestByObject
//
//	 @param uri
//	 @param request_target : the port object
//	 @param response_target: the **point** of response object
//								you can also get plain text through a string point
//	 @return int
//	 @return error
func PostRequestByObject(uri string, request_target interface{}, response_target interface{}) (int, error) {
	jsonData, err := json.Marshal(request_target)
	if err != nil {
		fmt.Printf("PostRequestByObject: Marshal object failed, err msg: %s,n", err.Error())
		return 0, err
	}
	response, err := http.Post(uri, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("PostRequestByObject: Post object failed, err msg: %s\n" + err.Error())
		return 0, err
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
			fmt.Printf("PostRequestByObject: Decode response failed, err msg: %s\n", err.Error())
			return 0, err
		}
		*va = string(text)
	} else {
		// 如果response_target为其他类型指针，则decode json为该类型结构体
		err = json.NewDecoder(response.Body).Decode(response_target)
		if err != nil {
			fmt.Printf("PostRequestByObject: Decode response failed, err msg: %s\n", err.Error())
			return 0, err
		}
	}

	return response.StatusCode, nil
}
