package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Tuliime/tulime-backend/internal/packages"
)

func NetHttpRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientIP string = ""

		if os.Getenv("GO_ENV") == "production" {
			clientIP := r.Header.Get("CF-Connecting-IP") // Cloudflare's real IP header
			log.Println("IP Cloudflare header: ", clientIP)
		}

		if clientIP == "" {
			forwardedFor := r.Header.Get("X-Forwarded-For")
			if forwardedFor != "" {
				clientIP = strings.Split(forwardedFor, ",")[0]
				log.Println("IP X-Forwarded-For: ", clientIP)
			}
		}

		if clientIP == "" {
			clientIP = r.Header.Get("X-Real-IP")
			log.Println("IP X-Real-IP: ", clientIP)
		}

		if clientIP == "" {
			clientIP = r.RemoteAddr
			log.Println("IP default: ", clientIP)
		}

		log.Println("clientIP address:", clientIP)

		if !rateLimiter.AllowRequest(clientIP) {
			packages.AppError("You have made too many requests!, Please try again later.", 429, w)
			log.Println("Request denied for:", clientIP)
			return
		}

		next.ServeHTTP(w, r)
	})
}
