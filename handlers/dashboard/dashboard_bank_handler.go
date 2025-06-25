package handlers

import (
	"net/http"
	"strings"

	"davet.link/configs/logconfig"
	"davet.link/models"
	"davet.link/pkg/flashmessages"
	"davet.link/pkg/queryparams"
	"davet.link/pkg/renderer"
	"davet.link/requests"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DashboardBankHandler struct {
	bankService services.IBankService
}

func NewDashboardBankHandler() *DashboardBankHandler {
	svc := services.NewBankService()
	return &DashboardBankHandler{bankService: svc}
}

func (h *DashboardBankHandler) ListBanks(c *fiber.Ctx) error {
	var params queryparams.ListParams
	if err := c.QueryParser(&params); err != nil {
		logconfig.Log.Warn("Banka listesi: Query parametreleri parse edilemedi, varsayılanlar kullanılıyor.", zap.Error(err))
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

	paginatedResult, dbErr := h.bankService.GetAllBanks(params)

	renderData := fiber.Map{
		"Title":  "Bankalar",
		"Result": paginatedResult,
		"Params": params,
	}
	if dbErr != nil {
		logconfig.Log.Error("Banka listesi DB Hatası", zap.Error(dbErr))
		renderData[renderer.FlashErrorKeyView] = "Bankalar getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.Bank{},
			Meta: queryparams.PaginationMeta{
				CurrentPage: params.Page, PerPage: params.PerPage,
			},
		}
	}
	return renderer.Render(c, "dashboard/banks/list", "layouts/dashboard", renderData, http.StatusOK)
}

func (h *DashboardBankHandler) ShowCreateBank(c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/banks/create", "layouts/dashboard", fiber.Map{
		"Title": "Yeni Banka Ekle",
	})
}

func (h *DashboardBankHandler) CreateBank(c *fiber.Ctx) error {
	if err := requests.ValidateBankRequest(c); err != nil {
		return err
	}
	req := c.Locals("bankRequest").(requests.BankRequest)
	bank := &models.Bank{
		Name:     req.Name,
		IsActive: req.IsActive == "true",
	}
	if err := h.bankService.CreateBank(c.UserContext(), bank); err != nil {
		return renderBankFormError("dashboard/banks/create", "Yeni Banka Ekle", req, "Banka oluşturulamadı: "+err.Error(), c)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Banka başarıyla oluşturuldu.")
	return c.Redirect("/dashboard/banks", fiber.StatusFound)
}

func (h *DashboardBankHandler) ShowUpdateBank(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	bank, err := h.bankService.GetBankByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Banka bulunamadı.")
		return c.Redirect("/dashboard/banks", fiber.StatusSeeOther)
	}
	return renderer.Render(c, "dashboard/banks/update", "layouts/dashboard", fiber.Map{
		"Title": "Banka Düzenle",
		"Bank":  bank,
	})
}

func (h *DashboardBankHandler) UpdateBank(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := requests.ValidateBankRequest(c); err != nil {
		return err
	}
	req := c.Locals("bankRequest").(requests.BankRequest)
	bank := &models.Bank{
		Name:     req.Name,
		IsActive: req.IsActive == "true",
	}
	// Get userID from context
	userID, _ := c.Locals("userID").(uint)
	if err := h.bankService.UpdateBank(c.UserContext(), uint(id), bank, userID); err != nil {
		return renderBankFormError("dashboard/banks/update", "Banka Güncelle", req, "Banka güncellenemedi: "+err.Error(), c)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Banka başarıyla güncellendi.")
	return c.Redirect("/dashboard/banks", fiber.StatusFound)
}

func (h *DashboardBankHandler) DeleteBank(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	bankID := uint(id)

	if err := h.bankService.DeleteBank(c.UserContext(), bankID); err != nil {
		errMsg := "Banka silinemedi: " + err.Error()
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/dashboard/banks", fiber.StatusSeeOther)
	}

	if strings.Contains(c.Get("Accept"), "application/json") {
		return c.JSON(fiber.Map{"message": "Banka başarıyla silindi."})
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Banka başarıyla silindi.")
	return c.Redirect("/dashboard/banks", fiber.StatusFound)
}

func renderBankFormError(template string, title string, req any, message string, c *fiber.Ctx) error {
	return renderer.Render(c, template, "layouts/dashboard", fiber.Map{
		"Title":                    title,
		renderer.FlashErrorKeyView: message,
		renderer.FormDataKey:       req,
	}, http.StatusBadRequest)
}
