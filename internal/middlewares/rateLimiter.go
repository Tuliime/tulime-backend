package middlewares

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RateLimiter struct {
	requests     map[string]int
	blockedUntil map[string]time.Time
	limit        int
	window       time.Duration
	mutex        sync.Mutex
}

var rateLimiter = &RateLimiter{
	requests:     make(map[string]int),
	blockedUntil: make(map[string]time.Time),
	window:       60 * time.Second, // 1 min
	limit:        40,
	mutex:        sync.Mutex{},
}

func (rl *RateLimiter) resetCount(clientIp string) {
	time.Sleep(rl.window)
	rl.mutex.Lock()
	delete(rl.requests, clientIp)
	delete(rl.blockedUntil, clientIp)
	rl.mutex.Unlock()
}

// AllowRequest blocks requests until the window expires
func (rl *RateLimiter) AllowRequest(clientIp string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// If the IP is already blocked, check if the block has expired
	if unblockTime, blocked := rl.blockedUntil[clientIp]; blocked {
		if time.Now().Before(unblockTime) {
			return false
		}

		delete(rl.blockedUntil, clientIp)
		delete(rl.requests, clientIp)
	}

	count := rl.requests[clientIp]
	if count >= rl.limit {
		// Block the IP and set the unblock time
		rl.blockedUntil[clientIp] = time.Now().Add(rl.window)
		return false
	}

	// Increment request count and start the reset timer if first request
	if count == 0 {
		go rl.resetCount(clientIp)
	}
	rl.requests[clientIp]++

	return true
}

func RateLimit(c *fiber.Ctx) error {
	var clientIP string = ""

	if os.Getenv("GO_ENV") == "production" {
		clientIP = c.Get("CF-Connecting-IP") // Cloudflare's real IP header
		log.Println("IP Cloudflare header: ", clientIP)
	}

	if clientIP == "" {
		forwardedFor := c.Get("X-Forwarded-For")
		if forwardedFor != "" {
			clientIP = strings.Split(forwardedFor, ",")[0]
			log.Println("IP X-Forwarded-For: ", clientIP)
		}
	}

	if clientIP == "" {
		clientIP = c.Get("X-Real-IP")
		log.Println("IP X-Real-IP: ", clientIP)
	}

	if clientIP == "" {
		clientIP = c.IP()
		log.Println("IP default: ", clientIP)
	}

	log.Println("clientIP address:", clientIP)

	if !rateLimiter.AllowRequest(clientIP) {
		log.Println("Request denied for:", clientIP)
		return fiber.NewError(fiber.StatusTooManyRequests,
			"You have made too many requests! Please try again later.")
	}

	c.Locals("clientIP", clientIP)

	return c.Next()
}
