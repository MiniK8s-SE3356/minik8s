package cmdline

import (
	"bytes"
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
