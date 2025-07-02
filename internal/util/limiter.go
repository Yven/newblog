package util

import (
	"sync"

	"golang.org/x/time/rate"
)

type Visitors struct {
	Bucket map[string]*rate.Limiter
	mu     sync.Mutex
}

func NewVisitors() *Visitors {
	return &Visitors{
		Bucket: make(map[string]*rate.Limiter),
		mu:     sync.Mutex{},
	}
}

// 每个 IP 获取自己的限流器
func (v *Visitors) GetVisitor(ip string) *rate.Limiter {
	v.mu.Lock()
	defer v.mu.Unlock()

	limiter, exists := v.Bucket[ip]
	if !exists {
		limiter = rate.NewLimiter(5, 10)
		v.Bucket[ip] = limiter
	}

	return limiter
}
