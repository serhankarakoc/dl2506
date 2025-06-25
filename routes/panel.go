package routes

import (
	handlers "davet.link/handlers/panel"
	"davet.link/middlewares"
	"davet.link/models"

	"github.com/gofiber/fiber/v2"
)

func registerPanelRoutes(app *fiber.App) {
	panelGroup := app.Group("/panel")
	panelGroup.Use(
		middlewares.AuthMiddleware,
		middlewares.StatusMiddleware,
		middlewares.TypeMiddleware(models.Panel),
		middlewares.VerifiedMiddleware,
	)

	panelGroup.Get("/home", handlers.PanelHomeHandler)

	panelCardHandler := handlers.NewPanelCardHandler()
	panelGroup.Get("/cards", panelCardHandler.ListCards)
	panelGroup.Get("/cards/create", panelCardHandler.ShowCreateCard)
	panelGroup.Post("/cards/create", panelCardHandler.CreateCard)
	panelGroup.Get("/cards/update/:id", panelCardHandler.ShowUpdateCard)
	panelGroup.Post("/cards/update/:id", panelCardHandler.UpdateCard)
	panelGroup.Delete("/cards/delete/:id", panelCardHandler.DeleteCard)

	panelInvitationHandler := handlers.NewPanelInvitationHandler()
	panelGroup.Get("/invitations", panelInvitationHandler.ListInvitations)
	panelGroup.Get("/invitations/create", panelInvitationHandler.ShowCreateInvitation)
	panelGroup.Post("/invitations/create", panelInvitationHandler.CreateInvitation)
	panelGroup.Get("/invitations/update/:id", panelInvitationHandler.ShowUpdateInvitation)
	panelGroup.Post("/invitations/update/:id", panelInvitationHandler.UpdateInvitation)
	panelGroup.Delete("/invitations/delete/:id", panelInvitationHandler.DeleteInvitation)
	panelGroup.Get("/invitations/participants/:id", panelInvitationHandler.ListParticipants)
}
