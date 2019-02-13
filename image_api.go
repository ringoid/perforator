package main

import (
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/json"
	"github.com/ringoid/commons"
	"bytes"
	"sync"
)

var catImage []byte
var dogImage []byte

func MakePutRequestWithContent(url string, source []byte) {
	request, err := http.NewRequest("PUT", url, bytes.NewReader(source))
	if err != nil {
		ERROR.Fatalf("error create PUT request by [%s] : %v", url, err)
	}

	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		ERROR.Fatalf("error execute PUT request by [%s] : %v", url, err)
	}
	defer httpResponse.Body.Close()

	_, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		ERROR.Fatalf("error reading response by [%s] : %v", url, err)
	}

	if httpResponse.StatusCode != 200 {
		ERROR.Fatalf("not OK response from [%s], status code is [%d]", url, httpResponse.StatusCode)
	}
}

//Image service
func GenerateCatOrDog(isItDog bool, mutex *sync.Mutex) (bool, []byte) {
	if isItDog && len(dogImage) > 0 {
		return true, dogImage
	} else if !isItDog && len(catImage) > 0 {
		return true, catImage
	}

	mutex.Lock()
	defer mutex.Unlock()

	if isItDog && len(dogImage) > 0 {
		return true, dogImage
	} else if !isItDog && len(catImage) > 0 {
		return true, catImage
	}
	var result []byte
	for {
		DEBUG.Printf("Try to upload [isItDog=%v] image first time", isItDog)
		urlStr := "https://api.thecatapi.com/v1/images/search"
		if isItDog {
			urlStr = "https://dog.ceo/api/breeds/image/random"
		}
		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			ERROR.Fatalf("error create GET request to [%s] : %v", urlStr, err)
		}
		client := &http.Client{}
		httpResponse, err := client.Do(req)
		if err != nil {
			ERROR.Fatalf("error execute GET request to [%s] : %v", urlStr, err)
		}
		defer httpResponse.Body.Close()

		if httpResponse.StatusCode != 200 {
			ERROR.Printf("not OK response from [%s], status code is [%d]", urlStr, httpResponse.StatusCode)
			return false, nil
		}

		respBody, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			ERROR.Fatalf("error reading body response from [%s] : %v", urlStr, err)
		}

		var finalUrl string
		if isItDog {
			var dogResp DogResponse
			err = json.Unmarshal(respBody, &dogResp)
			if err != nil {
				ERROR.Fatalf("error unmarshal response body [%s] : %v", string(respBody), err)
			}
			if !strings.HasSuffix(dogResp.Message, ".jpg") {
				continue
			}

			finalUrl = dogResp.Message

		} else {
			var arrResp []CatResponse
			err = json.Unmarshal(respBody, &arrResp)
			if err != nil {
				ERROR.Fatalf("error unmarshal response body [%s] : %v", string(respBody), err)
			}
			if len(arrResp) != 1 {
				ERROR.Printf("error call generate image, 0 image returned")
				return false, nil
			}

			if !strings.HasSuffix(arrResp[0].Url, ".jpg") {
				continue
			}

			finalUrl = arrResp[0].Url
		}

		req, err = http.NewRequest("GET", finalUrl, nil)
		if err != nil {
			ERROR.Fatalf("error create GET request to final url [%s] : %v", finalUrl, err)
		}
		httpResponse2, err := client.Do(req)
		if err != nil {
			ERROR.Fatalf("error execute GET request to final url [%s] : %v", finalUrl, err)
		}
		defer httpResponse2.Body.Close()

		if httpResponse2.StatusCode != 200 {
			ERROR.Printf("not OK response from final url [%s], status code is [%d]", finalUrl, httpResponse2.StatusCode)
			return false, nil
		}

		respBody, err = ioutil.ReadAll(httpResponse2.Body)
		if err != nil {
			ERROR.Fatalf("error reading response body from final url [%s] : %v", finalUrl, err)
		}
		result = respBody
		break
	}

	if isItDog {
		dogImage = result
	} else {
		catImage = result
	}
	return true, result
}

type CatResponse struct {
	Url string `json:"url"`
}

type DogResponse struct {
	Message string `json:"message"`
}

func GetPresignUrl(accessToken string) (bool, *GetPresignUrlResp) {
	request := GetPresignUrlReq{
		AccessToken:   accessToken,
		Extension:     "jpg",
		ClientPhotoId: "fakeClientId",
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		ERROR.Fatalf("error marshaling request [%v] : %v", request, err)
	}

	ok, respBody := MakePostRequest(ImageApiEndpoint, jsonBody, "/get_presigned")
	if !ok {
		return ok, nil
	}
	response := GetPresignUrlResp{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		ERROR.Fatalf("error unmarshaling response [%s] : %v", string(respBody), err)
	}

	if len(response.ErrorCode) != 0 {
		ERROR.Printf("error code not nil in reponse [%v]", response)
		return false, nil
	}

	return true, &response
}

func GetOwnPhotos(accessToken, resolution string) (bool, *GetOwnPhotosResp) {
	params := make(map[string]string)
	params["accessToken"] = accessToken
	params["resolution"] = resolution
	ok, respBody := MakeGetRequest(ImageApiEndpoint, params, "/get_own_photos")
	if !ok {
		return ok, nil
	}

	response := GetOwnPhotosResp{}
	err := json.Unmarshal(respBody, &response)
	if err != nil {
		ERROR.Fatalf("error unmarshaling response [%s] : %v", string(respBody), err)
	}

	if len(response.ErrorCode) != 0 {
		ERROR.Printf("error code not nil in reponse [%v]", response)
		return false, nil
	}

	return true, &response
}

func DeletePhoto(accessToken, photoId string) (bool, *commons.BaseResponse) {
	request := DeletePhotoReq{
		AccessToken: accessToken,
		PhotoId:     photoId,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		ERROR.Fatalf("error marshaling request [%v]", request)
	}

	ok, respBody := MakePostRequest(ImageApiEndpoint, jsonBody, "/delete_photo")
	if !ok {
		return ok, nil
	}

	response := commons.BaseResponse{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		ERROR.Fatalf("error unmarshaling response [%s] : %v", string(respBody), err)
	}

	if len(response.ErrorCode) != 0 {
		ERROR.Printf("error code not nil in reponse [%v]", response)
		return false, nil
	}

	return true, &response
}
