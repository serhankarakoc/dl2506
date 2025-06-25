package handlers

import (
	"net/http"
	"strings"

	"davet.link/models"
	"davet.link/pkg/flashmessages"
	"davet.link/pkg/queryparams"
	"davet.link/pkg/renderer"
	"davet.link/requests"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
)

type DashboardSocialMediaHandler struct {
	socialMediaService services.ISocialMediaService
}

func NewDashboardSocialMediaHandler() *DashboardSocialMediaHandler {
	svc := services.NewSocialMediaService()
	return &DashboardSocialMediaHandler{socialMediaService: svc}
}

func (h *DashboardSocialMediaHandler) ListSocialMedias(c *fiber.Ctx) error {
	var params queryparams.ListParams
	if err := c.QueryParser(&params); err != nil {
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

	paginatedResult, dbErr := h.socialMediaService.GetAllSocialMedias(params)

	renderData := fiber.Map{
		"Title":  "Sosyal Medya",
		"Result": paginatedResult,
		"Params": params,
	}
	if dbErr != nil {
		renderData[renderer.FlashErrorKeyView] = "Sosyal medya kayıtları getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.SocialMedia{},
			Meta: queryparams.PaginationMeta{
				CurrentPage: params.Page, PerPage: params.PerPage,
			},
		}
	}
	return renderer.Render(c, "dashboard/social-media/list", "layouts/dashboard", renderData, http.StatusOK)
}

func (h *DashboardSocialMediaHandler) ShowCreateSocialMedia(c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/social-media/create", "layouts/dashboard", fiber.Map{
		"Title": "Yeni Sosyal Medya Ekle",
	})
}

func (h *DashboardSocialMediaHandler) CreateSocialMedia(c *fiber.Ctx) error {
	if err := requests.ValidateSocialMediaRequest(c); err != nil {
		return renderSocialMediaFormError("dashboard/social-media/create", "Yeni Sosyal Medya Ekle", c.Locals("socialMediaRequest"), err.Error(), c)
	}
	req := c.Locals("socialMediaRequest").(requests.SocialMediaRequest)
	socialMedia := &models.SocialMedia{
		Name:     req.Name,
		Icon:     req.Icon,
		IsActive: req.IsActive == "true",
	}
	if err := h.socialMediaService.CreateSocialMedia(c.UserContext(), socialMedia); err != nil {
		return renderSocialMediaFormError("dashboard/social-media/create", "Yeni Sosyal Medya Ekle", req, "Kayıt oluşturulamadı: "+err.Error(), c)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kayıt başarıyla oluşturuldu.")
	return c.Redirect("/dashboard/social-media", fiber.StatusFound)
}

func (h *DashboardSocialMediaHandler) ShowUpdateSocialMedia(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	socialMedia, err := h.socialMediaService.GetSocialMediaByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kayıt bulunamadı.")
		return c.Redirect("/dashboard/social-media", fiber.StatusSeeOther)
	}
	return renderer.Render(c, "dashboard/social-media/update", "layouts/dashboard", fiber.Map{
		"Title":       "Sosyal Medya Düzenle",
		"SocialMedia": socialMedia,
	})
}

func (h *DashboardSocialMediaHandler) UpdateSocialMedia(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := requests.ValidateSocialMediaRequest(c); err != nil {
		return err
	}
	req := c.Locals("socialMediaRequest").(requests.SocialMediaRequest)
	socialMedia := &models.SocialMedia{
		Name:     req.Name,
		Icon:     req.Icon,
		IsActive: req.IsActive == "true",
	}
	userID, _ := c.Locals("userID").(uint)
	if err := h.socialMediaService.UpdateSocialMedia(c.UserContext(), uint(id), socialMedia, userID); err != nil {
		return renderSocialMediaFormError("dashboard/social-media/update", "Sosyal Medya Güncelle", req, "Kayıt güncellenemedi: "+err.Error(), c)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kayıt başarıyla güncellendi.")
	return c.Redirect("/dashboard/social-media", fiber.StatusFound)
}

func (h *DashboardSocialMediaHandler) DeleteSocialMedia(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	socialMediaID := uint(id)

	if err := h.socialMediaService.DeleteSocialMedia(c.UserContext(), socialMediaID); err != nil {
		errMsg := "Kayıt silinemedi: " + err.Error()
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/dashboard/social-media", fiber.StatusSeeOther)
	}

	if strings.Contains(c.Get("Accept"), "application/json") {
		return c.JSON(fiber.Map{"message": "Kayıt başarıyla silindi."})
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kayıt başarıyla silindi.")
	return c.Redirect("/dashboard/social-media", fiber.StatusFound)
}

func renderSocialMediaFormError(template string, title string, req any, message string, c *fiber.Ctx) error {
	return renderer.Render(c, template, "layouts/dashboard", fiber.Map{
		"Title":                    title,
		renderer.FlashErrorKeyView: message,
		renderer.FormDataKey:       req,
	}, http.StatusBadRequest)
}
