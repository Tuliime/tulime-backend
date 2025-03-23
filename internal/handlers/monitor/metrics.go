package monitor

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var GetMetrics = func(c *fiber.Ctx) error {
	adaptor.HTTPHandler(promhttp.Handler())(c)
	return nil
}
