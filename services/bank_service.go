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

type IBankService interface {
	GetAllBanks(params queryparams.ListParams) (*queryparams.PaginatedResult, error)
	GetBankByID(id uint) (*models.Bank, error)
	CreateBank(ctx context.Context, bank *models.Bank) error
	UpdateBank(ctx context.Context, id uint, bankData *models.Bank, updatedBy uint) error
	DeleteBank(ctx context.Context, id uint) error
	GetBankCount() (int64, error)
}

type BankService struct {
	repo repositories.IBankRepository
}

func NewBankService() IBankService {
	return &BankService{repo: repositories.NewBankRepository()}
}

func (s *BankService) GetAllBanks(params queryparams.ListParams) (*queryparams.PaginatedResult, error) {
	banks, totalCount, err := s.repo.GetAllBanks(params)
	if err != nil {
		logconfig.Log.Error("Bankalar alınamadı", zap.Error(err))
		return nil, errors.New("bankalar getirilirken bir hata oluştu")
	}
	result := &queryparams.PaginatedResult{
		Data: banks,
		Meta: queryparams.PaginationMeta{
			CurrentPage: params.Page,
			PerPage:     params.PerPage,
			TotalItems:  totalCount,
			TotalPages:  queryparams.CalculateTotalPages(totalCount, params.PerPage),
		},
	}
	return result, nil
}

func (s *BankService) GetBankByID(id uint) (*models.Bank, error) {
	bank, err := s.repo.GetBankByID(id)
	if err != nil {
		logconfig.Log.Warn("Banka bulunamadı", zap.Uint("bank_id", id), zap.Error(err))
		return nil, errors.New("banka bulunamadı")
	}
	return bank, nil
}

func (s *BankService) CreateBank(ctx context.Context, bank *models.Bank) error {
	return s.repo.CreateBank(ctx, bank)
}

func (s *BankService) UpdateBank(ctx context.Context, id uint, bankData *models.Bank, updatedBy uint) error {
	_, err := s.repo.GetBankByID(id)
	if err != nil {
		return errors.New("banka bulunamadı")
	}
	updateData := map[string]interface{}{
		"name":      bankData.Name,
		"is_active": bankData.IsActive,
	}
	return s.repo.UpdateBank(ctx, id, updateData, updatedBy)
}

func (s *BankService) DeleteBank(ctx context.Context, id uint) error {
	return s.repo.DeleteBank(ctx, id)
}

func (s *BankService) GetBankCount() (int64, error) {
	return s.repo.GetBankCount()
}

var _ IBankService = (*BankService)(nil)
