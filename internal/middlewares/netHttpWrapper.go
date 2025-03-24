package middlewares

import (
	"net/http"
)

func NetHttpWrapper(next http.Handler) http.Handler {
	return NetHttpCors(NetHttpLogger(NetHttpRateLimit(NetHttpAuth(next))))
}
