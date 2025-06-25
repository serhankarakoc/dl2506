package handlers

import (
	"net/http"

	"davet.link/pkg/renderer"

	"github.com/gofiber/fiber/v2"
)

type WebsiteHandler struct {
}

func NewWebsiteHandler() *WebsiteHandler {
	return &WebsiteHandler{}
}

func (h *WebsiteHandler) ShowHomePage(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/home", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) ShowTermsOfUse(c *fiber.Ctx) error {
	return renderer.Render(c, "website/terms_of_use", "layouts/website", fiber.Map{}, http.StatusOK)
}

func (h *WebsiteHandler) ShowStaticPage(c *fiber.Ctx) error {
	page := c.Params("staticPageName")
	template := "website/" + page
	return renderer.Render(c, template, "layouts/website", fiber.Map{}, http.StatusOK)
}

func (h *WebsiteHandler) ShowInvitation(c *fiber.Ctx) error {
	invitationKey := c.Params("invitationKey")
	// TODO: Davetiye verisini çek ve render et
	return renderer.Render(c, "website/invitation", "layouts/website", fiber.Map{"InvitationKey": invitationKey}, http.StatusOK)
}

func (h *WebsiteHandler) ShowCard(c *fiber.Ctx) error {
	cardSlug := c.Params("cardSlug")
	// TODO: Kartvizit verisini çek ve render et
	return renderer.Render(c, "website/card", "layouts/website", fiber.Map{"CardSlug": cardSlug}, http.StatusOK)
}
