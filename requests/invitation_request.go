package requests

import (
	"davet.link/pkg/flashmessages"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type InvitationRequest struct {
	Image      string `form:"image"`
	CategoryID uint   `form:"category_id" validate:"required,gt=0"`
	Title      string `form:"title" validate:"required,min=3"`
	Type       string `form:"type" validate:"required"`
	Template   string `form:"template" validate:"required"`
	Date       string `form:"date" validate:"required,datetime=2006-01-02"`
	Time       string `form:"time" validate:"required,datetime=15:04"`
	Venue      string `form:"venue"`
	Address    string `form:"address"`
	Location   string `form:"location" validate:"omitempty,url"`
	Telephone  string `form:"telephone"`
	Detail     InvitationDetailRequest `form:"detail" validate:"required,dive"`
}

type InvitationDetailRequest struct {
	Title              string `form:"title"`
	BrideName          string `form:"bride_name"`
	BrideSurname       string `form:"bride_surname"`
	BrideMotherName    string `form:"bride_mother_name"`
	BrideMotherSurname string `form:"bride_mother_surname"`
	BrideFatherName    string `form:"bride_father_name"`
	BrideFatherSurname string `form:"bride_father_surname"`
	GroomName          string `form:"groom_name"`
	GroomSurname       string `form:"groom_surname"`
	GroomMotherName    string `form:"groom_mother_name"`
	GroomMotherSurname string `form:"groom_mother_surname"`
	GroomFatherName    string `form:"groom_father_name"`
	GroomFatherSurname string `form:"groom_father_surname"`
	Person             string `form:"person"`
	MotherName         string `form:"mother_name"`
	MotherSurname      string `form:"mother_surname"`
	FatherName         string `form:"father_name"`
	FatherSurname      string `form:"father_surname"`
	IsMotherLive       string `form:"is_mother_live"`
	IsFatherLive       string `form:"is_father_live"`
	IsBrideMotherLive  string `form:"is_bride_mother_live"`
	IsBrideFatherLive  string `form:"is_bride_father_live"`
	IsGroomMotherLive  string `form:"is_groom_mother_live"`
	IsGroomFatherLive  string `form:"is_groom_father_live"`
}

func ValidateInvitationRequest(c *fiber.Ctx) error {
	var req InvitationRequest

	if err := c.BodyParser(&req); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı.")
		return fmt.Errorf("body parser error: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.Locals("invitationRequest", req)
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Lütfen formdaki tüm zorunlu alanları doğru bir şekilde doldurun.")
		return fmt.Errorf("validation error: %w", err)
	}
	
	c.Locals("invitationRequest", req)
	return nil
}