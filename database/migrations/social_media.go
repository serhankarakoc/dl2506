package migrations

import (
	"davet.link/configs/logconfig"
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateSocialMediaTable(db *gorm.DB) error {
	logconfig.SLog.Info("SocialMedia tablosu migrate ediliyor...")
	if err := db.AutoMigrate(&models.SocialMedia{}); err != nil {
		return err
	}
	logconfig.SLog.Info("SocialMedia tablosu migrate işlemi tamamlandı.")
	return nil
}
