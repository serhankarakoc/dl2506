package seeders

import (
	"davet.link/configs/logconfig"
	"davet.link/models"
	"gorm.io/gorm"
)

func SeedInvitationCategories(db *gorm.DB) error {
	// Kategori listesi
	categories := []models.InvitationCategory{
		{Name: "Açılış Davetiyesi", Template: "title", Icon: "fa-solid fa-shop", IsActive: true},
		{Name: "After Party Davetiyesi", Template: "title", Icon: "fa-solid fa-champagne-glasses", IsActive: true},
		{Name: "Asker Eğlencesi Davetiyesi", Template: "person-family", Icon: "fa-solid fa-helmet-safety", IsActive: true},
		{Name: "Baby Shower Davetiyesi", Template: "person-family", Icon: "fa-solid fa-baby", IsActive: true},
		{Name: "Balo Davetiyesi", Template: "title", Icon: "fa-solid fa-mask", IsActive: true},
		{Name: "Bekarlığa Veda Partisi Davetiyesi", Template: "person", Icon: "fa-solid fa-martini-glass", IsActive: true},
		{Name: "Cinsiyet Partisi Davetiyesi", Template: "person-family", Icon: "fa-solid fa-venus-mars", IsActive: true},
		{Name: "Defile Davetiyesi", Template: "title", Icon: "fa-solid fa-person-walking", IsActive: true},
		{Name: "Dini Tören Davetiyesi", Template: "title", Icon: "fa-solid fa-mosque", IsActive: true},
		{Name: "Doğum Günü Davetiyesi", Template: "person", Icon: "fa-solid fa-cake-candles", IsActive: true},
		{Name: "Düğün Davetiyesi", Template: "wedding", Icon: "fa-solid fa-rings-wedding", IsActive: true},
		{Name: "Film Galası Davetiyesi", Template: "title", Icon: "fa-solid fa-film", IsActive: true},
		{Name: "Gelin Hamamı Davetiyesi", Template: "person", Icon: "fa-solid fa-spray-can-sparkles", IsActive: true},
		{Name: "Kına Gecesi Davetiyesi", Template: "wedding", Icon: "fa-solid fa-moon", IsActive: true},
		{Name: "Kongre Davetiyesi", Template: "title", Icon: "fa-solid fa-users-rectangle", IsActive: true},
		{Name: "Konser Davetiyesi", Template: "title", Icon: "fa-solid fa-music", IsActive: true},
		{Name: "Lansman Davetiyesi", Template: "title", Icon: "fa-solid fa-rocket", IsActive: true},
		{Name: "Mezuniyet Töreni Davetiyesi", Template: "title", Icon: "fa-solid fa-graduation-cap", IsActive: true},
		{Name: "Nikâh Töreni Davetiyesi", Template: "wedding", Icon: "fa-solid fa-ring", IsActive: true},
		{Name: "Nişan Davetiyesi", Template: "wedding", Icon: "fa-solid fa-hand-holding-heart", IsActive: true},
		{Name: "Seminer Davetiyesi", Template: "title", Icon: "fa-solid fa-chalkboard-user", IsActive: true},
		{Name: "Sergi Davetiyesi", Template: "title", Icon: "fa-solid fa-image", IsActive: true},
		{Name: "Spor Müsabakası Davetiyesi", Template: "title", Icon: "fa-solid fa-trophy", IsActive: true},
		{Name: "Sünnet Düğünü Davetiyesi", Template: "person-family", Icon: "fa-solid fa-child", IsActive: true},
		{Name: "Tanıtım Davetiyesi", Template: "title", Icon: "fa-solid fa-bullhorn", IsActive: true},
		{Name: "Tiyatro Gösterisi Davetiyesi", Template: "title", Icon: "fa-solid fa-masks-theater", IsActive: true},
		{Name: "Toplantı Davetiyesi", Template: "title", Icon: "fa-solid fa-handshake", IsActive: true},
		{Name: "Veda Partisi Davetiyesi", Template: "title", Icon: "fa-solid fa-door-open", IsActive: true},
		{Name: "Yılbaşı Partisi Davetiyesi", Template: "title", Icon: "fa-solid fa-calendar-star", IsActive: true},
		{Name: "Yıldönümü Davetiyesi", Template: "title", Icon: "fa-solid fa-calendar-heart", IsActive: true},
	}

	logconfig.SLog.Info("Davetiye kategorileri yükleniyor...")

	// Her bir kategori için
	for _, category := range categories {
		// Kategori zaten var mı kontrol et
		var existingCategory models.InvitationCategory
		if err := db.Where("name = ?", category.Name).First(&existingCategory).Error; err == gorm.ErrRecordNotFound {
			// Kategori yoksa ekle
			if err := db.Create(&category).Error; err != nil {
				logconfig.SLog.Error("Kategori eklenirken hata: " + category.Name)
				return err
			}
			logconfig.SLog.Info("Kategori eklendi: " + category.Name)
		}
	}

	logconfig.SLog.Info("Davetiye kategorileri yükleme işlemi tamamlandı.")
	return nil
}
