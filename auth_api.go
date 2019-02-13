package main

import (
	"time"
	"encoding/json"
	"github.com/ringoid/commons"
)

//Auth service
func CreateUserProfile(yearOfBirth int, sex string) (bool, *CreateResp) {
	request := CreateReq{
		YearOfBirth:                yearOfBirth,
		Sex:                        sex,
		Locale:                     "ru",
		DateTimeLegalAge:           time.Now().Round(time.Millisecond).UnixNano() / 1000000,
		DateTimePrivacyNotes:       time.Now().Round(time.Millisecond).UnixNano() / 1000000,
		DateTimeTermsAndConditions: time.Now().Round(time.Millisecond).UnixNano() / 1000000,
		DeviceModel:                "test device",
		OsVersion:                  "test android",
	}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		ERROR.Fatalf("error marshaling request [%v]", request)
	}

	ok, respBody := MakePostRequest(AuthApiEndpoint, jsonBody, "/create_profile")
	if !ok {
		return ok, nil
	}

	response := CreateResp{}
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

func DeleteUserProfile(accessToken string, ) (bool, *commons.BaseResponse) {
	request := DeleteReq{
		AccessToken: accessToken,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		ERROR.Fatalf("error marshaling request [%v]", request)
	}

	ok, respBody := MakePostRequest(AuthApiEndpoint, jsonBody, "/delete")
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

func UpdateUserSettings(accessToken string, safeDistanceInMeter int, pushMessages, pushMatches bool,
	pushLikes string) (bool, *commons.BaseResponse) {

	request := UpdateSettingsReq{
		AccessToken:         accessToken,
		SafeDistanceInMeter: safeDistanceInMeter,
		PushMessages:        pushMessages,
		PushMatches:         pushMatches,
		PushLikes:           pushLikes,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		ERROR.Fatalf("error marshaling request [%v]", request)
	}

	ok, respBody := MakePostRequest(AuthApiEndpoint, jsonBody, "/update_settings")
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

func GetUserSettings(accessToken string) (bool, *GetSettingsResp) {
	params := make(map[string]string)
	params["accessToken"] = accessToken

	ok, respBody := MakeGetRequest(AuthApiEndpoint, params, "/get_settings")
	if !ok {
		return ok, nil
	}

	response := GetSettingsResp{}
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
