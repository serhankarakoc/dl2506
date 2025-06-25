package repositories

import (
	"context"
	"davet.link/configs/databaseconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"
	"gorm.io/gorm"
)

type ICardRepository interface {
	GetAllCards(params queryparams.ListParams) ([]models.Card, int64, error)
	GetCardByID(id uint) (*models.Card, error)
	CreateCardWithRelations(ctx context.Context, card *models.Card) error
	UpdateCardWithRelations(ctx context.Context, card *models.Card) error
	DeleteCardWithRelations(ctx context.Context, id uint) error
	GetCardCount() (int64, error)
	IsSlugAvailable(slug string, excludeID uint) (bool, error)
}

type CardRepository struct {
	base IBaseRepository[models.Card]
	db   *gorm.DB
}

func NewCardRepository() ICardRepository {
	base := NewBaseRepository[models.Card](databaseconfig.GetDB())
	base.SetAllowedSortColumns([]string{"id", "name", "slug", "is_active", "created_at"})
	base.SetPreloads("CardBanks.Bank", "CardSocialMedia.SocialMedia")
	return &CardRepository{base: base, db: databaseconfig.GetDB()}
}

func (r *CardRepository) GetAllCards(params queryparams.ListParams) ([]models.Card, int64, error) {
	return r.base.GetAll(params)
}

func (r *CardRepository) GetCardByID(id uint) (*models.Card, error) {
	return r.base.GetByID(id)
}

func (r *CardRepository) CreateCardWithRelations(ctx context.Context, card *models.Card) error {
	return r.base.CreateWithRelations(ctx, card)
}

func (r *CardRepository) UpdateCardWithRelations(ctx context.Context, card *models.Card) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var currentBankIDs []uint
	for _, bank := range card.CardBanks {
		if bank.ID != 0 {
			currentBankIDs = append(currentBankIDs, bank.ID)
		}
	}

	if len(currentBankIDs) > 0 {
		if err := tx.Where("card_id = ? AND id NOT IN ?", card.ID, currentBankIDs).Delete(&models.CardBank{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := tx.Where("card_id = ?", card.ID).Delete(&models.CardBank{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	var currentSocialMediaIDs []uint
	for _, sm := range card.CardSocialMedia {
		if sm.ID != 0 {
			currentSocialMediaIDs = append(currentSocialMediaIDs, sm.ID)
		}
	}

	if len(currentSocialMediaIDs) > 0 {
		if err := tx.Where("card_id = ? AND id NOT IN ?", card.ID, currentSocialMediaIDs).Delete(&models.CardSocialMedia{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := tx.Where("card_id = ?", card.ID).Delete(&models.CardSocialMedia{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Save(card).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *CardRepository) DeleteCardWithRelations(ctx context.Context, id uint) error {
	return r.base.DeleteWithRelations(ctx, id)
}

func (r *CardRepository) GetCardCount() (int64, error) {
	return r.base.GetCount()
}

func (r *CardRepository) IsSlugAvailable(slug string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Card{}).Where("slug = ?", slug)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}