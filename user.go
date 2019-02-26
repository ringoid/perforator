package main

import (
	"math/rand"
	"time"
	"sync"
)

func createUserJob(sex string, tokenChan chan<- string, mutex *sync.Mutex) {
	token := CreateUserMethod(sex, 3+rand.Intn(2), mutex)
	tokenChan <- token
}

func userJob(token string, finish chan<- int) {
	lastActionTime := int64(0)
	nfProfiles := NewFacesMethod(token, lastActionTime)
	if len(nfProfiles) == 0 {
		ERROR.Fatalf("0 new faces returned first time")
	}
	DEBUG.Printf("receive [%d] profiles for new faces feed first time", len(nfProfiles))
	for {

		//-------------- New Faces -------------//
		//lets view and like in new faces
		actions := make([]Action, 0)
		for _, each := range nfProfiles {
			actions = ProfileActionMethod(actions, each, ViewActionType, NewFacesSourceFeed)
			actions = ProfileActionMethod(actions, each, LikeActionType, NewFacesSourceFeed)
		}
		//send actions
		if tmpL := ActionMethod(token, actions); tmpL > 0 {
			lastActionTime = tmpL
		}
		DEBUG.Printf("successfully send [%d] actions from new faces feed, lastActionTime [%v]", len(actions), lastActionTime)
		time.Sleep(time.Millisecond * 500)

		//--------------Who liked me -------------//
		//now go lmm feed and like somebody who liked me before
		likes, matches, messages := LMMMethod(token, lastActionTime)
		DEBUG.Printf("receive [%d] profiles in lmm (likes you) feed", len(likes))
		actions = make([]Action, 0)
		for _, each := range likes {
			actions = ProfileActionMethod(actions, each, ViewActionType, WhoLikedMeSourceFeed)
			actions = ProfileActionMethod(actions, each, LikeActionType, WhoLikedMeSourceFeed)
		}
		//send actions
		if tmpL := ActionMethod(token, actions); tmpL > 0 {
			DEBUG.Printf("successfully send [%d] actions from who liked me feed, lastActionTime [%v]", len(actions), lastActionTime)
			lastActionTime = tmpL
			time.Sleep(time.Millisecond * 500)
			likes, matches, messages = LMMMethod(token, lastActionTime)
		} else {
			DEBUG.Printf("there is no actions from who liked me feed, lastActionTime [%v]", lastActionTime)
		}

		//-------------- Matches  -------------//
		DEBUG.Printf("receive [%d] profiles in lmm (matches) feed", len(matches))
		actions = make([]Action, 0)
		for _, each := range matches {
			actions = ProfileActionMethod(actions, each, ViewActionType, MatchesSourceFeed)
			actions = ProfileActionMethod(actions, each, MessageActionType, MatchesSourceFeed)
		}
		//send actions
		if tmpL := ActionMethod(token, actions); tmpL > 0 {
			DEBUG.Printf("successfully send [%d] actions from matches feed, lastActionTime [%v]", len(actions), lastActionTime)
			lastActionTime = tmpL
			time.Sleep(time.Millisecond * 500)
			likes, matches, messages = LMMMethod(token, lastActionTime)
		} else {
			DEBUG.Printf("there is no actions from who liked me feed, lastActionTime [%v]", lastActionTime)
		}

		//-------------- Messages  -------------//
		DEBUG.Printf("receive [%d] profiles in lmm (messages) feed", len(messages))
		actions = make([]Action, 0)
		for _, each := range messages {
			actions = ProfileActionMethod(actions, each, ViewActionType, MessagesSourceFeed)
			actions = ProfileActionMethod(actions, each, MessageActionType, MessagesSourceFeed)
		}
		//send actions
		if tmpL := ActionMethod(token, actions); tmpL > 0 {
			lastActionTime = tmpL
		}
		DEBUG.Printf("successfully send [%d] actions from messages feed, lastActionTime [%v]", len(actions), lastActionTime)
		time.Sleep(time.Millisecond * 500)

		//-------------- New Faces -------------//
		nfProfiles = NewFacesMethod(token, lastActionTime)
		DEBUG.Printf("receive [%d] profiles for new faces feed", len(nfProfiles))
		if len(messages) == concurrentUsers/2 {
			INFO.Printf("user has a chat with everyone, terminate")
			finish <- 1
			return
		}
	}
}
