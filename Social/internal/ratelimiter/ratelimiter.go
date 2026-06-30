package ratelimiter

import "time"


type Limiter interface{
	Allow(ip string) (bool, time.Duration) //time.Duration is for letting user know when he can make another request 
	//i.e how much time do he has to
}

type Config struct {
	RequestPerWindow int 
	TimeWindow time.Duration
	Enabled bool
}