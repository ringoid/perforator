package main

import (
	"time"
	"expvar"
	"github.com/zserge/metric"
	"github.com/ringoid/commons"
	"sync"
)

const (
	PhotoResolution480x640   = "480x640"
	PhotoResolution720x960   = "720x960"
	PhotoResolution1080x1440 = "1080x1440"
	PhotoResolution1440x1920 = "1440x1920"
)

const (
	NewFacesSourceFeed   = "new_faces"
	WhoLikedMeSourceFeed = "who_liked_me"
	MatchesSourceFeed    = "matches"
	MessagesSourceFeed   = "messages"

	ViewActionType    = "VIEW"
	LikeActionType    = "LIKE"
	UnlikeActionType  = "UNLIKE"
	BlockActionType   = "BLOCK"
	MessageActionType = "MESSAGE"
)

func NewFacesMethod(token string, lastActionTime int64, wasActionSent bool) []commons.Profile {
	startTime := time.Now().Round(time.Millisecond).UnixNano() / 1000000
	for {
		ok, resp := GetNewFaces(token, PhotoResolution1080x1440, lastActionTime)
		if ok && resp.RepeatRequestAfter == 0 {
			finishTime := time.Now().Round(time.Millisecond).UnixNano() / 1000000
			expvar.Get(SUCCESSFULLY_NEW_FACES_REQUEST_COUNTER).(metric.Metric).Add(1)
			expvar.Get(NEW_FACES_REQUEST_COUNTER).(metric.Metric).Add(1)
			expvar.Get(NEW_FACES_RESPONSE_TIME).(metric.Metric).Add(float64(finishTime - startTime))
			DEBUG.Printf("new faces request was successfull with profiles num [%d]", len(resp.Profiles))
			if wasActionSent {
				expvar.Get(NEW_FACES_AFTER_ACTION_REQUEST_COUNTER).(metric.Metric).Add(1)
				expvar.Get(NEW_FACES_AFTER_ACTION_RESPONSE_TIME).(metric.Metric).Add(float64(finishTime - startTime))
			}
			return resp.Profiles
		}
		if ok && resp != nil && resp.RepeatRequestAfter != 0 {
			DEBUG.Printf("new faces return repeat after sec [%v]", resp.RepeatRequestAfter)
			//time.Sleep(time.Millisecond * time.Duration(resp.RepeatRequestAfter))
			time.Sleep(time.Millisecond * 100)
			expvar.Get(SUCCESSFULLY_NEW_FACES_REQUEST_COUNTER).(metric.Metric).Add(1)
		}
		if !ok {
			DEBUG.Printf("failed new faces request, drop counter and repeat")
			expvar.Get(FAILED_NEW_FACES_REQUEST_COUNTER).(metric.Metric).Add(1)
			startTime = time.Now().Round(time.Millisecond).UnixNano() / 1000000
			time.Sleep(time.Millisecond * 500)
		}
	}
}

//return likes, matched and messages
func LMMMethod(token string, lastActionTime int64, wasActionSent bool) ([]commons.Profile, []commons.Profile, []commons.Profile) {
	startTime := time.Now().Round(time.Millisecond).UnixNano() / 1000000
	for {
		ok, resp := GetLMM(token, PhotoResolution1440x1920, lastActionTime)
		if ok && resp.RepeatRequestAfter == 0 {
			finishTime := time.Now().Round(time.Millisecond).UnixNano() / 1000000
			expvar.Get(SUCCESSFULLY_LMM_REQUEST_COUNTER).(metric.Metric).Add(1)
			expvar.Get(LMM_REQUEST_COUNTER).(metric.Metric).Add(1)
			expvar.Get(LMM_RESPONSE_TIME).(metric.Metric).Add(float64(finishTime - startTime))
			DEBUG.Printf("lmm request was successfull with likes num [%d], matches num [%d], messages num [%d]",
				len(resp.LikesYou), len(resp.Matches), len(resp.Messages))
			if wasActionSent {
				expvar.Get(LMM_AFTER_ACTION_REQUEST_COUNTER).(metric.Metric).Add(1)
				expvar.Get(LMM_RESPONSE_AFTER_ACTION_TIME).(metric.Metric).Add(float64(finishTime - startTime))
			}
			return resp.LikesYou, resp.Matches, resp.Messages
		}
		if resp != nil && resp.RepeatRequestAfter != 0 {
			DEBUG.Printf("lmm return repeat after sec [%v]", resp.RepeatRequestAfter)
			expvar.Get(SUCCESSFULLY_LMM_REQUEST_COUNTER).(metric.Metric).Add(1)
			time.Sleep(time.Millisecond * 100)
			//time.Sleep(time.Millisecond * time.Duration(resp.RepeatRequestAfter))
		}
		if !ok {
			DEBUG.Printf("failed llm request, drop counter and repeat")
			expvar.Get(FAILED_LMM_REQUEST_COUNTER).(metric.Metric).Add(1)
			time.Sleep(time.Millisecond * 500)
			startTime = time.Now().Round(time.Millisecond).UnixNano() / 1000000
		}
	}
}

func ActionMethod(token string, sourceActions []Action) int64 {
	if len(sourceActions) == 0 {
		DEBUG.Printf("empty action list, return")
		return -1
	}

	startTime := time.Now().Round(time.Millisecond).UnixNano() / 1000000
	for {
		ok, resp := Actions(token, sourceActions)
		if ok {
			finishTime := time.Now().Round(time.Millisecond).UnixNano() / 1000000
			expvar.Get(SUBMITED_ACTIONS_COUNTER).(metric.Metric).Add(float64(len(sourceActions)))
			expvar.Get(ACTION_REQUEST_COUNTER).(metric.Metric).Add(1)
			expvar.Get(ACTION_REQUEST_TIME).(metric.Metric).Add(float64(finishTime - startTime))
			DEBUG.Printf("actions request was successfull with actions num [%d]", len(sourceActions))
			return resp.LastActionTime
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func ProfileActionMethod(sourceActions []Action, profile commons.Profile, actionType, source string) []Action {
	if sourceActions == nil {
		sourceActions = make([]Action, 0)
	}
	for _, eachPhoto := range profile.Photos {
		switch actionType {
		case ViewActionType:
			sourceActions = append(sourceActions, Action{
				SourceFeed:     source,
				ActionType:     ViewActionType,
				TargetPhotoId:  eachPhoto.PhotoId,
				TargetUserId:   profile.UserId,
				LikeCount:      0,
				ViewCount:      1,
				ViewTimeMillis: 1,
				ActionTime:     time.Now().Round(time.Millisecond).UnixNano() / 1000000,
			})
		case LikeActionType:
			sourceActions = append(sourceActions, Action{
				SourceFeed:     source,
				ActionType:     LikeActionType,
				TargetPhotoId:  eachPhoto.PhotoId,
				TargetUserId:   profile.UserId,
				LikeCount:      1,
				ViewCount:      0,
				ViewTimeMillis: 0,
				ActionTime:     time.Now().Round(time.Millisecond).UnixNano() / 1000000,
			})
		case MessageActionType:
			sourceActions = append(sourceActions, Action{
				SourceFeed:    source,
				ActionType:    MessageActionType,
				TargetPhotoId: eachPhoto.PhotoId,
				TargetUserId:  profile.UserId,
				Text:          "hello ass!",
				ActionTime:    time.Now().Round(time.Millisecond).UnixNano() / 1000000,
			})
		case UnlikeActionType:
			sourceActions = append(sourceActions, Action{
				SourceFeed:    source,
				ActionType:    UnlikeActionType,
				TargetPhotoId: eachPhoto.PhotoId,
				TargetUserId:  profile.UserId,
				ActionTime:    time.Now().Round(time.Millisecond).UnixNano() / 1000000,
			})
		case BlockActionType:
			sourceActions = append(sourceActions, Action{
				SourceFeed:     source,
				ActionType:     BlockActionType,
				TargetPhotoId:  eachPhoto.PhotoId,
				TargetUserId:   profile.UserId,
				BlockReasonNum: 7,
				ActionTime:     time.Now().Round(time.Millisecond).UnixNano() / 1000000,
			})
		}
	}
	return sourceActions
}

func CreateUserMethod(sex string, photoNum int, mutex *sync.Mutex) string {
	startTime := time.Now().Round(time.Millisecond).UnixNano() / 1000000
	var ok bool
	var accessToken string

	for {
		ok, createProfileResp := CreateUserProfile(1982, sex)
		if ok {
			accessToken = createProfileResp.AccessToken
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	DEBUG.Printf("user profile was created")
	//upload images
	for i := 0; i < photoNum; i++ {
		var presignUri string
		for {
			ok, getPresignResp := GetPresignUrl(accessToken)
			if ok {
				presignUri = getPresignResp.Uri
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		DEBUG.Printf("presign link for photo [%d] was received", i)
		var image []byte
		for {
			isItDog := sex == "male"
			ok, image = GenerateCatOrDog(isItDog, mutex)
			if ok {
				break
			}
		}
		DEBUG.Printf("image for photo [%d] was generated", i)
		MakePutRequestWithContent(presignUri, image)
		DEBUG.Printf("image for photo [%d] was uploaded", i)
		expvar.Get(PHOTO_COUNTER).(metric.Metric).Add(1)
	} //end images upload

	//todo:later add something related to settings

	finishTime := time.Now().Round(time.Millisecond).UnixNano() / 1000000

	expvar.Get(USER_COUNTER).(metric.Metric).Add(1)
	expvar.Get(CREATE_USER_TIME).(metric.Metric).Add(float64(finishTime - startTime))
	return accessToken
}
