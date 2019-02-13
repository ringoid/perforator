package main

import (
	"encoding/json"
)

func Actions(accessToken string, actions []Action) (bool, *ActionResponse) {
	request := ActionReq{
		AccessToken: accessToken,
		Actions:     actions,
	}

	jsonBody, err := json.Marshal(request)
	if err != nil {
		ERROR.Fatalf("error marshaling request [%v]", request)
	}

	ok, respBody := MakePostRequest(ActionsApiEndpoint, jsonBody, "/actions")
	if !ok {
		return ok, nil
	}

	response := ActionResponse{}
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
