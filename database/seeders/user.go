package seeders

import (
	"context"

	"davet.link/configs/logconfig"
	"davet.link/models"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetSystemUserConfig() models.User {
	return models.User{
		Name:     "davet.link",
		Email:    "davet.link@davet.link",
		Type:     models.Dashboard,
		Password: "davet.link",
	}
}

func SeedSystemUser(db *gorm.DB) error {
	systemUserConfig := GetSystemUserConfig()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(systemUserConfig.Password), bcrypt.DefaultCost)
	if err != nil {
		logconfig.Log.Error("Sistem kullanıcısının şifresi hash'lenirken hata oluştu",
			zap.String("email", systemUserConfig.Email),
			zap.Error(err),
		)
		return err
	}

	userToSeed := models.User{
		Name:          systemUserConfig.Name,
		Email:         systemUserConfig.Email,
		Type:          systemUserConfig.Type,
		Password:      string(hashedPassword),
		Status:        true,
		EmailVerified: true,
	}

	var existingUser models.User
	result := db.Where("email = ? AND type = ?", userToSeed.Email, userToSeed.Type).First(&existingUser)

	if result.Error == nil {
		logconfig.SLog.Info("Sistem kullanıcısı '%s' zaten mevcut. Güncelleme gerekip gerekmediği kontrol ediliyor...", userToSeed.Email)

		updateFields := make(map[string]interface{})
		needsUpdate := false

		if existingUser.Name != userToSeed.Name {
			updateFields["name"] = userToSeed.Name
			needsUpdate = true
		}
		if !existingUser.Status {
			updateFields["status"] = true
			needsUpdate = true
		}

		if needsUpdate {
			logconfig.SLog.Info("Mevcut sistem kullanıcısı '%s' güncelleniyor...", userToSeed.Email)

			ctx := context.WithValue(context.Background(), "user_id", uint(1))
			err := db.WithContext(ctx).Model(&existingUser).Updates(updateFields).Error
			if err != nil {
				logconfig.Log.Error("Mevcut sistem kullanıcısı güncellenemedi",
					zap.String("email", userToSeed.Email),
					zap.Error(err),
				)
				return err
			}
			logconfig.SLog.Info("Mevcut sistem kullanıcısı '%s' başarıyla güncellendi.", userToSeed.Email)
		} else {
			logconfig.SLog.Info("Mevcut sistem kullanıcısı '%s' için güncelleme gerekmiyor.", userToSeed.Email)
		}
		return nil

	} else if result.Error != gorm.ErrRecordNotFound {
		logconfig.Log.Error("Sistem kullanıcısı kontrol edilirken veritabanı hatası",
			zap.String("email", userToSeed.Email),
			zap.Error(result.Error),
		)
		return result.Error
	}

	logconfig.SLog.Info("Sistem kullanıcısı '%s' bulunamadı. Oluşturuluyor...", userToSeed.Email)

	ctx := context.WithValue(context.Background(), "user_id", uint(1))
	err = db.WithContext(ctx).Create(&userToSeed).Error
	if err != nil {
		logconfig.Log.Error("Sistem kullanıcısı oluşturulamadı",
			zap.String("email", userToSeed.Email),
			zap.Error(err),
		)
		return err
	}

	logconfig.SLog.Info("Sistem kullanıcısı '%s' başarıyla oluşturuldu.", userToSeed.Email)
	return nil
}
