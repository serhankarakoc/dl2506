package repositories

import (
	"context"
	"davet.link/configs/databaseconfig"
	"davet.link/models"
)

type ICardSocialMediaRepository interface {
	DeleteByCardID(ctx context.Context, cardID uint) error
	BulkCreate(ctx context.Context, cardSocialMedia []models.CardSocialMedia) error
}

type CardSocialMediaRepository struct{}

func NewCardSocialMediaRepository() ICardSocialMediaRepository {
	return &CardSocialMediaRepository{}
}

func (r *CardSocialMediaRepository) DeleteByCardID(ctx context.Context, cardID uint) error {
	db := databaseconfig.GetDB()
	return db.WithContext(ctx).Where("card_id = ?", cardID).Delete(&models.CardSocialMedia{}).Error
}

func (r *CardSocialMediaRepository) BulkCreate(ctx context.Context, cardSocialMedia []models.CardSocialMedia) error {
	if len(cardSocialMedia) == 0 {
		return nil
	}
	db := databaseconfig.GetDB()
	return db.WithContext(ctx).Create(&cardSocialMedia).Error
}
