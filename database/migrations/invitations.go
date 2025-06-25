package migrations

import (
	"davet.link/configs/logconfig"
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateInvitationsTable(db *gorm.DB) error {
	logconfig.SLog.Info("Invitation tablosu migrate ediliyor...")
	if err := db.AutoMigrate(&models.Invitation{}); err != nil {
		return err
	}
	logconfig.SLog.Info("Invitation tablosu migrate işlemi tamamlandı.")
	return nil
}
