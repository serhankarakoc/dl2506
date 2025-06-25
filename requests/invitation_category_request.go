package requests

import (
	"davet.link/pkg/flashmessages"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type InvitationCategoryRequest struct {
	Name     string `form:"name" validate:"required,min=2"`
	Icon     string `form:"icon" validate:"required"`
	Template string `form:"template" validate:"required"`
	IsActive string `form:"is_active" validate:"required"`
}

func validateInvitationCategoryRequest(c *fiber.Ctx, req interface{}, errorMessages map[string]string, redirectPath string) error {
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
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz kategori bilgileri")
		}
		return c.Redirect(redirectPath, fiber.StatusSeeOther)
	}
	return nil
}

func ValidateInvitationCategoryRequest(c *fiber.Ctx) error {
	var req InvitationCategoryRequest
	errorMessages := map[string]string{
		"Name_required":     "Kategori adı zorunludur",
		"Name_min":          "Kategori adı en az 2 karakter olmalıdır",
		"Icon_required":     "İkon zorunludur",
		"Template_required": "Şablon zorunludur",
		"IsActive_required": "Durum (Aktif/Pasif) seçilmelidir",
	}
	if err := validateInvitationCategoryRequest(c, &req, errorMessages, "/dashboard/invitation-categories/create"); err != nil {
		return err
	}
	c.Locals("invitationCategoryRequest", req)
	return nil
}
