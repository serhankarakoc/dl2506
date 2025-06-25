package routes

import (
	handlers "davet.link/handlers/dashboard"
	"davet.link/middlewares"
	"davet.link/models"

	"github.com/gofiber/fiber/v2"
)

func registerDashboardRoutes(app *fiber.App) {
	dashboardGroup := app.Group("/dashboard")
	dashboardGroup.Use(
		middlewares.AuthMiddleware,
		middlewares.StatusMiddleware,
		middlewares.TypeMiddleware(models.Dashboard),
	)

	dashboardHomeHandler := handlers.NewDashboardHomeHandler()
	dashboardGroup.Get("/home", dashboardHomeHandler.HomePage)

	userHandler := handlers.NewDashboardUserHandler()
	dashboardGroup.Get("/users", userHandler.ListUsers)
	dashboardGroup.Get("/users/create", userHandler.ShowCreateUser)
	dashboardGroup.Post("/users/create", userHandler.CreateUser)
	dashboardGroup.Get("/users/update/:id", userHandler.ShowUpdateUser)
	dashboardGroup.Post("/users/update/:id", userHandler.UpdateUser)
	dashboardGroup.Delete("/users/delete/:id", userHandler.DeleteUser)

	invitationCategoryHandler := handlers.NewDashboardInvitationCategoryHandler()
	dashboardGroup.Get("/invitation-categories", invitationCategoryHandler.ListCategories)
	dashboardGroup.Get("/invitation-categories/create", invitationCategoryHandler.ShowCreateCategory)
	dashboardGroup.Post("/invitation-categories/create", invitationCategoryHandler.CreateCategory)
	dashboardGroup.Get("/invitation-categories/update/:id", invitationCategoryHandler.ShowUpdateCategory)
	dashboardGroup.Post("/invitation-categories/update/:id", invitationCategoryHandler.UpdateCategory)
	dashboardGroup.Delete("/invitation-categories/delete/:id", invitationCategoryHandler.DeleteCategory)

	bankHandler := handlers.NewDashboardBankHandler()
	dashboardGroup.Get("/banks", bankHandler.ListBanks)
	dashboardGroup.Get("/banks/create", bankHandler.ShowCreateBank)
	dashboardGroup.Post("/banks/create", bankHandler.CreateBank)
	dashboardGroup.Get("/banks/update/:id", bankHandler.ShowUpdateBank)
	dashboardGroup.Post("/banks/update/:id", bankHandler.UpdateBank)
	dashboardGroup.Delete("/banks/delete/:id", bankHandler.DeleteBank)

	socialMediaHandler := handlers.NewDashboardSocialMediaHandler()
	dashboardGroup.Get("/social-media", socialMediaHandler.ListSocialMedias)
	dashboardGroup.Get("/social-media/create", socialMediaHandler.ShowCreateSocialMedia)
	dashboardGroup.Post("/social-media/create", socialMediaHandler.CreateSocialMedia)
	dashboardGroup.Get("/social-media/update/:id", socialMediaHandler.ShowUpdateSocialMedia)
	dashboardGroup.Post("/social-media/update/:id", socialMediaHandler.UpdateSocialMedia)
	dashboardGroup.Delete("/social-media/delete/:id", socialMediaHandler.DeleteSocialMedia)

	cardHandler := handlers.NewDashboardCardHandler()
	dashboardGroup.Get("/cards", cardHandler.ListCards)
	dashboardGroup.Get("/cards/create", cardHandler.ShowCreateCard)
	dashboardGroup.Post("/cards/create", cardHandler.CreateCard)
	dashboardGroup.Get("/cards/update/:id", cardHandler.ShowUpdateCard)
	dashboardGroup.Post("/cards/update/:id", cardHandler.UpdateCard)
	dashboardGroup.Delete("/cards/delete/:id", cardHandler.DeleteCard)
	dashboardGroup.Get("/cards/slug-check", cardHandler.SlugCheck)

	invitationHandler := handlers.NewDashboardInvitationHandler()
	dashboardGroup.Get("/invitations", invitationHandler.ListInvitations)
	dashboardGroup.Get("/invitations/create", invitationHandler.ShowCreateInvitation)
	dashboardGroup.Post("/invitations/create", invitationHandler.CreateInvitation)
	dashboardGroup.Get("/invitations/update/:id", invitationHandler.ShowUpdateInvitation)
	dashboardGroup.Post("/invitations/update/:id", invitationHandler.UpdateInvitation)
	dashboardGroup.Delete("/invitations/delete/:id", invitationHandler.DeleteInvitation)
	dashboardGroup.Get("/invitations/participants/:id", invitationHandler.ListParticipants)
}
