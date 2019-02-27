package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"expvar"
	"github.com/zserge/metric"
)

var AuthApiEndpoint string
var ImageApiEndpoint string
var ActionsApiEndpoint string
var FeedsApiEndpoint string

var INFO *log.Logger
var ERROR *log.Logger
var DEBUG *log.Logger

const (
	USER_COUNTER              = "UserCounter"
	CREATE_USER_TIME          = "CreateUserTime"
	PHOTO_COUNTER             = "PhotoCounter"
	NEW_FACES_REQUEST_COUNTER = "NewFacesRequestCounter"
	NEW_FACES_RESPONSE_TIME   = "NewFacesResponseTime"
	LMM_REQUEST_COUNTER       = "LmmRequestCounter"
	LMM_RESPONSE_TIME         = "LmmResponseTime"
	SUBMITED_ACTIONS_COUNTER  = "ActionsCounter"
	ACTION_REQUEST_COUNTER    = "ActionRequestCounter"
	ACTION_REQUEST_TIME       = "ActionRequestTime"

	SUCCESSFULLY_NEW_FACES_REQUEST_COUNTER = "SuccessfullyNewFacesRequestCounter"
	FAILED_NEW_FACES_REQUEST_COUNTER = "FailedNewFacesRequestCounter"

	SUCCESSFULLY_LMM_REQUEST_COUNTER = "SuccessfullyLmmRequestCounter"
	FAILED_LMM_REQUEST_COUNTER = "FailedLmmRequestCounter"
)

func InitConfig(env string) {
	AuthApiEndpoint = fmt.Sprintf("https://%s.ringoidapp.com/auth", env)
	ImageApiEndpoint = fmt.Sprintf("https://%s.ringoidapp.com/image", env)
	ActionsApiEndpoint = fmt.Sprintf("https://%s.ringoidapp.com/actions", env)
	FeedsApiEndpoint = fmt.Sprintf("https://%s.ringoidapp.com/feeds", env)

	file, err := os.OpenFile(fmt.Sprintf("%s-log.out", env), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error open file : %v", err)
		os.Exit(-1)
	}
	multi := io.MultiWriter(os.Stdout, file)
	INFO = log.New(multi, fmt.Sprintf("%s-INFO: ", env), log.Ldate|log.Ltime|log.Lshortfile)
	ERROR = log.New(multi, fmt.Sprintf("%s-ERROR: ", env), log.Ldate|log.Ltime|log.Lshortfile)
	DEBUG = log.New(os.Stdout, fmt.Sprintf("%s-DEBUG: ", env), log.Ldate|log.Ltime|log.Lshortfile)

	expvar.Publish(USER_COUNTER, metric.NewCounter())
	expvar.Publish(CREATE_USER_TIME, metric.NewHistogram())
	expvar.Publish(PHOTO_COUNTER, metric.NewCounter())
	expvar.Publish(NEW_FACES_REQUEST_COUNTER, metric.NewCounter())
	expvar.Publish(NEW_FACES_RESPONSE_TIME, metric.NewHistogram())
	expvar.Publish(LMM_REQUEST_COUNTER, metric.NewCounter())
	expvar.Publish(LMM_RESPONSE_TIME, metric.NewHistogram())
	expvar.Publish(SUBMITED_ACTIONS_COUNTER, metric.NewCounter())
	expvar.Publish(ACTION_REQUEST_COUNTER, metric.NewCounter())
	expvar.Publish(ACTION_REQUEST_TIME, metric.NewHistogram())

	expvar.Publish(SUCCESSFULLY_NEW_FACES_REQUEST_COUNTER, metric.NewCounter())
	expvar.Publish(FAILED_NEW_FACES_REQUEST_COUNTER, metric.NewCounter())
	expvar.Publish(SUCCESSFULLY_LMM_REQUEST_COUNTER, metric.NewCounter())
	expvar.Publish(FAILED_LMM_REQUEST_COUNTER, metric.NewCounter())
}
