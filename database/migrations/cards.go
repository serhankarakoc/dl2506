package migrations

import (
	"davet.link/configs/logconfig"
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateCardsTable(db *gorm.DB) error {
	logconfig.SLog.Info("Card tablosu migrate ediliyor...")
	if err := db.AutoMigrate(&models.Card{}); err != nil {
		return err
	}
	logconfig.SLog.Info("Card tablosu migrate işlemi tamamlandı.")
	return nil
}
