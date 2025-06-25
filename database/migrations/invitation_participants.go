package migrations

import (
	"davet.link/configs/logconfig"
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateInvitationParticipantsTable(db *gorm.DB) error {
	logconfig.SLog.Info("InvitationParticipant tablosu migrate ediliyor...")
	if err := db.AutoMigrate(&models.InvitationParticipant{}); err != nil {
		return err
	}
	logconfig.SLog.Info("InvitationParticipant tablosu migrate işlemi tamamlandı.")
	return nil
}
