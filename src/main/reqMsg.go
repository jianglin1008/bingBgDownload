package main

import "time"

//format=js&idx=1&n=1&nc=1421493718026&pid=hp
type ReqMsg struct {
	Format string
	Idx    int
	N      int
	Nc     int64
	Pid    string
}

var DefaultReqMsg = &ReqMsg{"js", 1, 1, time.Now().Unix(), "hp"}

func NewReqMsg(idx int) *ReqMsg {
	ret := DefaultReqMsg
	ret.Idx = idx
	return ret
}
