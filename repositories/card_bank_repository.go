package repositories

import (
	"context"
	"davet.link/configs/databaseconfig"
	"davet.link/models"
)

type ICardBankRepository interface {
	DeleteByCardID(ctx context.Context, cardID uint) error
	BulkCreate(ctx context.Context, cardBanks []models.CardBank) error
}

type CardBankRepository struct{}

func NewCardBankRepository() ICardBankRepository {
	return &CardBankRepository{}
}

func (r *CardBankRepository) DeleteByCardID(ctx context.Context, cardID uint) error {
	db := databaseconfig.GetDB()
	return db.WithContext(ctx).Where("card_id = ?", cardID).Delete(&models.CardBank{}).Error
}

func (r *CardBankRepository) BulkCreate(ctx context.Context, cardBanks []models.CardBank) error {
	if len(cardBanks) == 0 {
		return nil
	}
	db := databaseconfig.GetDB()
	return db.WithContext(ctx).Create(&cardBanks).Error
}
