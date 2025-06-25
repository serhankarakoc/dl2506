package dashboard

import (
	"context"
	"davet.link/models"
	"davet.link/pkg/filemanager"
	"davet.link/pkg/flashmessages"
	"davet.link/pkg/queryparams"
	"davet.link/pkg/renderer"
	"davet.link/requests"
	"davet.link/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DashboardInvitationHandler struct {
	invitationService         services.IInvitationService
	invitationCategoryService services.IInvitationCategoryService
}

func NewDashboardInvitationHandler() *DashboardInvitationHandler {
	return &DashboardInvitationHandler{
		invitationService:         services.NewInvitationService(),
		invitationCategoryService: services.NewInvitationCategoryService(),
	}
}

func (h *DashboardInvitationHandler) ListInvitations(c *fiber.Ctx) error {
	var params queryparams.ListParams
	if err := c.QueryParser(&params); err != nil {
		logconfig.Log.Warn("Davetiye listesi: Query parametreleri parse edilemedi, varsayılanlar kullanılıyor.", zap.Error(err))
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

	result, err := h.invitationService.GetAllInvitations(params)
	renderData := fiber.Map{"Title": "Davetiyeler", "Result": result, "Params": params}

	if err != nil {
		renderData[renderer.FlashErrorKeyView] = "Davetiyeler getirilirken bir hata oluştu."
	}

	return renderer.Render(c, "dashboard/invitations/list", "layouts/dashboard", renderData)
}

func (h *DashboardInvitationHandler) ShowCreateInvitation(c *fiber.Ctx) error {
	categories, err := h.invitationCategoryService.GetAllCategories()
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kategoriler getirilemedi.")
	}
	return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
		"Title":      "Yeni Davetiye Oluştur",
		"Categories": categories,
	})
}

func (h *DashboardInvitationHandler) CreateInvitation(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)

	if err := requests.ValidateInvitationRequest(c); err != nil {
		req, _ := c.Locals("invitationRequest").(requests.InvitationRequest)
		categories, _ := h.invitationCategoryService.GetAllCategories()
		return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
			"Title": "Yeni Davetiye Oluştur", "Categories": categories, "FormData": req,
		})
	}

	req := c.Locals("invitationRequest").(requests.InvitationRequest)

	newFileName, err := filemanager.UploadFile(c, "image", "invitations")
	if err != nil {
		flashMsg := "Resim yüklenemedi: " + err.Error()
		if err == filemanager.ErrFileNotProvided {
			flashMsg = "Davetiye resmi yüklemek zorunludur."
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, flashMsg)
		categories, _ := h.invitationCategoryService.GetAllCategories()
		return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
			"Title": "Yeni Davetiye Oluştur", "Categories": categories, "FormData": req,
		})
	}

	date, _ := time.Parse("2006-01-02", req.Date)
	timeVal, _ := time.Parse("15:04", req.Time)
	eventDateTime := date.Add(time.Hour*time.Duration(timeVal.Hour()) + time.Minute*time.Duration(timeVal.Minute()))

	invitation := &models.Invitation{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Template:   req.Template,
		Type:       req.Type,
		Title:      req.Title,
		Image:      newFileName,
		Venue:      req.Venue,
		Address:    req.Address,
		Location:   req.Location,
		Telephone:  req.Telephone,
		Date:       eventDateTime,
	}

	detail := models.InvitationDetail{
		Title:              req.Detail.Title, BrideName: req.Detail.BrideName, BrideSurname: req.Detail.BrideSurname,
		BrideMotherName:    req.Detail.BrideMotherName, BrideMotherSurname: req.Detail.BrideMotherSurname,
		BrideFatherName:    req.Detail.BrideFatherName, BrideFatherSurname: req.Detail.BrideFatherSurname,
		GroomName:          req.Detail.GroomName, GroomSurname: req.Detail.GroomSurname, GroomMotherName: req.Detail.GroomMotherName,
		GroomMotherSurname: req.Detail.GroomMotherSurname, GroomFatherName: req.Detail.GroomFatherName, GroomFatherSurname: req.Detail.GroomFatherSurname,
		Person:             req.Detail.Person, MotherName: req.Detail.MotherName, MotherSurname: req.Detail.MotherSurname,
		FatherName:         req.Detail.FatherName, FatherSurname: req.Detail.FatherSurname,
		IsMotherLive:       req.Detail.IsMotherLive == "true", IsFatherLive: req.Detail.IsFatherLive == "true",
		IsBrideMotherLive:  req.Detail.IsBrideMotherLive == "true", IsBrideFatherLive: req.Detail.IsBrideFatherLive == "true",
		IsGroomMotherLive:  req.Detail.IsGroomMotherLive == "true", IsGroomFatherLive: req.Detail.IsGroomFatherLive == "true",
	}
	invitation.InvitationDetail = &detail

	if err := h.invitationService.CreateInvitation(c.UserContext(), invitation); err != nil {
		filemanager.DeleteFile("invitations", newFileName)
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye oluşturulamadı: "+err.Error())
		categories, _ := h.invitationCategoryService.GetAllCategories()
		return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
			"Title": "Yeni Davetiye Oluştur", "Categories": categories, "FormData": req,
		})
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla oluşturuldu.")
	return c.Redirect("/dashboard/invitations", http.StatusFound)
}

func (h *DashboardInvitationHandler) ShowUpdateInvitation(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz davetiye ID'si.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}

	invitation, err := h.invitationService.GetInvitationByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye bulunamadı.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}
	
	categories, _ := h.invitationCategoryService.GetAllCategories()
	
	return renderer.Render(c, "dashboard/invitations/update", "layouts/dashboard", fiber.Map{
		"Title":      "Davetiye Düzenle",
		"FormData":   invitation,
		"Categories": categories,
	})
}

func (h *DashboardInvitationHandler) UpdateInvitation(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz davetiye ID'si.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}
	
	redirectURL := fmt.Sprintf("/dashboard/invitations/update/%d", id)

	existingInvitation, err := h.invitationService.GetInvitationByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Güncellenecek davetiye bulunamadı.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}
	
	if err := requests.ValidateInvitationRequest(c); err != nil {
		req, _ := c.Locals("invitationRequest").(requests.InvitationRequest)
		categories, _ := h.invitationCategoryService.GetAllCategories()
		return renderer.Render(c, "dashboard/invitations/update", "layouts/dashboard", fiber.Map{
			"Title": "Davetiye Düzenle", "Categories": categories, "FormData": req,
		})
	}
	
	req := c.Locals("invitationRequest").(requests.InvitationRequest)

	newFileName, err := filemanager.UploadFile(c, "image", "invitations")
	if err != nil && err != filemanager.ErrFileNotProvided {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Yeni resim yüklenemedi: "+err.Error())
		return c.Redirect(redirectURL, http.StatusSeeOther)
	}
	
	var oldPhotoToDelete string
	if newFileName != "" {
		oldPhotoToDelete = existingInvitation.Image
		existingInvitation.Image = newFileName
	}
	
	date, _ := time.Parse("2006-01-02", req.Date)
	timeVal, _ := time.Parse("15:04", req.Time)
	eventDateTime := date.Add(time.Hour*time.Duration(timeVal.Hour()) + time.Minute*time.Duration(timeVal.Minute()))

	existingInvitation.CategoryID = req.CategoryID
	existingInvitation.Template = req.Template
	existingInvitation.Type = req.Type
	existingInvitation.Title = req.Title
	existingInvitation.Venue = req.Venue
	existingInvitation.Address = req.Address
	existingInvitation.Location = req.Location
	existingInvitation.Telephone = req.Telephone
	existingInvitation.Date = eventDateTime
	
	if existingInvitation.InvitationDetail != nil {
		existingInvitation.InvitationDetail.Title = req.Detail.Title
		existingInvitation.InvitationDetail.BrideName = req.Detail.BrideName
		existingInvitation.InvitationDetail.BrideSurname = req.Detail.BrideSurname
		existingInvitation.InvitationDetail.BrideMotherName = req.Detail.BrideMotherName
		existingInvitation.InvitationDetail.BrideMotherSurname = req.Detail.BrideMotherSurname
		existingInvitation.InvitationDetail.BrideFatherName = req.Detail.BrideFatherName
		existingInvitation.InvitationDetail.BrideFatherSurname = req.Detail.BrideFatherSurname
		existingInvitation.InvitationDetail.GroomName = req.Detail.GroomName
		existingInvitation.InvitationDetail.GroomSurname = req.Detail.GroomSurname
		existingInvitation.InvitationDetail.GroomMotherName = req.Detail.GroomMotherName
		existingInvitation.InvitationDetail.GroomMotherSurname = req.Detail.GroomMotherSurname
		existingInvitation.InvitationDetail.GroomFatherName = req.Detail.GroomFatherName
		existingInvitation.InvitationDetail.GroomFatherSurname = req.Detail.GroomFatherSurname
		existingInvitation.InvitationDetail.Person = req.Detail.Person
		existingInvitation.InvitationDetail.MotherName = req.Detail.MotherName
		existingInvitation.InvitationDetail.MotherSurname = req.Detail.MotherSurname
		existingInvitation.InvitationDetail.FatherName = req.Detail.FatherName
		existingInvitation.InvitationDetail.FatherSurname = req.Detail.FatherSurname
		existingInvitation.InvitationDetail.IsMotherLive = req.Detail.IsMotherLive == "true"
		existingInvitation.InvitationDetail.IsFatherLive = req.Detail.IsFatherLive == "true"
		existingInvitation.InvitationDetail.IsBrideMotherLive = req.Detail.IsBrideMotherLive == "true"
		existingInvitation.InvitationDetail.IsBrideFatherLive = req.Detail.IsBrideFatherLive == "true"
		existingInvitation.InvitationDetail.IsGroomMotherLive = req.Detail.IsGroomMotherLive == "true"
		existingInvitation.InvitationDetail.IsGroomFatherLive = req.Detail.IsGroomFatherLive == "true"
	}

	if err := h.invitationService.UpdateInvitation(c.UserContext(), existingInvitation); err != nil {
		if newFileName != "" {
			filemanager.DeleteFile("invitations", newFileName)
		}
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye güncellenirken bir hata oluştu.")
		return c.Redirect(redirectURL, http.StatusSeeOther)
	}
	
	if oldPhotoToDelete != "" {
		filemanager.DeleteFile("invitations", oldPhotoToDelete)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla güncellendi.")
	return c.Redirect("/dashboard/invitations", http.StatusFound)
}

func (h *DashboardInvitationHandler) DeleteInvitation(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.invitationService.DeleteInvitation(c.UserContext(), uint(id)); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye silinemedi: "+err.Error())
		return c.Redirect("/dashboard/invitations", fiber.StatusSeeOther)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla silindi.")
	return c.Redirect("/dashboard/invitations", http.StatusFound)
}