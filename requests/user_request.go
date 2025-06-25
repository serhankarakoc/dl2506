package requests

import (
	"davet.link/pkg/flashmessages"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserRequest struct {
	Name              string `form:"name" validate:"required,min=2"`
	Email             string `form:"email" validate:"required,email"`
	Password          string `form:"password"`
	Status            string `form:"status" validate:"required"`
	Type              string `form:"type" validate:"required,oneof=dashboard panel"`
	ResetToken        string `form:"reset_token"`
	EmailVerified     string `form:"email_verified"`
	VerificationToken string `form:"verification_token"`
	Provider          string `form:"provider"`
	ProviderID        string `form:"provider_id"`
}

func validateUserRequest(c *fiber.Ctx, req interface{}, errorMessages map[string]string, redirectPath string) error {
	if err := c.BodyParser(req); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı")
		return c.Redirect(redirectPath, fiber.StatusSeeOther)
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		err := err.(validator.ValidationErrors)[0]
		if msg, ok := errorMessages[err.Field()+"_"+err.Tag()]; ok {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, msg)
		} else {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz kullanıcı bilgileri")
		}
		return c.Redirect(redirectPath, fiber.StatusSeeOther)
	}
	return nil
}

func ValidateUserRequest(c *fiber.Ctx) error {
	var req UserRequest
	errorMessages := map[string]string{
		"Name_required": "Kullanıcı adı zorunludur",
		"Name_min":     "Kullanıcı adı en az 2 karakter olmalıdır",
		"Email_required": "E-posta adresi zorunludur",
		"Email_email":    "Geçerli bir e-posta adresi giriniz",
		"Password_required_without": "Şifre zorunludur",
		"Password_min": "Şifre en az 6 karakter olmalıdır",
		"Status_required": "Durum seçilmelidir",
		"Type_required":   "Kullanıcı tipi seçilmelidir",
		"Type_oneof":     "Kullanıcı tipi geçersiz",
	}
	if err := validateUserRequest(c, &req, errorMessages, c.OriginalURL()); err != nil {
		return err
	}
	c.Locals("userRequest", req)
	return nil
}
