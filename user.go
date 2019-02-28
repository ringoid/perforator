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

func userJob(index int, token string, finish chan<- int) {
	wasActionsSent := false
	wait := 2 + rand.Intn(10)
	time.Sleep(time.Second * time.Duration(int64(wait)))

	lastActionTime := int64(0)
	nfProfiles := NewFacesMethod(token, lastActionTime, false)
	//if len(nfProfiles) == 0 {
	//	ERROR.Fatalf("0 new faces returned first time")
	//}
	DEBUG.Printf("[user%d] : receive [%d] profiles for new faces feed first time", index, len(nfProfiles))
	for {
		wasActionsSent = false
		//-------------- New Faces -------------//
		//lets view and like in new faces
		actions := make([]Action, 0)
		var tmpL int64
		for _, each := range nfProfiles {
			actions = ProfileActionMethod(actions, each, ViewActionType, NewFacesSourceFeed)
			actions = ProfileActionMethod(actions, each, LikeActionType, NewFacesSourceFeed)
			if len(actions) == 10 {
				time.Sleep(time.Second * 2)
				if tmpL = ActionMethod(token, actions); tmpL > 0 {
					lastActionTime = tmpL
					wasActionsSent = true
				}
				DEBUG.Printf("[user%d] : successfully send [%d] view and like actions from new faces feed, lastActionTime [%v]", index, len(actions), lastActionTime)
				actions = make([]Action, 0)
			}
		}
		//send actions
		if tmpL = ActionMethod(token, actions); tmpL > 0 {
			lastActionTime = tmpL
			wasActionsSent = true
			DEBUG.Printf("[user%d] : successfully send [%d] view and like actions from new faces feed, lastActionTime [%v]", index, len(actions), lastActionTime)
		}

		if !wasActionsSent {
			DEBUG.Printf("[user%d] : there is no actions from new faces feed, lastActionTime [%v]", index, lastActionTime)
		}

		anyLmmActionWasSend := false
		//--------------Who liked me -------------//
		//now go lmm feed and like somebody who liked me before
		likes, matches, messages := LMMMethod(token, lastActionTime, wasActionsSent)
		DEBUG.Printf("[user%d] : receive [%d] profiles in lmm (likes you) feed", index, len(likes))
		wasActionsSent = false
		actions = make([]Action, 0)
		for _, each := range likes {
			actions = ProfileActionMethod(actions, each, ViewActionType, WhoLikedMeSourceFeed)
			actions = ProfileActionMethod(actions, each, LikeActionType, WhoLikedMeSourceFeed)
			if len(actions) == 10 {
				time.Sleep(time.Second * 2)
				if tmpL = ActionMethod(token, actions); tmpL > 0 {
					lastActionTime = tmpL
					wasActionsSent = true
					anyLmmActionWasSend = true
				}
				actions = make([]Action, 0)
				DEBUG.Printf("[user%d] : successfully send [%d] view and like actions from who liked me feed, lastActionTime [%v]", index, len(actions), lastActionTime)
			}
		}
		//send actions
		if tmpL := ActionMethod(token, actions); tmpL > 0 {
			DEBUG.Printf("[user%d] : successfully send [%d] view and like actions from who liked me feed, lastActionTime [%v]", index, len(actions), lastActionTime)
			lastActionTime = tmpL
			wasActionsSent = true
			anyLmmActionWasSend = true
			//time.Sleep(time.Millisecond * 500)
		}
		if !wasActionsSent {
			DEBUG.Printf("[user%d] : there is no actions from who liked me feed, lastActionTime [%v]", index, lastActionTime)
		}

		if wasActionsSent{//only in case if sent something from likes (likes tab was not empty)
			likes, matches, messages = LMMMethod(token, lastActionTime, wasActionsSent)
		}

		wasActionsSent = false
		//-------------- Matches  -------------//
		DEBUG.Printf("[user%d] : receive [%d] profiles in lmm (matches) feed", index, len(matches))
		actions = make([]Action, 0)
		for _, each := range matches {
			actions = ProfileActionMethod(actions, each, ViewActionType, MatchesSourceFeed)
			actions = ProfileActionMethod(actions, each, MessageActionType, MatchesSourceFeed)
			if len(actions) == 10 {
				time.Sleep(time.Second * 2)
				if tmpL = ActionMethod(token, actions); tmpL > 0 {
					lastActionTime = tmpL
					wasActionsSent = true
					anyLmmActionWasSend = true
				}
				actions = make([]Action, 0)
				DEBUG.Printf("[user%d] : successfully send [%d] view and message actions from matches feed, lastActionTime [%v]", index, len(actions), lastActionTime)
			}
		}
		//send actions
		if tmpL := ActionMethod(token, actions); tmpL > 0 {
			DEBUG.Printf("[user%d] : successfully send [%d] view and message actions from matches feed, lastActionTime [%v]", index, len(actions), lastActionTime)
			lastActionTime = tmpL
			wasActionsSent = true
			anyLmmActionWasSend = true
			//time.Sleep(time.Millisecond * 500)
		}

		if !wasActionsSent {
			DEBUG.Printf("[user%d] : there is no actions from who liked me feed, lastActionTime [%v]", index, lastActionTime)
		}

		if !wasActionsSent{//ask only if we send something from matches
			likes, matches, messages = LMMMethod(token, lastActionTime, wasActionsSent)
		}

		wasActionsSent = false
		//-------------- Messages  -------------//
		DEBUG.Printf("[user%d] : receive [%d] profiles in lmm (messages) feed", index, len(messages))
		actions = make([]Action, 0)
		for _, each := range messages {
			actions = ProfileActionMethod(actions, each, ViewActionType, MessagesSourceFeed)
			actions = ProfileActionMethod(actions, each, MessageActionType, MessagesSourceFeed)
			if len(actions) == 10 {
				time.Sleep(time.Second * 2)
				if tmpL = ActionMethod(token, actions); tmpL > 0 {
					lastActionTime = tmpL
					wasActionsSent = true
					anyLmmActionWasSend = true
				}
				actions = make([]Action, 0)
				DEBUG.Printf("[user%d] : successfully send [%d] view and message actions from messages feed, lastActionTime [%v]", index, len(actions), lastActionTime)
			}
		}
		//send actions
		if tmpL := ActionMethod(token, actions); tmpL > 0 {
			lastActionTime = tmpL
			wasActionsSent = true
			anyLmmActionWasSend = true
			DEBUG.Printf("[user%d] : successfully send [%d] view and message actions from messages feed, lastActionTime [%v]", index, len(actions), lastActionTime)
		}

		if !wasActionsSent {
			DEBUG.Printf("[user%d] : there is no actions from messages feed, lastActionTime [%v]", index, lastActionTime)
		}
		//time.Sleep(time.Millisecond * 500)

		//-------------- New Faces -------------//
		nfProfiles = NewFacesMethod(token, lastActionTime, anyLmmActionWasSend)
		DEBUG.Printf("[user%d] : receive [%d] profiles for new faces feed", index, len(nfProfiles))
		if len(messages) == concurrentUsers/2 {
			INFO.Printf("[user%d] : user has a chat with everyone, terminate", index)
			finish <- 1
			return
		}
	}
}
