package middlewares

import (
	"davet.link/configs/sessionconfig"
	"davet.link/models"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
)

func GuestMiddleware(c *fiber.Ctx) error {
	userID, err := sessionconfig.GetUserIDFromSession(c)
	if err != nil || userID == 0 {
		return c.Next()
	}

	authService := services.NewAuthService()
	user, err := authService.GetUserProfile(userID)
	if err != nil {
		sessionconfig.DestroySession(c)
		return c.Next()
	}

	var redirectURL string
	switch user.Type {
	case models.Panel:
		redirectURL = "/panel/home"
	case models.Dashboard:
		redirectURL = "/dashboard/home"
	default:
		sessionconfig.DestroySession(c)
		return c.Next()
	}

	return c.Redirect(redirectURL)
}
