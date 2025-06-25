package repositories

import (
	"context"

	"davet.link/configs/databaseconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"

	"gorm.io/gorm"
)

type IBankRepository interface {
	GetAllBanks(params queryparams.ListParams) ([]models.Bank, int64, error)
	GetBankByID(id uint) (*models.Bank, error)
	CreateBank(ctx context.Context, bank *models.Bank) error
	BulkCreateBanks(ctx context.Context, banks []models.Bank) error
	UpdateBank(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error
	BulkUpdateBanks(ctx context.Context, condition map[string]interface{}, data map[string]interface{}, updatedBy uint) error
	DeleteBank(ctx context.Context, id uint) error
	BulkDeleteBanks(ctx context.Context, condition map[string]interface{}) error
	GetBankCount() (int64, error)
}

type BankRepository struct {
	base IBaseRepository[models.Bank]
	db   *gorm.DB
}

func NewBankRepository() IBankRepository {
	base := NewBaseRepository[models.Bank](databaseconfig.GetDB())
	base.SetAllowedSortColumns([]string{"id", "name", "is_active", "created_at"})
	return &BankRepository{base: base, db: databaseconfig.GetDB()}
}

func (r *BankRepository) GetAllBanks(params queryparams.ListParams) ([]models.Bank, int64, error) {
	return r.base.GetAll(params)
}

func (r *BankRepository) GetBankByID(id uint) (*models.Bank, error) {
	return r.base.GetByID(id)
}

func (r *BankRepository) CreateBank(ctx context.Context, bank *models.Bank) error {
	return r.base.Create(ctx, bank)
}

func (r *BankRepository) BulkCreateBanks(ctx context.Context, banks []models.Bank) error {
	return r.base.BulkCreate(ctx, banks)
}

func (r *BankRepository) UpdateBank(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error {
	return r.base.Update(ctx, id, data, updatedBy)
}

func (r *BankRepository) BulkUpdateBanks(ctx context.Context, condition map[string]interface{}, data map[string]interface{}, updatedBy uint) error {
	return r.base.BulkUpdate(ctx, condition, data, updatedBy)
}

func (r *BankRepository) DeleteBank(ctx context.Context, id uint) error {
	return r.base.Delete(ctx, id)
}

func (r *BankRepository) BulkDeleteBanks(ctx context.Context, condition map[string]interface{}) error {
	return r.base.BulkDelete(ctx, condition)
}

func (r *BankRepository) GetBankCount() (int64, error) {
	return r.base.GetCount()
}

var _ IBankRepository = (*BankRepository)(nil)
var _ IBaseRepository[models.Bank] = (*BaseRepository[models.Bank])(nil)
