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

type ICardService interface {
	GetAllCards(params queryparams.ListParams) (*queryparams.PaginatedResult, error)
	GetCardByID(id uint) (*models.Card, error)
	CreateCardWithRelations(ctx context.Context, card *models.Card) error
	UpdateCardWithRelations(ctx context.Context, card *models.Card) error
	DeleteCardWithRelations(ctx context.Context, id uint) error
	GetCardCount() (int64, error)
	IsSlugAvailable(slug string, excludeID uint) (bool, error)
}

type CardService struct {
	repo repositories.ICardRepository
}

func NewCardService() ICardService {
	return &CardService{repo: repositories.NewCardRepository()}
}

func (s *CardService) GetAllCards(params queryparams.ListParams) (*queryparams.PaginatedResult, error) {
	cards, totalCount, err := s.repo.GetAllCards(params)
	if err != nil {
		logconfig.Log.Error("Kartlar alınamadı", zap.Error(err))
		return nil, errors.New("kartlar getirilirken bir hata oluştu")
	}
	result := &queryparams.PaginatedResult{
		Data: cards,
		Meta: queryparams.PaginationMeta{
			CurrentPage: params.Page,
			PerPage:     params.PerPage,
			TotalItems:  totalCount,
			TotalPages:  queryparams.CalculateTotalPages(totalCount, params.PerPage),
		},
	}
	return result, nil
}

func (s *CardService) GetCardByID(id uint) (*models.Card, error) {
	card, err := s.repo.GetCardByID(id)
	if err != nil {
		logconfig.Log.Warn("Kart bulunamadı", zap.Uint("card_id", id), zap.Error(err))
		return nil, errors.New("kart bulunamadı")
	}
	return card, nil
}

func (s *CardService) CreateCardWithRelations(ctx context.Context, card *models.Card) error {
	return s.repo.CreateCardWithRelations(ctx, card)
}

func (s *CardService) UpdateCardWithRelations(ctx context.Context, card *models.Card) error {
	return s.repo.UpdateCardWithRelations(ctx, card)
}

func (s *CardService) DeleteCardWithRelations(ctx context.Context, id uint) error {
	return s.repo.DeleteCardWithRelations(ctx, id)
}

func (s *CardService) GetCardCount() (int64, error) {
	return s.repo.GetCardCount()
}

func (s *CardService) IsSlugAvailable(slug string, excludeID uint) (bool, error) {
	return s.repo.IsSlugAvailable(slug, excludeID)
}

var _ ICardService = (*CardService)(nil)