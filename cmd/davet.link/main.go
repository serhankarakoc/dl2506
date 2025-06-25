package main

import (
	"os"
	"os/signal"
	"syscall"

	"davet.link/configs/csrfconfig"
	"davet.link/configs/databaseconfig"
	"davet.link/configs/fileconfig"
	"davet.link/configs/logconfig"
	"davet.link/configs/sessionconfig"
	"davet.link/pkg/flashmessages"
	"davet.link/pkg/templatehelpers"
	"davet.link/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	logconfig.InitLogger()
	defer logconfig.SyncLogger()

	logconfig.SLog.Debugw("Ortam değişkenleri yüklendi ve logger başlatıldı")

	databaseconfig.InitDB()
	defer databaseconfig.CloseDB()

	sessionconfig.InitSession()

	fileconfig.InitFileConfig()

	fileconfig.Config.SetAllowedExtensions("cards", []string{"jpg", "png", "webp"})
	fileconfig.Config.SetAllowedExtensions("invitations", []string{"jpeg", "png"})

	engine := html.New("./views", ".html")
	engine.AddFunc("getFlashMessages", flashmessages.GetFlashMessages)
	engine.AddFuncMap(templatehelpers.TemplateHelpers())

	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			logconfig.Log.Error("Fiber request error",
				zap.Error(err),
				zap.Int("status_code", code),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.String("ip", c.IP()),
			)

			return c.Status(code).JSON(fiber.Map{"error": message})
		},
	})

	app.Static("/", "./public")
	app.Static("/uploads", "./uploads")
	app.Use(csrfconfig.SetupCSRF())
	routes.SetupRoutes(app)

	startServer(app)
}

func startServer(app *fiber.App) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		port := os.Getenv("APP_PORT")
		if port == "" {
			port = "3000"
		}
		address := ":" + port
		fullAddress := "http://localhost" + address

		logconfig.Log.Info("Uygulama başlatılıyor",
			zap.String("address", fullAddress),
			zap.String("port", port),
		)

		if err := app.Listen(address); err != nil {
			logconfig.Log.Fatal("Sunucu dinlenemedi",
				zap.String("address", address),
				zap.Error(err),
			)
		}
	}()

	<-shutdown
	logconfig.Log.Info("Kapatma sinyali alındı, uygulama kapatılıyor...")

	if err := app.Shutdown(); err != nil {
		logconfig.Log.Error("Sunucu kapatılırken hata oluştu", zap.Error(err))
	} else {
		logconfig.Log.Info("Sunucu başarıyla kapatıldı")
	}

	logconfig.Log.Info("Uygulama başarıyla sonlandırıldı.")
}
