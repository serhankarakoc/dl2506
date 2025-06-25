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

type DashboardUserHandler struct {
	userService services.IUserService
}

func NewDashboardUserHandler() *DashboardUserHandler {
	svc := services.NewUserService()
	return &DashboardUserHandler{userService: svc}
}

func (h *DashboardUserHandler) ListUsers(c *fiber.Ctx) error {
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

	paginatedResult, dbErr := h.userService.GetAllUsers(params)

	renderData := fiber.Map{
		"Title":  "Kullanıcılar",
		"Result": paginatedResult,
		"Params": params,
	}
	if dbErr != nil {
		renderData[renderer.FlashErrorKeyView] = "Kullanıcılar getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.User{},
			Meta: queryparams.PaginationMeta{
				CurrentPage: params.Page, PerPage: params.PerPage,
			},
		}
	}
	return renderer.Render(c, "dashboard/users/list", "layouts/dashboard", renderData, http.StatusOK)
}

func (h *DashboardUserHandler) ShowCreateUser(c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/users/create", "layouts/dashboard", fiber.Map{
		"Title": "Yeni Kullanıcı Ekle",
	})
}

func (h *DashboardUserHandler) CreateUser(c *fiber.Ctx) error {
	if err := requests.ValidateUserRequest(c); err != nil {
		return err
	}
	req := c.Locals("userRequest").(requests.UserRequest)
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Status:   req.Status == "true",
		Type:     models.UserType(req.Type),
		ResetToken: req.ResetToken,
		EmailVerified: req.EmailVerified == "true",
		VerificationToken: req.VerificationToken,
		Provider: req.Provider,
		ProviderID: req.ProviderID,
	}
	if user.Type != models.Dashboard && user.Type != models.Panel {
		return renderUserFormError("dashboard/users/create", "Yeni Kullanıcı Ekle", req, "Geçersiz kullanıcı tipi seçildi.", c)
	}
	if err := h.userService.CreateUser(c.UserContext(), user); err != nil {
		return renderUserFormError("dashboard/users/create", "Yeni Kullanıcı Ekle", req, "Kullanıcı oluşturulamadı: "+err.Error(), c)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kullanıcı başarıyla oluşturuldu.")
	return c.Redirect("/dashboard/users", fiber.StatusFound)
}

func (h *DashboardUserHandler) ShowUpdateUser(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı bulunamadı.")
		return c.Redirect("/dashboard/users", fiber.StatusSeeOther)
	}
	return renderer.Render(c, "dashboard/users/update", "layouts/dashboard", fiber.Map{
		"Title": "Kullanıcı Düzenle",
		"User":  user,
	})
}

func (h *DashboardUserHandler) UpdateUser(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	if err := requests.ValidateUserRequest(c); err != nil {
		return err
	}
	req := c.Locals("userRequest").(requests.UserRequest)
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Status:   req.Status == "true",
		Type:     models.UserType(req.Type),
		ResetToken: req.ResetToken,
		EmailVerified: req.EmailVerified == "true",
		VerificationToken: req.VerificationToken,
		Provider: req.Provider,
		ProviderID: req.ProviderID,
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if user.Type != models.Dashboard && user.Type != models.Panel {
		return renderUserFormError("dashboard/users/update", "Kullanıcı Düzenle", req, "Geçersiz kullanıcı tipi seçildi.", c)
	}
	// Get userID from context
	userID, _ := c.Locals("userID").(uint)
	if err := h.userService.UpdateUser(c.UserContext(), uint(id), user, userID); err != nil {
		return renderUserFormError("dashboard/users/update", "Kullanıcı Düzenle", req, "Kullanıcı güncellenemedi: "+err.Error(), c)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kullanıcı başarıyla güncellendi.")
	return c.Redirect("/dashboard/users", fiber.StatusFound)
}

func (h *DashboardUserHandler) DeleteUser(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	userID := uint(id)

	if err := h.userService.DeleteUser(c.UserContext(), userID); err != nil {
		errMsg := "Kullanıcı silinemedi: " + err.Error()
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/dashboard/users", fiber.StatusSeeOther)
	}

	if strings.Contains(c.Get("Accept"), "application/json") {
		return c.JSON(fiber.Map{"message": "Kullanıcı başarıyla silindi."})
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kullanıcı başarıyla silindi.")
	return c.Redirect("/dashboard/users", fiber.StatusFound)
}

func renderUserFormError(template string, title string, req any, message string, c *fiber.Ctx) error {
	return renderer.Render(c, template, "layouts/dashboard", fiber.Map{
		"Title":                    title,
		renderer.FlashErrorKeyView: message,
		renderer.FormDataKey:       req,
	}, http.StatusBadRequest)
}
