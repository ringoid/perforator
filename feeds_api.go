package main

import (
	"fmt"
	"encoding/json"
)

func GetNewFaces(accessToken, resolution string, lastActionTime int64) (bool, *GetNewFacesFeedResp) {
	params := make(map[string]string)
	params["accessToken"] = accessToken
	params["resolution"] = resolution
	params["lastActionTime"] = fmt.Sprintf("%v", lastActionTime)
	ok, respBody := MakeGetRequest(FeedsApiEndpoint, params, "/get_new_faces")
	if !ok {
		return ok, nil
	}
	response := GetNewFacesFeedResp{}
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

func GetLMM(accessToken, resolution string, lastActionTime int64) (bool, *LMMFeedResp) {
	params := make(map[string]string)
	params["accessToken"] = accessToken
	params["resolution"] = resolution
	params["lastActionTime"] = fmt.Sprintf("%v", lastActionTime)
	ok, respBody := MakeGetRequest(FeedsApiEndpoint, params, "/get_lmm")
	if !ok {
		return ok, nil
	}

	response := LMMFeedResp{}
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
