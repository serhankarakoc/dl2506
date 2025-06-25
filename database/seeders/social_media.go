package seeders

import (
	"davet.link/configs/logconfig"
	"davet.link/models"
	"gorm.io/gorm"
)

func SeedSocialMedia(db *gorm.DB) error {
	// Sosyal medya listesi
	socialMedias := []models.SocialMedia{
		{Name: "Facebook", Icon: "fa-brands fa-facebook-f", IsActive: true},
		{Name: "Instagram", Icon: "fa-brands fa-instagram", IsActive: true},
		{Name: "X (Twitter)", Icon: "fa-brands fa-x", IsActive: true},
		{Name: "LinkedIn", Icon: "fa-brands fa-linkedin-in", IsActive: true},
		{Name: "YouTube", Icon: "fa-brands fa-youtube", IsActive: true},
		{Name: "TikTok", Icon: "fa-brands fa-tiktok", IsActive: true},
		{Name: "GitHub", Icon: "fa-brands fa-github", IsActive: true},
	}

	logconfig.SLog.Info("Sosyal medya verileri yükleniyor...")

	// Her bir sosyal medya platformu için
	for _, socialMedia := range socialMedias {
		// Platform zaten var mı kontrol et
		var existingSocialMedia models.SocialMedia
		if err := db.Where("name = ?", socialMedia.Name).First(&existingSocialMedia).Error; err == gorm.ErrRecordNotFound {
			// Platform yoksa ekle
			if err := db.Create(&socialMedia).Error; err != nil {
				logconfig.SLog.Error("Sosyal medya platformu eklenirken hata: " + socialMedia.Name)
				return err
			}
			logconfig.SLog.Info("Sosyal medya platformu eklendi: " + socialMedia.Name)
		}
	}

	logconfig.SLog.Info("Sosyal medya verileri yükleme işlemi tamamlandı.")
	return nil
}
