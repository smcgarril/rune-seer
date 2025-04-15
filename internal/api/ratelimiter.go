package api

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors        = make(map[string]*visitor)
	mu              = sync.Mutex{}
	rateLimit       = rate.Every(1 * time.Second)
	burst           = 5
	cleanupInterval = 5 * time.Minute
	visitorTTL      = 10 * time.Minute
)

func init() {
	go cleanupVisitors()
}

func cleanupVisitors() {
	time.Sleep(cleanupInterval)
	mu.Lock()
	defer mu.Unlock()

	for ip, v := range visitors {
		if time.Since(v.lastSeen) > visitorTTL {
			delete(visitors, ip)
		}
	}
}

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		limiter := getVisitor(ip)

		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		lim := rate.NewLimiter(rateLimit, burst)
		visitors[ip] = &visitor{limiter: lim, lastSeen: time.Now()}
		return lim
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func getIP(r *http.Request) string {
	// check Fly.io header
	if flyIP := r.Header.Get("X-Fly-Client-IP"); flyIP != "" {
		return flyIP
	}

	// Check generic proxy header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check nginx or custom proxy
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}

	// Fallback to remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
