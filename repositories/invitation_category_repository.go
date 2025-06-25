package repositories

import (
	"context"

	"davet.link/configs/databaseconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"

	"gorm.io/gorm"
)

type IInvitationCategoryRepository interface {
	GetAllCategories(params queryparams.ListParams) ([]models.InvitationCategory, int64, error)
	GetCategoryByID(id uint) (*models.InvitationCategory, error)
	CreateCategory(ctx context.Context, category *models.InvitationCategory) error
	BulkCreateCategories(ctx context.Context, categories []models.InvitationCategory) error
	UpdateCategory(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error
	BulkUpdateCategories(ctx context.Context, condition map[string]interface{}, data map[string]interface{}, updatedBy uint) error
	DeleteCategory(ctx context.Context, id uint) error
	BulkDeleteCategories(ctx context.Context, condition map[string]interface{}) error
	GetCategoryCount() (int64, error)
}

type InvitationCategoryRepository struct {
	base IBaseRepository[models.InvitationCategory]
	db   *gorm.DB
}

func NewInvitationCategoryRepository() IInvitationCategoryRepository {
	base := NewBaseRepository[models.InvitationCategory](databaseconfig.GetDB())
	base.SetAllowedSortColumns([]string{"id", "name", "is_active", "created_at"})

	return &InvitationCategoryRepository{base: base, db: databaseconfig.GetDB()}
}

func (r *InvitationCategoryRepository) GetAllCategories(params queryparams.ListParams) ([]models.InvitationCategory, int64, error) {
	return r.base.GetAll(params)
}

func (r *InvitationCategoryRepository) GetCategoryByID(id uint) (*models.InvitationCategory, error) {
	return r.base.GetByID(id)
}

func (r *InvitationCategoryRepository) CreateCategory(ctx context.Context, category *models.InvitationCategory) error {
	return r.base.Create(ctx, category)
}

func (r *InvitationCategoryRepository) BulkCreateCategories(ctx context.Context, categories []models.InvitationCategory) error {
	return r.base.BulkCreate(ctx, categories)
}

func (r *InvitationCategoryRepository) UpdateCategory(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error {
	return r.base.Update(ctx, id, data, updatedBy)
}

func (r *InvitationCategoryRepository) BulkUpdateCategories(ctx context.Context, condition map[string]interface{}, data map[string]interface{}, updatedBy uint) error {
	return r.base.BulkUpdate(ctx, condition, data, updatedBy)
}

func (r *InvitationCategoryRepository) DeleteCategory(ctx context.Context, id uint) error {
	return r.base.Delete(ctx, id)
}

func (r *InvitationCategoryRepository) BulkDeleteCategories(ctx context.Context, condition map[string]interface{}) error {
	return r.base.BulkDelete(ctx, condition)
}

func (r *InvitationCategoryRepository) GetCategoryCount() (int64, error) {
	return r.base.GetCount()
}

var _ IInvitationCategoryRepository = (*InvitationCategoryRepository)(nil)
var _ IBaseRepository[models.InvitationCategory] = (*BaseRepository[models.InvitationCategory])(nil)
