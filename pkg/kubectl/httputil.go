package cmdline

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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

	return "", nil
}
