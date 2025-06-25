package routes

import (
	handlers "davet.link/handlers/website"

	"github.com/gofiber/fiber/v2"
)

func registerWebsiteRoutes(app *fiber.App) {
	websiteHandler := handlers.NewWebsiteHandler()
	app.Get("/", websiteHandler.ShowHomePage)
	app.Get("/kullanim-sartlari", websiteHandler.ShowTermsOfUse)
	// Statik sayfalar için tek bir route
	app.Get("/:staticPageName", websiteHandler.ShowStaticPage)
	// Davetiye rotası (ör: /123asd1)
	app.Get("/:invitationKey", websiteHandler.ShowInvitation)
	// Kartvizit rotası (ör: /@serhan)
	app.Get("/@:cardSlug", websiteHandler.ShowCard)
}
