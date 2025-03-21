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
	requests map[string]int
	limit    int
	window   time.Duration
	mutex    sync.Mutex
}

var rateLimiter = &RateLimiter{
	requests: make(map[string]int),
	window:   60 * time.Second, // 1 min
	limit:    20,
	mutex:    sync.Mutex{},
}

func (rl *RateLimiter) resetCount(clientIp string) {
	time.Sleep(rl.window)
	rl.mutex.Lock()
	delete(rl.requests, clientIp)
	rl.mutex.Unlock()
}

func (rl *RateLimiter) AllowRequest(clientIp string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Check the current count for the client
	count, found := rl.requests[clientIp]
	if !found || count < rl.limit {
		// Reset the count for a new window or increment for the current window
		if !found {
			go rl.resetCount(clientIp)
		}
		rl.requests[clientIp]++
		return true
	}

	return false
}

func RateLimit(c *fiber.Ctx) error {
	env := os.Getenv("GO_ENV")
	var clientIP string = ""

	if env == "production" {
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
