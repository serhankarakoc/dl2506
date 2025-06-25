package middlewares

import (
	"time"

	"davet.link/configs/logconfig"
	"github.com/gofiber/fiber/v2"
)

func ZapLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()

		logconfig.SLog.Infow("request",
			"method", c.Method(),
			"path", c.OriginalURL(),
			"status", c.Response().StatusCode(),
			"latency", stop.Sub(start).String(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
		)
		return err
	}
}
