package requests

import (
	"davet.link/pkg/flashmessages"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"regexp"
	"sort"
	"strconv"
)

type CardRequest struct {
	Name            string                   `form:"name" validate:"required,min=2"`
	Slug            string                   `form:"slug" validate:"required"`
	Title           string                   `form:"title"`
	Photo           string                   `form:"photo"`
	Telephone       string                   `form:"telephone"`
	Email           string                   `form:"email" validate:"omitempty,email"`
	Location        string                   `form:"location" validate:"omitempty,url"`
	WebsiteUrl      string                   `form:"website_url" validate:"omitempty,url"`
	StoreUrl        string                   `form:"store_url" validate:"omitempty,url"`
	IsActive        string                   `form:"is_active" validate:"required,boolean"`
	CardBanks       []CardBankRequest        `validate:"dive"`
	CardSocialMedia []CardSocialMediaRequest `validate:"dive"`
}

type CardBankRequest struct {
	ID     uint   `validate:"-"`
	BankID uint   `validate:"required"`
	IBAN   string `validate:"required"`
}

type CardSocialMediaRequest struct {
	ID            uint   `validate:"-"`
	SocialMediaID uint   `validate:"required"`
	URL           string `validate:"required,url"`
}

func parseBanksFromMap(formValues map[string][]string) []CardBankRequest {
	banksData := make(map[int]*CardBankRequest)
	re := regexp.MustCompile(`^card_banks\[(\d+)\]\[(bank_id|iban|id)\]$`)

	for key, values := range formValues {
		if len(values) == 0 {
			continue
		}
		value := values[0]
		matches := re.FindStringSubmatch(key)
		if len(matches) == 3 {
			index, _ := strconv.Atoi(matches[1])
			if _, ok := banksData[index]; !ok {
				banksData[index] = &CardBankRequest{}
			}
			switch matches[2] {
			case "bank_id":
				id, _ := strconv.ParseUint(value, 10, 64)
				banksData[index].BankID = uint(id)
			case "iban":
				banksData[index].IBAN = value
			case "id":
				id, _ := strconv.ParseUint(value, 10, 64)
				banksData[index].ID = uint(id)
			}
		}
	}

	var indexes []int
	for k := range banksData {
		indexes = append(indexes, k)
	}
	sort.Ints(indexes)
	var cardBanks []CardBankRequest
	for _, i := range indexes {
		if banksData[i] != nil && banksData[i].BankID > 0 && banksData[i].IBAN != "" {
			cardBanks = append(cardBanks, *banksData[i])
		}
	}
	return cardBanks
}

func parseSocialMediaFromMap(formValues map[string][]string) []CardSocialMediaRequest {
	socialData := make(map[int]*CardSocialMediaRequest)
	re := regexp.MustCompile(`^card_social_media\[(\d+)\]\[(social_media_id|url|id)\]$`)

	for key, values := range formValues {
		if len(values) == 0 {
			continue
		}
		value := values[0]
		matches := re.FindStringSubmatch(key)
		if len(matches) == 3 {
			index, _ := strconv.Atoi(matches[1])
			if _, ok := socialData[index]; !ok {
				socialData[index] = &CardSocialMediaRequest{}
			}
			switch matches[2] {
			case "social_media_id":
				id, _ := strconv.ParseUint(value, 10, 64)
				socialData[index].SocialMediaID = uint(id)
			case "url":
				socialData[index].URL = value
			case "id":
				id, _ := strconv.ParseUint(value, 10, 64)
				socialData[index].ID = uint(id)
			}
		}
	}

	var indexes []int
	for k := range socialData {
		indexes = append(indexes, k)
	}
	sort.Ints(indexes)
	var cardSocialMedia []CardSocialMediaRequest
	for _, i := range indexes {
		if socialData[i] != nil && socialData[i].SocialMediaID > 0 && socialData[i].URL != "" {
			cardSocialMedia = append(cardSocialMedia, *socialData[i])
		}
	}
	return cardSocialMedia
}

func mainValidateCardRequest(c *fiber.Ctx) (*CardRequest, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("multipart form parse error: %w", err)
	}
	formValues := form.Value

	req := CardRequest{}
	if val, ok := formValues["name"]; ok && len(val) > 0 { req.Name = val[0] }
	if val, ok := formValues["slug"]; ok && len(val) > 0 { req.Slug = val[0] }
	if val, ok := formValues["title"]; ok && len(val) > 0 { req.Title = val[0] }
	if val, ok := formValues["telephone"]; ok && len(val) > 0 { req.Telephone = val[0] }
	if val, ok := formValues["email"]; ok && len(val) > 0 { req.Email = val[0] }
	if val, ok := formValues["location"]; ok && len(val) > 0 { req.Location = val[0] }
	if val, ok := formValues["website_url"]; ok && len(val) > 0 { req.WebsiteUrl = val[0] }
	if val, ok := formValues["store_url"]; ok && len(val) > 0 { req.StoreUrl = val[0] }
	if val, ok := formValues["is_active"]; ok && len(val) > 0 { req.IsActive = val[0] }

	req.CardBanks = parseBanksFromMap(formValues)
	req.CardSocialMedia = parseSocialMediaFromMap(formValues)

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return &req, fmt.Errorf("validation error: %w", err)
	}

	return &req, nil
}

func ValidateCardRequest(c *fiber.Ctx) error {
	req, err := mainValidateCardRequest(c)
	if err != nil {
		if req != nil {
			c.Locals("cardRequest", *req)
		}

		if _, ok := err.(validator.ValidationErrors); ok {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Lütfen formdaki tüm zorunlu alanları doğru bir şekilde doldurun.")
		} else {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "İsteğiniz işlenirken bir sorun oluştu.")
		}
		
		return err
	}

	c.Locals("cardRequest", *req)
	return nil
}

func ValidateCardRequestWithPath(c *fiber.Ctx, redirectPath string) error {
	return ValidateCardRequest(c)
}