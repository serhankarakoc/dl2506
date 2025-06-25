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

type ISocialMediaService interface {
	GetAllSocialMedias(params queryparams.ListParams) (*queryparams.PaginatedResult, error)
	GetSocialMediaByID(id uint) (*models.SocialMedia, error)
	CreateSocialMedia(ctx context.Context, socialMedia *models.SocialMedia) error
	UpdateSocialMedia(ctx context.Context, id uint, socialMediaData *models.SocialMedia, updatedBy uint) error
	DeleteSocialMedia(ctx context.Context, id uint) error
	GetSocialMediaCount() (int64, error)
}

type SocialMediaService struct {
	repo repositories.ISocialMediaRepository
}

func NewSocialMediaService() ISocialMediaService {
	return &SocialMediaService{repo: repositories.NewSocialMediaRepository()}
}

func (s *SocialMediaService) GetAllSocialMedias(params queryparams.ListParams) (*queryparams.PaginatedResult, error) {
	items, totalCount, err := s.repo.GetAllSocialMedias(params)
	if err != nil {
		logconfig.Log.Error("Sosyal medya kayıtları alınamadı", zap.Error(err))
		return nil, errors.New("sosyal medya kayıtları getirilirken bir hata oluştu")
	}
	result := &queryparams.PaginatedResult{
		Data: items,
		Meta: queryparams.PaginationMeta{
			CurrentPage: params.Page,
			PerPage:     params.PerPage,
			TotalItems:  totalCount,
			TotalPages:  queryparams.CalculateTotalPages(totalCount, params.PerPage),
		},
	}
	return result, nil
}

func (s *SocialMediaService) GetSocialMediaByID(id uint) (*models.SocialMedia, error) {
	item, err := s.repo.GetSocialMediaByID(id)
	if err != nil {
		logconfig.Log.Warn("Sosyal medya kaydı bulunamadı", zap.Uint("social_media_id", id), zap.Error(err))
		return nil, errors.New("sosyal medya kaydı bulunamadı")
	}
	return item, nil
}

func (s *SocialMediaService) CreateSocialMedia(ctx context.Context, socialMedia *models.SocialMedia) error {
	return s.repo.CreateSocialMedia(ctx, socialMedia)
}

func (s *SocialMediaService) UpdateSocialMedia(ctx context.Context, id uint, socialMediaData *models.SocialMedia, updatedBy uint) error {
	_, err := s.repo.GetSocialMediaByID(id)
	if err != nil {
		return errors.New("sosyal medya kaydı bulunamadı")
	}
	updateData := map[string]interface{}{
		"name":      socialMediaData.Name,
		"icon":      socialMediaData.Icon,
		"is_active": socialMediaData.IsActive,
	}
	return s.repo.UpdateSocialMedia(ctx, id, updateData, updatedBy)
}

func (s *SocialMediaService) DeleteSocialMedia(ctx context.Context, id uint) error {
	return s.repo.DeleteSocialMedia(ctx, id)
}

func (s *SocialMediaService) GetSocialMediaCount() (int64, error) {
	return s.repo.GetSocialMediaCount()
}

var _ ISocialMediaService = (*SocialMediaService)(nil)
