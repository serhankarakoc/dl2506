package requests

import (
	"davet.link/pkg/flashmessages"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type InvitationParticipantRequest struct {
	Title       string `form:"title" validate:"required,min=2"`
	PhoneNumber string `form:"phone_number" validate:"required,min=10"`
	GuestCount  int    `form:"guest_count" validate:"required,min=1"`
}

func validateInvitationParticipantRequest(c *fiber.Ctx, req interface{}, errorMessages map[string]string, redirectPath string) error {
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
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz katılımcı bilgileri")
		}
		return c.Redirect(redirectPath, fiber.StatusSeeOther)
	}
	return nil
}

func ValidateInvitationParticipantRequest(c *fiber.Ctx) error {
	var req InvitationParticipantRequest
	errorMessages := map[string]string{
		"Title_required":       "Ad Soyad zorunludur",
		"Title_min":            "Ad Soyad en az 2 karakter olmalıdır",
		"PhoneNumber_required": "Telefon numarası zorunludur",
		"PhoneNumber_min":      "Telefon numarası en az 10 karakter olmalıdır",
		"GuestCount_required":  "Kişi sayısı zorunludur",
		"GuestCount_min":       "Kişi sayısı en az 1 olmalıdır",
	}
	if err := validateInvitationParticipantRequest(c, &req, errorMessages, ""); err != nil {
		return err
	}
	c.Locals("invitationParticipantRequest", req)
	return nil
}
