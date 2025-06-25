package services

import (
	"context"
	"errors"

	"davet.link/configs/logconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"
	"davet.link/repositories"

	"go.uber.org/zap"
)

type IInvitationCategoryService interface {
	GetAllCategories(params queryparams.ListParams) (*queryparams.PaginatedResult, error)
	GetCategoryByID(id uint) (*models.InvitationCategory, error)
	CreateCategory(ctx context.Context, category *models.InvitationCategory) error
	UpdateCategory(ctx context.Context, id uint, categoryData *models.InvitationCategory, updatedBy uint) error
	DeleteCategory(ctx context.Context, id uint) error
	GetCategoryCount() (int64, error)
}

type InvitationCategoryService struct {
	repo repositories.IInvitationCategoryRepository
}

func NewInvitationCategoryService() IInvitationCategoryService {
	return &InvitationCategoryService{repo: repositories.NewInvitationCategoryRepository()}
}

func (s *InvitationCategoryService) GetAllCategories(params queryparams.ListParams) (*queryparams.PaginatedResult, error) {
	categories, totalCount, err := s.repo.GetAllCategories(params)
	if err != nil {
		logconfig.Log.Error("Davet kategorileri alınamadı", zap.Error(err))
		return nil, errors.New("davet kategorileri getirilirken bir hata oluştu")
	}
	result := &queryparams.PaginatedResult{
		Data: categories,
		Meta: queryparams.PaginationMeta{
			CurrentPage: params.Page,
			PerPage:     params.PerPage,
			TotalItems:  totalCount,
			TotalPages:  queryparams.CalculateTotalPages(totalCount, params.PerPage),
		},
	}
	return result, nil
}

func (s *InvitationCategoryService) GetCategoryByID(id uint) (*models.InvitationCategory, error) {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		logconfig.Log.Warn("Davet kategorisi bulunamadı", zap.Uint("category_id", id), zap.Error(err))
		return nil, errors.New("davet kategorisi bulunamadı")
	}
	return category, nil
}

func (s *InvitationCategoryService) CreateCategory(ctx context.Context, category *models.InvitationCategory) error {
	return s.repo.CreateCategory(ctx, category)
}

func (s *InvitationCategoryService) UpdateCategory(ctx context.Context, id uint, categoryData *models.InvitationCategory, updatedBy uint) error {
	_, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return errors.New("davet kategorisi bulunamadı")
	}
	updateData := map[string]interface{}{
		"name":      categoryData.Name,
		"icon":      categoryData.Icon,
		"template":  categoryData.Template,
		"is_active": categoryData.IsActive,
	}
	return s.repo.UpdateCategory(ctx, id, updateData, updatedBy)
}

func (s *InvitationCategoryService) DeleteCategory(ctx context.Context, id uint) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *InvitationCategoryService) GetCategoryCount() (int64, error) {
	return s.repo.GetCategoryCount()
}

var _ IInvitationCategoryService = (*InvitationCategoryService)(nil)
