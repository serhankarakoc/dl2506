package limiterconfig

import (
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func GetLimiterConfig() limiter.Config {
	return limiter.Config{
		Max:        1000,
		Expiration: 60,
	}
}
