package requests

import (
	"davet.link/pkg/flashmessages"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SocialMediaRequest struct {
	Name     string `form:"name" validate:"required,min=2"`
	Icon     string `form:"icon" validate:"required"`
	IsActive string `form:"is_active" validate:"required"`
}

func validateSocialMediaRequest(c *fiber.Ctx, req interface{}, errorMessages map[string]string, redirectPath string) error {
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
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz sosyal medya bilgileri")
		}
		return c.Redirect(redirectPath, fiber.StatusSeeOther)
	}
	return nil
}

func ValidateSocialMediaRequest(c *fiber.Ctx) error {
	var req SocialMediaRequest
	errorMessages := map[string]string{
		"Name_required":     "Sosyal medya adı zorunludur",
		"Name_min":          "Sosyal medya adı en az 2 karakter olmalıdır",
		"Icon_required":     "İkon zorunludur",
		"IsActive_required": "Durum (Aktif/Pasif) seçilmelidir",
	}
	if err := validateSocialMediaRequest(c, &req, errorMessages, "/dashboard/social-media/create"); err != nil {
		return err
	}
	c.Locals("socialMediaRequest", req)
	return nil
}
