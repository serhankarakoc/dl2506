package middlewares

import (
	"context"

	"davet.link/configs/sessionconfig"
	"davet.link/pkg/flashmessages"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	userID, err := sessionconfig.GetUserIDFromSession(c)
	if err != nil || userID == 0 {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Oturum bilgileri geçersiz")
		return c.Redirect("/auth/login")
	}

	authService := services.NewAuthService()
	user, err := authService.GetUserProfile(userID)
	if err != nil {
		sessionconfig.DestroySession(c)
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı bulunamadı")
		return c.Redirect("/auth/login")
	}

	ctx := context.WithValue(c.Context(), "user_id", userID)
	ctx = context.WithValue(ctx, "user_type", user.Type)
	ctx = context.WithValue(ctx, "user_email", user.Email)
	c.SetUserContext(ctx)

	c.Locals("userID", userID)
	c.Locals("userType", user.Type)
	c.Locals("userEmail", user.Email)

	return c.Next()
}
