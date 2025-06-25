package handlers

import (
	"net/http"

	"davet.link/configs/logconfig"
	"davet.link/pkg/renderer"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DashboardHomeHandler struct {
	userService services.IUserService
	cardService  services.ICardService
	invitationService services.IInvitationService
}

func NewDashboardHomeHandler() *DashboardHomeHandler {
	userSvc := services.NewUserService()
	cardSvc := services.NewCardService()
	invitationSvc := services.NewInvitationService()
	return &DashboardHomeHandler{
		userService:      userSvc,
		cardService:     cardSvc,
		invitationService: invitationSvc,
	}
}

func (h *DashboardHomeHandler) HomePage(c *fiber.Ctx) error {
	userCount, userErr := h.userService.GetUserCount()
	cardCount, cardErr := h.cardService.GetCardCount()
	invitationCount, invitationErr := h.invitationService.GetInvitationCount()
	if userErr != nil {
		logconfig.Log.Error("Anasayfa: Kullanıcı sayısı alınamadı", zap.Error(userErr))
		userCount = 0
	}
	if cardErr != nil {
		logconfig.Log.Error("Anasayfa: Kart sayısı alınamadı", zap.Error(cardErr))
		cardCount = 0
	}
	if invitationErr != nil {
		logconfig.Log.Error("Anasayfa: Davetiye sayısı alınamadı", zap.Error(invitationErr))
		invitationCount = 0
	}

	mapData := fiber.Map{
		"Title":         "Dashboard",
		"UserCount":     userCount,
		"CardCount":     cardCount,
		"InvitationCount": invitationCount,
	}
	return renderer.Render(c, "dashboard/home/home", "layouts/dashboard", mapData, http.StatusOK)
}
