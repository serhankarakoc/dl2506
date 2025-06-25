package sessionconfig

import (
	"encoding/gob"
	"time"

	"davet.link/configs/envconfig"
	"davet.link/configs/logconfig"
	"davet.link/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Session *session.Store

func InitSession() {
	Session = createSessionStore()
	registerGobTypes()
	logconfig.SLog.Info("Oturum (session) sistemi başlatıldı ve utils içinde kayıt edildi.")
}

func SetupSession() *session.Store {
	if Session == nil {
		logconfig.SLog.Warn("Session store isteniyor ancak henüz başlatılmamış, şimdi başlatılıyor.")
		InitSession()
	}
	return Session
}

func createSessionStore() *session.Store {
	sessionExpirationHours := envconfig.GetEnvAsInt("SESSION_EXPIRATION_HOURS", 24)
	cookieSecure := envconfig.IsProduction()

	store := session.New(session.Config{
		CookieHTTPOnly: false,
		CookieSecure:   cookieSecure,
		Expiration:     time.Duration(sessionExpirationHours) * time.Hour,
		KeyLookup:      "cookie:session_id",
		CookieSameSite: "Lax",
	})

	logconfig.SLog.Info("Cookie tabanlı session sistemi %d saatlik süreyle yapılandırıldı.", sessionExpirationHours)
	return store
}

func registerGobTypes() {
	gob.Register(models.UserType(""))
	gob.Register(&models.User{})
	logconfig.SLog.Debug("Session için gob türleri kaydedildi: models.UserType, *models.User")
}

func SessionStart(c *fiber.Ctx) (*session.Session, error) {
	if Session == nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "session store not initialized")
	}
	return Session.Get(c)
}

func DestroySession(c *fiber.Ctx) error {
	sess, err := SessionStart(c)
	if err != nil {
		return err
	}
	return sess.Destroy()
}

func GetUserTypeFromSession(sess *session.Session) (models.UserType, error) {
	userType, ok := sess.Get("user_type").(models.UserType)
	if !ok {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Geçersiz oturum veya kullanıcı tipi")
	}
	return userType, nil
}

func GetUserIDFromSession(c *fiber.Ctx) (uint, error) {
	sess, err := SessionStart(c)
	if err != nil {
		return 0, err
	}
	userIDValue := sess.Get("user_id")
	switch v := userIDValue.(type) {
	case uint:
		return v, nil
	case int:
		return uint(v), nil
	case float64:
		return uint(v), nil
	default:
		return 0, fiber.ErrUnauthorized
	}
}

func GetUserStatusFromSession(sess *session.Session) (bool, error) {
	userStatus, ok := sess.Get("user_status").(bool)
	if !ok {
		return false, fiber.NewError(fiber.StatusUnauthorized, "Geçersiz oturum veya kullanıcı durumu")
	}
	return userStatus, nil
}

func SetSessionValue(c *fiber.Ctx, key string, value interface{}) error {
	sess, err := SessionStart(c)
	if err != nil {
		return err
	}
	sess.Set(key, value)
	return sess.Save()
}
