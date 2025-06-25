package seeders

import (
	"davet.link/configs/logconfig"
	"davet.link/models"
	"gorm.io/gorm"
)

func SeedBanks(db *gorm.DB) error {
	// Banka listesi
	banks := []models.Bank{
		{Name: "AKBANK T.A.Ş.", IsActive: true},
		{Name: "AKTİF YATIRIM BANKASI A.Ş.", IsActive: true},
		{Name: "AHLATCI ÖDEME VE ELEKTRONİK PARA HİZMETLERİ A.Ş.", IsActive: true},
		{Name: "ALBARAKA TÜRK KATILIM BANKASI A.Ş.", IsActive: true},
		{Name: "ALTERNATİFBANK A.Ş.", IsActive: true},
		{Name: "ANADOLUBANK A.Ş.", IsActive: true},
		{Name: "BELBİM ELEKTRONİK PARA VE ÖDEME HİZMETLERİ A.Ş.", IsActive: true},
		{Name: "BURGAN BANK A.Ş.", IsActive: true},
		{Name: "DENİZBANK A.Ş.", IsActive: true},
		{Name: "DÜNYA KATILIM BANKASI A.Ş.", IsActive: true},
		{Name: "ENPARA BANK A.Ş.", IsActive: true},
		{Name: "FİBABANKA A.Ş.", IsActive: true},
		{Name: "GOLDEN GLOBAL YATIRIM BANKASI A.Ş.", IsActive: true},
		{Name: "HAYAT FİNANS KATILIM BANKASI A.Ş.", IsActive: true},
		{Name: "ING BANK A.Ş.", IsActive: true},
		{Name: "İNİNAL ÖDEME VE ELEKTRONİK PARA HİZMETLERİ A.Ş.", IsActive: true},
		{Name: "İYZİ ÖDEME VE ELEKTRONİK PARA HİZMETLERİ A.Ş.", IsActive: true},
		{Name: "KUVEYT TÜRK KATILIM BANKASI A.Ş.", IsActive: true},
		{Name: "LYDIANS ELEKTRONİK PARA VE ÖDEME HİZMETLERİ A.Ş.", IsActive: true},
		{Name: "MİSYON YATIRIM BANKASI A.Ş.", IsActive: true},
		{Name: "MOKA UNİTED ÖDEME HİZMETLERİ VE ELEKTRONİK PARA KURULUŞU A.Ş.", IsActive: true},
		{Name: "ODEA BANK A.Ş.", IsActive: true},
		{Name: "PAPARA ELEKTRONİK PARA A.Ş.", IsActive: true},
		{Name: "PAROLAPARA ELEKTRONİK PARA VE ÖDEME HİZMETLERİ A.Ş.", IsActive: true},
		{Name: "PAY FİX ELEKTRONİK PARA VE ÖDEME HİZMETLERİ A.Ş.", IsActive: true},
		{Name: "POSTA VE TELGRAF TEŞKİLATI A.Ş.", IsActive: true},
		{Name: "QNB BANK A.Ş.", IsActive: true},
		{Name: "SİPAY ELEKTRONİK PARA VE ÖDEME HİZMETLERİ A.Ş.", IsActive: true},
		{Name: "ŞEKERBANK T.A.Ş.", IsActive: true},
		{Name: "T.C. ZİRAAT BANKASI A.Ş.", IsActive: true},
		{Name: "T. EKONOMİ BANKASI A.Ş.", IsActive: true},
		{Name: "T. GARANTİ BANKASI A.Ş.", IsActive: true},
		{Name: "T. HALK BANKASI A.Ş.", IsActive: true},
		{Name: "T. İŞ BANKASI A.Ş.", IsActive: true},
		{Name: "T.O.M. KATILIM BANKASI A.Ş.", IsActive: true},
		{Name: "T. VAKIFLAR BANKASI T.A.O.", IsActive: true},
		{Name: "TURK ELEKTRONİK PARA A.Ş.", IsActive: true},
		{Name: "TURKCELL ÖDEME VE ELEKTRONİK PARA HİZMETLERİ A.Ş.", IsActive: true},
		{Name: "TÜRKİYE EMLAK KATILIM BANKASI A.Ş.", IsActive: true},
		{Name: "TÜRKİYE FİNANS KATILIM BANKASI A.Ş.", IsActive: true},
		{Name: "VAKIF KATILIM BANKASI A.Ş.", IsActive: true},
		{Name: "YAPI VE KREDİ BANKASI A.Ş.", IsActive: true},
		{Name: "ZİRAAT KATILIM BANKASI A.Ş.", IsActive: true},
	}

	logconfig.SLog.Info("Banka verileri yükleniyor...")

	// Her bir banka için
	for _, bank := range banks {
		// Banka zaten var mı kontrol et
		var existingBank models.Bank
		if err := db.Where("name = ?", bank.Name).First(&existingBank).Error; err == gorm.ErrRecordNotFound {
			// Banka yoksa ekle
			if err := db.Create(&bank).Error; err != nil {
				logconfig.SLog.Error("Banka eklenirken hata: " + bank.Name)
				return err
			}
			logconfig.SLog.Info("Banka eklendi: " + bank.Name)
		}
	}

	logconfig.SLog.Info("Banka verileri yükleme işlemi tamamlandı.")
	return nil
}
