package middlewares

import (
	"davet.link/configs/sessionconfig"
	"davet.link/models"
	"davet.link/pkg/flashmessages"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
)

func TypeMiddleware(requiredType models.UserType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := sessionconfig.GetUserIDFromSession(c)
		if err != nil || userID == 0 {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Yetkili oturum bulunamadı")
			return c.Redirect("/auth/login")
		}

		authService := services.NewAuthService()
		user, err := authService.GetUserProfile(userID)
		if err != nil {
			sessionconfig.DestroySession(c)
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı bulunamadı")
			return c.Redirect("/auth/login")
		}

		if user.Type != requiredType {
			sessionconfig.DestroySession(c)
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Bu sayfaya erişim izniniz yok")
			return c.Redirect("/auth/login")
		}

		return c.Next()
	}
}
