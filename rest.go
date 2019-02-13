package main

import (
	"strings"
	"net/http"
	"io/ioutil"
	"bytes"
	"fmt"
)

func MakeGetRequest(baseApi string, params map[string]string, urlPart string) (bool, []byte) {
	url := baseApi + urlPart + "?"
	for k, v := range params {
		url += fmt.Sprintf("%s=%s&", k, v)
	}
	index := strings.LastIndex(url, "&")
	url = url[:index]

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ERROR.Fatalf("error create new GET request to url [%s] : %v", url, err)
	}
	req.Header.Set("x-ringoid-android-buildnum", "1000")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	httpResponse, err := client.Do(req)
	if err != nil {
		ERROR.Fatalf("error execute new GET request to url [%s] : %v", url, err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		ERROR.Printf("not OK response from GET request to [%s], status code [%d]", url, httpResponse.StatusCode)
		return false, nil
	}

	respBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		ERROR.Printf("error reading response body, GET request to [%s] : %v", url, err)
		return false, nil
	}

	return true, respBody
}

func MakePostRequest(baseApi string, jsonBody []byte, urlPart string) (bool, []byte) {
	url := baseApi + urlPart
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		ERROR.Fatalf("error create new POST request to url [%s] with json body [%s] : %v", url, string(jsonBody), err)
	}
	req.Header.Set("x-ringoid-android-buildnum", "1000")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	httpResponse, err := client.Do(req)
	if err != nil {
		ERROR.Fatalf("error execute new POST request to url [%s] with json body [%s]: %v", url, string(jsonBody), err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		ERROR.Printf("not OK response from POST request to [%s] with json body [%s], status code [%d]", url, string(jsonBody), httpResponse.StatusCode)
		return false, nil
	}

	respBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		ERROR.Printf("error reading response body, POST request to [%s] with json body [%s] : %v", url, string(jsonBody), err)
		return false, nil
	}

	return true, respBody
}
