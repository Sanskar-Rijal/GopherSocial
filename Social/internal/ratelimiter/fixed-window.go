package ratelimiter

import (
	"sync"
	"time"
)


type FixedWindowRateLimiter struct {
	sync.RWMutex // lock — prevents race condition
	clients map[string] int // tracks requests per IP
	limit int  // max people allowed — 20
	window time.Duration // How long before the request count resets
}


// clients map looks like this in memory:
// clients = {
//     "192.168.1.1": 5,   // this IP made 5 requests so far
//     "192.168.1.2": 12,  // this IP made 12 requests so far
// }

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowRateLimiter{
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}


func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration){

// imagine TWO requests arrive at the EXACT same millisecond
// both try to do clients["192.168.1.1"]++ at the same time
// without lock — race condition — count gets corrupted
// same race condition concept you learned earlier!
	rl.Lock()
	//grab how many times the user have made request and if it don't exists it will return false
	count, exists := rl.clients[ip]
	rl.Unlock()

	//"Allow request if user is new OR still under limit"
	if !exists || count < rl.limit{
		rl.Lock()
		//if user is new reset count
		if !exists{
			go rl.resetCount(ip)
		}
		//user exists so increase the api request 
		rl.clients[ip]++;
		rl.Unlock()
		return true,0
	}
	
	//user has crossed limit send false and time left 

	return false,  rl.window
}

func (rl *FixedWindowRateLimiter) resetCount(ip string){
	time.Sleep(rl.window)
	rl.Lock()
	delete(rl.clients,ip)
	rl.Unlock()
}