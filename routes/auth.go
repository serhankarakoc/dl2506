package routes

import (
	handlers "davet.link/handlers/auth"
	"davet.link/middlewares"
	"davet.link/requests"

	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(app *fiber.App) {
	authHandler := handlers.NewAuthHandler()

	authGroup := app.Group("/auth")

	authGroup.Get("/login", authHandler.ShowLogin)
	authGroup.Post("/login", middlewares.GuestMiddleware, requests.ValidateLoginRequest, authHandler.Login)

	authGroup.Get("/logout", middlewares.AuthMiddleware, authHandler.Logout)
	authGroup.Get("/profile", middlewares.AuthMiddleware, authHandler.Profile)
	authGroup.Post("/profile/update-password", middlewares.AuthMiddleware, requests.ValidateUpdatePasswordRequest, authHandler.UpdatePassword)
	authGroup.Get("/register", authHandler.ShowRegister)
	authGroup.Post("/register", middlewares.GuestMiddleware, requests.ValidateRegisterRequest, authHandler.Register)
	authGroup.Get("/forgot-password", authHandler.ShowForgotPassword)
	authGroup.Post("/forgot-password", middlewares.GuestMiddleware, requests.ValidateForgotPasswordRequest, authHandler.ForgotPassword)
	authGroup.Get("/reset-password", authHandler.ShowResetPassword)
	authGroup.Post("/reset-password", middlewares.GuestMiddleware, requests.ValidateResetPasswordRequest, authHandler.ResetPassword)
	authGroup.Get("/verify-email", authHandler.VerifyEmail)
	authGroup.Get("/resend-verification", authHandler.ShowResendVerification)
	authGroup.Post("/resend-verification", requests.ValidateResendVerificationRequest, authHandler.ResendVerification)
	authGroup.Get("/google/login", handlers.GoogleLogin)
	authGroup.Get("/google/callback", handlers.GoogleCallback)
}
