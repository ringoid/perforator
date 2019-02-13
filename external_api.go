package main

import (
	"github.com/ringoid/commons"
	"fmt"
)

//------------------------------- Auth Service ----------------------------------------//
type CreateReq struct {
	WarmUpRequest              bool   `json:"warmUpRequest"`
	YearOfBirth                int    `json:"yearOfBirth"`
	Sex                        string `json:"sex"`
	Locale                     string `json:"locale"`
	DateTimeTermsAndConditions int64  `json:"dtTC"`
	DateTimePrivacyNotes       int64  `json:"dtPN"`
	DateTimeLegalAge           int64  `json:"dtLA"`
	DeviceModel                string `json:"deviceModel"`
	OsVersion                  string `json:"osVersion"`
}

func (req CreateReq) String() string {
	return fmt.Sprintf("%#v", req)
}

type CreateResp struct {
	commons.BaseResponse
	AccessToken string `json:"accessToken"`
	CustomerId  string `json:"customerId"`
}

func (resp CreateResp) String() string {
	return fmt.Sprintf("%#v", resp)
}

type DeleteReq struct {
	WarmUpRequest bool   `json:"warmUpRequest"`
	AccessToken   string `json:"accessToken"`
}

func (req DeleteReq) String() string {
	return fmt.Sprintf("%#v", req)
}

type UpdateSettingsReq struct {
	WarmUpRequest       bool   `json:"warmUpRequest"`
	AccessToken         string `json:"accessToken"`
	SafeDistanceInMeter int    `json:"safeDistanceInMeter"` // 0 (default for men) || 10 (default for women)
	PushMessages        bool   `json:"pushMessages"`        // true (default for men) || false (default for women)
	PushMatches         bool   `json:"pushMatches"`         // true (default)
	PushLikes           string `json:"pushLikes"`           //EVERY (default for men) || 10_NEW (default for women) || 100_NEW || NONE
}

func (req UpdateSettingsReq) String() string {
	return fmt.Sprintf("%#v", req)
}

type GetSettingsResp struct {
	commons.BaseResponse
	SafeDistanceInMeter int    `json:"safeDistanceInMeter"` // 0 (default for men) || 25 (default for women)
	PushMessages        bool   `json:"pushMessages"`        // true (default for men) || false (default for women)
	PushMatches         bool   `json:"pushMatches"`         // true (default)
	PushLikes           string `json:"pushLikes"`           //EVERY (default for men) || 10_NEW (default for women) || 100_NEW || NONE
}

func (resp GetSettingsResp) String() string {
	return fmt.Sprintf("%#v", resp)
}

//------------------------------- Image Service ----------------------------------------//
type GetPresignUrlReq struct {
	WarmUpRequest bool   `json:"warmUpRequest"`
	AccessToken   string `json:"accessToken"`
	Extension     string `json:"extension"`
	ClientPhotoId string `json:"clientPhotoId"`
}

func (req GetPresignUrlReq) String() string {
	return fmt.Sprintf("%#v", req)
}

type GetPresignUrlResp struct {
	commons.BaseResponse
	Uri           string `json:"uri"`
	OriginPhotoId string `json:"originPhotoId"`
	ClientPhotoId string `json:"clientPhotoId"`
}

func (resp GetPresignUrlResp) GoString() string {
	return fmt.Sprintf("%#v", resp)
}

type GetOwnPhotosResp struct {
	commons.BaseResponse
	Photos []OwnPhoto `json:"photos"`
}

func (resp GetOwnPhotosResp) String() string {
	return fmt.Sprintf("%#v", resp)
}

type OwnPhoto struct {
	PhotoId       string `json:"photoId"`
	PhotoUri      string `json:"photoUri"`
	Likes         int    `json:"likes"`
	OriginPhotoId string `json:"originPhotoId"`
	Blocked       bool   `json:"blocked"`
}

func (obj OwnPhoto) String() string {
	return fmt.Sprintf("%#v", obj)
}

type DeletePhotoReq struct {
	WarmUpRequest bool   `json:"warmUpRequest"`
	AccessToken   string `json:"accessToken"`
	PhotoId       string `json:"photoId"`
}

func (req DeletePhotoReq) String() string {
	return fmt.Sprintf("%#v", req)
}

//------------------------------- Actions Service ----------------------------------------//
type ActionReq struct {
	AccessToken string   `json:"accessToken"`
	Actions     []Action `json:"actions"`
}

func (req ActionReq) String() string {
	return fmt.Sprintf("%#v", req)
}

type Action struct {
	SourceFeed         string `json:"sourceFeed"`
	ActionType         string `json:"actionType"`
	TargetPhotoId      string `json:"targetPhotoId"`
	TargetUserId       string `json:"targetUserId"`
	Text               string `json:"text"`
	LikeCount          int    `json:"likeCount"`
	ViewCount          int    `json:"viewCount"`
	ViewTimeMillis     int64  `json:"viewTimeMillis"`
	OpenChatCount      int    `json:"openChatCount"`
	OpenChatTimeMillis int64  `json:"openChatTimeMillis"`
	BlockReasonNum     int    `json:"blockReasonNum"`
	ActionTime         int64  `json:"actionTime"`
}

func (req Action) String() string {
	return fmt.Sprintf("%#v", req)
}

type ActionResponse struct {
	commons.BaseResponse
	LastActionTime int64 `json:"lastActionTime"`
}

func (resp ActionResponse) String() string {
	return fmt.Sprintf("%#v", resp)
}

//------------------------------- Feeds Service ----------------------------------------//

type GetNewFacesFeedResp struct {
	commons.BaseResponse
	Profiles              []commons.Profile `json:"profiles"`
	RepeatRequestAfterSec int               `json:"repeatRequestAfterSec"`
}

func (resp GetNewFacesFeedResp) String() string {
	return fmt.Sprintf("%#v", resp)
}

type LMMFeedResp struct {
	commons.BaseResponse
	LikesYou              []commons.Profile `json:"likesYou"`
	Matches               []commons.Profile `json:"matches"`
	Messages              []commons.Profile `json:"messages"`
	RepeatRequestAfterSec int               `json:"repeatRequestAfterSec"`
}

func (resp LMMFeedResp) String() string {
	return fmt.Sprintf("%#v", resp)
}
