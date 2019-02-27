package main

import (
	"expvar"
	"github.com/zserge/metric"
	"time"
	"sync"
	"os"
	"fmt"
	"strconv"
)

var concurrentUsers int

func main() {
	var err error
	args := os.Args[1:]
	env := args[0]
	concurrentUsers, err = strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("error input args conversion")
		os.Exit(-1)
	}

	InitConfig(env)
	INFO.Printf("Start application in [%s] env with [%d] concurrent users", env, concurrentUsers)

	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				printResults()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	mutex := sync.Mutex{}
	tokenChan := make(chan string)
	tokens := make([]string, 0)
	for i := 0; i < concurrentUsers; i++ {
		sex := "male"
		if i%2 == 0 {
			sex = "female"
		}
		go createUserJob(sex, tokenChan, &mutex)
	}

	//receive tokens from the chanel
	for i := 0; i < concurrentUsers; i++ {
		DEBUG.Printf("start waiting token")
		token := <-tokenChan
		tokens = append(tokens, token)
		DEBUG.Printf("received token")
	}

	DEBUG.Printf("start waiting 30 sec")

	time.Sleep(time.Second * 30)

	DEBUG.Printf("run user lifecycle")
	//run user lifecycle
	finishUserJobChan := make(chan int)
	for _, each := range tokens {
		go userJob(each, finishUserJobChan)
	}

	//just wait while user's new faces feed will be empty for everybody
	for i := 0; i < concurrentUsers; i++ {
		_ = <-finishUserJobChan
	}

	quit <- struct{}{}
	time.Sleep(time.Second)

	printResults()
}

func printResults() {
	INFO.Printf("%s counter : %v", USER_COUNTER, expvar.Get(USER_COUNTER).(metric.Metric))
	INFO.Printf("%s counter : %v", PHOTO_COUNTER, expvar.Get(PHOTO_COUNTER).(metric.Metric))
	INFO.Printf("%s histogram : %v", CREATE_USER_TIME, expvar.Get(CREATE_USER_TIME).(metric.Metric))
	INFO.Printf("%s counter : %v", NEW_FACES_REQUEST_COUNTER, expvar.Get(NEW_FACES_REQUEST_COUNTER).(metric.Metric))
	INFO.Printf("%s histogram : %v", NEW_FACES_RESPONSE_TIME, expvar.Get(NEW_FACES_RESPONSE_TIME).(metric.Metric))
	INFO.Printf("%s counter : %v", LMM_REQUEST_COUNTER, expvar.Get(LMM_REQUEST_COUNTER).(metric.Metric))
	INFO.Printf("%s histogram : %v", LMM_RESPONSE_TIME, expvar.Get(LMM_RESPONSE_TIME).(metric.Metric))
	INFO.Printf("%s counter : %v", SUBMITED_ACTIONS_COUNTER, expvar.Get(SUBMITED_ACTIONS_COUNTER).(metric.Metric))
	INFO.Printf("%s counter : %v", ACTION_REQUEST_COUNTER, expvar.Get(ACTION_REQUEST_COUNTER).(metric.Metric))
	INFO.Printf("%s histogram : %v", ACTION_REQUEST_TIME, expvar.Get(ACTION_REQUEST_TIME).(metric.Metric))

	INFO.Printf("%s counter : %v", SUCCESSFULLY_NEW_FACES_REQUEST_COUNTER, expvar.Get(SUCCESSFULLY_NEW_FACES_REQUEST_COUNTER).(metric.Metric))
	INFO.Printf("%s counter : %v", FAILED_NEW_FACES_REQUEST_COUNTER, expvar.Get(FAILED_NEW_FACES_REQUEST_COUNTER).(metric.Metric))
	INFO.Printf("%s counter : %v", SUCCESSFULLY_LMM_REQUEST_COUNTER, expvar.Get(SUCCESSFULLY_LMM_REQUEST_COUNTER).(metric.Metric))
	INFO.Printf("%s counter : %v", FAILED_LMM_REQUEST_COUNTER, expvar.Get(FAILED_LMM_REQUEST_COUNTER).(metric.Metric))

}
