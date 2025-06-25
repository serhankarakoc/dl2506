package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"fmt"

	"davet.link/configs/logconfig"
	"davet.link/models"
	"davet.link/pkg/filemanager"
	"davet.link/pkg/flashmessages"
	"davet.link/pkg/queryparams"
	"davet.link/pkg/renderer"
	"davet.link/requests"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DashboardCardHandler struct {
	cardService services.ICardService
}

func NewDashboardCardHandler() *DashboardCardHandler {
	return &DashboardCardHandler{cardService: services.NewCardService()}
}

func (h *DashboardCardHandler) ListCards(c *fiber.Ctx) error {
	var params queryparams.ListParams
	if err := c.QueryParser(&params); err != nil {
		logconfig.Log.Warn("Kart listesi: Query parametreleri parse edilemedi, varsayılanlar kullanılıyor.", zap.Error(err))
		params = queryparams.DefaultListParams()
	}

	if params.Page <= 0 {
		params.Page = queryparams.DefaultPage
	}
	if params.PerPage <= 0 || params.PerPage > queryparams.MaxPerPage {
		params.PerPage = queryparams.DefaultPerPage
	}
	if params.SortBy == "" {
		params.SortBy = queryparams.DefaultSortBy
	}
	if params.OrderBy == "" {
		params.OrderBy = queryparams.DefaultOrderBy
	}

	result, err := h.cardService.GetAllCards(params)

	renderData := fiber.Map{
		"Title":  "Kartlar",
		"Result": result,
		"Params": params,
	}
	if err != nil {
		logconfig.Log.Error("Kart listesi DB Hatası", zap.Error(err))
		renderData[renderer.FlashErrorKeyView] = "Kartlar getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.Card{},
			Meta: queryparams.PaginationMeta{CurrentPage: params.Page, PerPage: params.PerPage},
		}
	}
	return renderer.Render(c, "dashboard/cards/list", "layouts/dashboard", renderData, http.StatusOK)
}

func (h *DashboardCardHandler) ShowCreateCard(c *fiber.Ctx) error {
	bankService := services.NewBankService()
	socialMediaService := services.NewSocialMediaService()
	banksResult, _ := bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
	socialMediasResult, _ := socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})
	return renderer.Render(c, "dashboard/cards/create", "layouts/dashboard", fiber.Map{
		"Title": "Yeni Kart Oluştur",
		"Banks": banksResult.Data,
		"SocialMedias": socialMediasResult.Data,
	})
}

func (h *DashboardCardHandler) CreateCard(c *fiber.Ctx) error {
	userIDValue := c.Locals("userID")
	userID, ok := userIDValue.(uint)
	if !ok {
		return c.Status(http.StatusUnauthorized).SendString("Geçerli kullanıcı bulunamadı.")
	}

	if err := requests.ValidateCardRequest(c); err != nil {
		req, _ := c.Locals("cardRequest").(requests.CardRequest)

		bankService := services.NewBankService()
		socialMediaService := services.NewSocialMediaService()
		banksResult, _ := bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
		socialMediasResult, _ := socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})

		return renderer.Render(c, "dashboard/cards/create", "layouts/dashboard", fiber.Map{
			"Title":        "Yeni Kart Oluştur",
			"Banks":        banksResult.Data,
			"SocialMedias": socialMediasResult.Data,
			"FormData":     req,
		})
	}

	req := c.Locals("cardRequest").(requests.CardRequest)

	card := &models.Card{
		Name:       req.Name,
		Slug:       req.Slug,
		Title:      req.Title,
		Telephone:  req.Telephone,
		Email:      req.Email,
		Location:   req.Location,
		WebsiteUrl: req.WebsiteUrl,
		StoreUrl:   req.StoreUrl,
		IsActive:   req.IsActive == "true",
		UserID:     userID,
	}

	newFileName, err := filemanager.UploadFile(c, "photo", "cards")
	if err != nil && err != filemanager.ErrFileNotProvided {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Fotoğraf yüklenemedi: "+err.Error())

		bankService := services.NewBankService()
		socialMediaService := services.NewSocialMediaService()
		banksResult, _ := bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
		socialMediasResult, _ := socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})

		return renderer.Render(c, "dashboard/cards/create", "layouts/dashboard", fiber.Map{
			"Title": "Yeni Kart Oluştur", "Banks": banksResult.Data, "SocialMedias": socialMediasResult.Data, "FormData": req,
		})
	}
	card.Photo = newFileName

	for _, cb := range req.CardBanks {
		card.CardBanks = append(card.CardBanks, models.CardBank{BankID: cb.BankID, IBAN: cb.IBAN})
	}
	for _, cs := range req.CardSocialMedia {
		card.CardSocialMedia = append(card.CardSocialMedia, models.CardSocialMedia{SocialMediaID: cs.SocialMediaID, URL: cs.URL})
	}

	if err := h.cardService.CreateCardWithRelations(c.UserContext(), card); err != nil {
		filemanager.DeleteFile("cards", newFileName)
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kart oluşturulamadı: "+err.Error())

		bankService := services.NewBankService()
		socialMediaService := services.NewSocialMediaService()
		banksResult, _ := bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
		socialMediasResult, _ := socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})
		
		return renderer.Render(c, "dashboard/cards/create", "layouts/dashboard", fiber.Map{
			"Title": "Yeni Kart Oluştur", "Banks": banksResult.Data, "SocialMedias": socialMediasResult.Data, "FormData": req,
		})
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kart başarıyla oluşturuldu.")
	return c.Redirect("/dashboard/cards", http.StatusFound)
}

func (h *DashboardCardHandler) ShowUpdateCard(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz kart ID'si.")
		return c.Redirect("/dashboard/cards", http.StatusSeeOther)
	}

	card, err := h.cardService.GetCardByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kart bulunamadı.")
		return c.Redirect("/dashboard/cards", http.StatusSeeOther)
	}

	bankService := services.NewBankService()
	socialMediaService := services.NewSocialMediaService()
	banksResult, _ := bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
	socialMediasResult, _ := socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})

	return renderer.Render(c, "dashboard/cards/update", "layouts/dashboard", fiber.Map{
		"Title":    "Kart Düzenle",
		"FormData": card,
		"Banks":    banksResult.Data,
		"SocialMedias": socialMediasResult.Data,
	})
}

func (h *DashboardCardHandler) UpdateCard(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz kart ID'si.")
		return c.Redirect("/dashboard/cards", http.StatusSeeOther)
	}

	redirectURL := fmt.Sprintf("/dashboard/cards/update/%d", id)

	existingCard, err := h.cardService.GetCardByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Güncellenecek kart bulunamadı.")
		return c.Redirect("/dashboard/cards", http.StatusSeeOther)
	}

	if err := requests.ValidateCardRequestWithPath(c, redirectURL); err != nil {
		req, _ := c.Locals("cardRequest").(requests.CardRequest)
		bankService := services.NewBankService()
		socialMediaService := services.NewSocialMediaService()
		banksResult, _ := bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
		socialMediasResult, _ := socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})

		return renderer.Render(c, "dashboard/cards/update", "layouts/dashboard", fiber.Map{
			"Title":        "Kart Düzenle",
			"FormData":     req,
			"Banks":        banksResult.Data,
			"SocialMedias": socialMediasResult.Data,
		})
	}

	req := c.Locals("cardRequest").(requests.CardRequest)
	userID, _ := c.Locals("userID").(uint)

	if req.Slug != existingCard.Slug {
		isAvailable, err := h.cardService.IsSlugAvailable(req.Slug, uint(id))
		if err != nil {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Slug kontrolü sırasında bir hata oluştu.")
			return c.Redirect(redirectURL, http.StatusSeeOther)
		}
		if !isAvailable {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Bu kullanıcı adı zaten alınmış.")
			c.Locals("cardRequest", req)
			bankService := services.NewBankService()
			socialMediaService := services.NewSocialMediaService()
			banksResult, _ := bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
			socialMediasResult, _ := socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})
			return renderer.Render(c, "dashboard/cards/update", "layouts/dashboard", fiber.Map{
				"Title": "Kart Düzenle", "FormData": req, "Banks": banksResult.Data, "SocialMedias": socialMediasResult.Data,
			})
		}
	}

	newFileName, err := filemanager.UploadFile(c, "photo", "cards")
	if err != nil && err != filemanager.ErrFileNotProvided {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Yeni fotoğraf yüklenemedi: "+err.Error())
		return c.Redirect(redirectURL, http.StatusSeeOther)
	}

	var oldPhotoToDelete string
	if newFileName != "" {
		oldPhotoToDelete = existingCard.Photo
		existingCard.Photo = newFileName
	}

	existingCard.Name = req.Name
	existingCard.Slug = req.Slug
	existingCard.Title = req.Title
	existingCard.Telephone = req.Telephone
	existingCard.Email = req.Email
	existingCard.Location = req.Location
	existingCard.WebsiteUrl = req.WebsiteUrl
	existingCard.StoreUrl = req.StoreUrl
	existingCard.IsActive = req.IsActive == "true"
	existingCard.UpdatedBy = userID

	existingCard.CardBanks = []models.CardBank{}
	for _, cb := range req.CardBanks {
		existingCard.CardBanks = append(existingCard.CardBanks, models.CardBank{
			BaseModel: models.BaseModel{ID: cb.ID},
			BankID:    cb.BankID,
			IBAN:      cb.IBAN,
		})
	}

	existingCard.CardSocialMedia = []models.CardSocialMedia{}
	for _, cs := range req.CardSocialMedia {
		existingCard.CardSocialMedia = append(existingCard.CardSocialMedia, models.CardSocialMedia{
			BaseModel:     models.BaseModel{ID: cs.ID},
			SocialMediaID: cs.SocialMediaID,
			URL:           cs.URL,
		})
	}

	if err := h.cardService.UpdateCardWithRelations(c.UserContext(), existingCard); err != nil {
		if newFileName != "" {
			filemanager.DeleteFile("cards", newFileName)
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kart güncellenirken bir veritabanı hatası oluştu.")
		return c.Redirect(redirectURL, http.StatusSeeOther)
	}

	if oldPhotoToDelete != "" {
		filemanager.DeleteFile("cards", oldPhotoToDelete)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kart başarıyla güncellendi.")
	return c.Redirect("/dashboard/cards", http.StatusFound)
}

func (h *DashboardCardHandler) DeleteCard(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.cardService.DeleteCardWithRelations(c.UserContext(), uint(id)); err != nil {
		errMsg := "Kart silinemedi: " + err.Error()
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/dashboard/cards", fiber.StatusSeeOther)
	}
	if strings.Contains(c.Get("Accept"), "application/json") {
		return c.JSON(fiber.Map{"message": "Kart başarıyla silindi."})
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kart başarıyla silindi.")
	return c.Redirect("/dashboard/cards", http.StatusFound)
}

func (h *DashboardCardHandler) SlugCheck(c *fiber.Ctx) error {
	slug := c.Query("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_available": false,
			"message":      "Slug parametresi eksik.",
		})
	}

	excludeIDStr := c.Query("exclude_id")
	var excludeID uint = 0
	if excludeIDStr != "" {
		id, err := strconv.ParseUint(excludeIDStr, 10, 32)
		if err == nil {
			excludeID = uint(id)
		}
	}

	isAvailable, err := h.cardService.IsSlugAvailable(slug, excludeID)
	if err != nil {
		logconfig.Log.Error("Slug kontrolü sırasında veritabanı hatası", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"is_available": false,
			"message":      "Sunucu hatası.",
		})
	}

	return c.JSON(fiber.Map{
		"is_available": isAvailable,
	})
}