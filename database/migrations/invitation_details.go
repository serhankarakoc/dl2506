package migrations

import (
	"davet.link/configs/logconfig"
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateInvitationDetailsTable(db *gorm.DB) error {
	logconfig.SLog.Info("InvitationDetail tablosu migrate ediliyor...")
	if err := db.AutoMigrate(&models.InvitationDetail{}); err != nil {
		return err
	}
	logconfig.SLog.Info("InvitationDetail tablosu migrate işlemi tamamlandı.")
	return nil
}
