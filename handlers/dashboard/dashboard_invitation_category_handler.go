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

type DashboardInvitationCategoryHandler struct {
	categoryService services.IInvitationCategoryService
}

func NewDashboardInvitationCategoryHandler() *DashboardInvitationCategoryHandler {
	svc := services.NewInvitationCategoryService()
	return &DashboardInvitationCategoryHandler{categoryService: svc}
}

func (h *DashboardInvitationCategoryHandler) ListCategories(c *fiber.Ctx) error {
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

	paginatedResult, dbErr := h.categoryService.GetAllCategories(params)

	renderData := fiber.Map{
		"Title":  "Davet Kategorileri",
		"Result": paginatedResult,
		"Params": params,
	}
	if dbErr != nil {
		renderData[renderer.FlashErrorKeyView] = "Kategoriler getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.InvitationCategory{},
			Meta: queryparams.PaginationMeta{
				CurrentPage: params.Page, PerPage: params.PerPage,
			},
		}
	}
	return renderer.Render(c, "dashboard/invitation-categories/list", "layouts/dashboard", renderData, http.StatusOK)
}

func (h *DashboardInvitationCategoryHandler) ShowCreateCategory(c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/invitation-categories/create", "layouts/dashboard", fiber.Map{
		"Title": "Yeni Kategori Ekle",
	})
}

func (h *DashboardInvitationCategoryHandler) CreateCategory(c *fiber.Ctx) error {
	if err := requests.ValidateInvitationCategoryRequest(c); err != nil {
		return renderCategoryFormError("dashboard/invitation-categories/create", "Yeni Kategori Ekle", c.Locals("invitationCategoryRequest"), err.Error(), c)
	}
	req := c.Locals("invitationCategoryRequest").(requests.InvitationCategoryRequest)
	category := &models.InvitationCategory{
		Name:     req.Name,
		Icon:     req.Icon,
		Template: req.Template,
		IsActive: req.IsActive == "true",
	}
	if err := h.categoryService.CreateCategory(c.UserContext(), category); err != nil {
		return renderCategoryFormError("dashboard/invitation-categories/create", "Yeni Kategori Ekle", req, "Kategori oluşturulamadı: "+err.Error(), c)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kategori başarıyla oluşturuldu.")
	return c.Redirect("/dashboard/invitation-categories", fiber.StatusFound)
}

func (h *DashboardInvitationCategoryHandler) ShowUpdateCategory(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	invitationCategory, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kategori bulunamadı.")
		return c.Redirect("/dashboard/invitation-categories", fiber.StatusSeeOther)
	}
	return renderer.Render(c, "dashboard/invitation-categories/update", "layouts/dashboard", fiber.Map{
		"Title":          "Kategori Düzenle",
		"InvitationCategory": invitationCategory,
	})
}

func (h *DashboardInvitationCategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := requests.ValidateInvitationCategoryRequest(c); err != nil {
		return err
	}
	req := c.Locals("invitationCategoryRequest").(requests.InvitationCategoryRequest)
	category := &models.InvitationCategory{
		Name:     req.Name,
		Icon:     req.Icon,
		Template: req.Template,
		IsActive: req.IsActive == "true",
	}
	userID, _ := c.Locals("userID").(uint)
	if err := h.categoryService.UpdateCategory(c.UserContext(), uint(id), category, userID); err != nil {
		return renderCategoryFormError("dashboard/invitation-categories/update", "Kategori Düzenle", req, "Kategori güncellenemedi: "+err.Error(), c)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kategori başarıyla güncellendi.")
	return c.Redirect("/dashboard/invitation-categories", fiber.StatusFound)
}

func (h *DashboardInvitationCategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	categoryID := uint(id)

	if err := h.categoryService.DeleteCategory(c.UserContext(), categoryID); err != nil {
		errMsg := "Kategori silinemedi: " + err.Error()
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/dashboard/invitation-categories", fiber.StatusSeeOther)
	}

	if strings.Contains(c.Get("Accept"), "application/json") {
		return c.JSON(fiber.Map{"message": "Kategori başarıyla silindi."})
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kategori başarıyla silindi.")
	return c.Redirect("/dashboard/invitation-categories", fiber.StatusFound)
}

func renderCategoryFormError(template string, title string, req any, message string, c *fiber.Ctx) error {
	return renderer.Render(c, template, "layouts/dashboard", fiber.Map{
		"Title":                    title,
		renderer.FlashErrorKeyView: message,
		renderer.FormDataKey:       req,
	}, http.StatusBadRequest)
}
