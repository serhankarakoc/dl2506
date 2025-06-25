package requests

import (
	"davet.link/pkg/flashmessages"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type BankRequest struct {
	Name     string `form:"name" validate:"required,min=2"`
	IsActive string `form:"is_active" validate:"required"`
}

func validateBankRequest(c *fiber.Ctx, req interface{}, errorMessages map[string]string, redirectPath string) error {
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
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz banka bilgileri")
		}
		return c.Redirect(redirectPath, fiber.StatusSeeOther)
	}
	return nil
}

func ValidateBankRequest(c *fiber.Ctx) error {
	var req BankRequest
	errorMessages := map[string]string{
		"Name_required": "Banka adı zorunludur",
		"Name_min":      "Banka adı en az 2 karakter olmalıdır",
		"IsActive_required": "Durum (Aktif/Pasif) seçilmelidir",
	}
	if err := validateBankRequest(c, &req, errorMessages, "/dashboard/banks/create"); err != nil {
		return err
	}
	c.Locals("bankRequest", req)
	return nil
}
