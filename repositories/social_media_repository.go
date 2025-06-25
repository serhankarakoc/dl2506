package repositories

import (
	"context"

	"davet.link/configs/databaseconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"

	"gorm.io/gorm"
)

type ISocialMediaRepository interface {
	GetAllSocialMedias(params queryparams.ListParams) ([]models.SocialMedia, int64, error)
	GetSocialMediaByID(id uint) (*models.SocialMedia, error)
	CreateSocialMedia(ctx context.Context, socialMedia *models.SocialMedia) error
	BulkCreateSocialMedias(ctx context.Context, socialMedias []models.SocialMedia) error
	UpdateSocialMedia(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error
	BulkUpdateSocialMedias(ctx context.Context, condition map[string]interface{}, data map[string]interface{}, updatedBy uint) error
	DeleteSocialMedia(ctx context.Context, id uint) error
	BulkDeleteSocialMedias(ctx context.Context, condition map[string]interface{}) error
	GetSocialMediaCount() (int64, error)
}

type SocialMediaRepository struct {
	base IBaseRepository[models.SocialMedia]
	db   *gorm.DB
}

func NewSocialMediaRepository() ISocialMediaRepository {
	base := NewBaseRepository[models.SocialMedia](databaseconfig.GetDB())
	base.SetAllowedSortColumns([]string{"id", "name", "is_active", "created_at"})
	return &SocialMediaRepository{base: base, db: databaseconfig.GetDB()}
}

func (r *SocialMediaRepository) GetAllSocialMedias(params queryparams.ListParams) ([]models.SocialMedia, int64, error) {
	return r.base.GetAll(params)
}

func (r *SocialMediaRepository) GetSocialMediaByID(id uint) (*models.SocialMedia, error) {
	return r.base.GetByID(id)
}

func (r *SocialMediaRepository) CreateSocialMedia(ctx context.Context, socialMedia *models.SocialMedia) error {
	return r.base.Create(ctx, socialMedia)
}

func (r *SocialMediaRepository) BulkCreateSocialMedias(ctx context.Context, socialMedias []models.SocialMedia) error {
	return r.base.BulkCreate(ctx, socialMedias)
}

func (r *SocialMediaRepository) UpdateSocialMedia(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error {
	return r.base.Update(ctx, id, data, updatedBy)
}

func (r *SocialMediaRepository) BulkUpdateSocialMedias(ctx context.Context, condition map[string]interface{}, data map[string]interface{}, updatedBy uint) error {
	return r.base.BulkUpdate(ctx, condition, data, updatedBy)
}

func (r *SocialMediaRepository) DeleteSocialMedia(ctx context.Context, id uint) error {
	return r.base.Delete(ctx, id)
}

func (r *SocialMediaRepository) BulkDeleteSocialMedias(ctx context.Context, condition map[string]interface{}) error {
	return r.base.BulkDelete(ctx, condition)
}

func (r *SocialMediaRepository) GetSocialMediaCount() (int64, error) {
	return r.base.GetCount()
}

var _ ISocialMediaRepository = (*SocialMediaRepository)(nil)
var _ IBaseRepository[models.SocialMedia] = (*BaseRepository[models.SocialMedia])(nil)
