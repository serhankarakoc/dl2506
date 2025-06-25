package services

import (
	"context"
	"davet.link/configs/logconfig"
	"davet.link/models"
	"davet.link/pkg/filemanager"
	"davet.link/pkg/queryparams"
	"davet.link/pkg/random"
	"davet.link/repositories"
	"errors"
	"go.uber.org/zap"
)

type IInvitationService interface {
	GetAllInvitations(params queryparams.ListParams) (*queryparams.PaginatedResult, error)
	GetInvitationByID(id uint) (*models.Invitation, error)
	GetInvitationByKey(key string) (*models.Invitation, error)
	CreateInvitation(ctx context.Context, invitation *models.Invitation) error
	UpdateInvitation(ctx context.Context, invitation *models.Invitation) error
	DeleteInvitation(ctx context.Context, id uint) error
}

type InvitationService struct {
	repo repositories.IInvitationRepository
}

func NewInvitationService() IInvitationService {
	return &InvitationService{repo: repositories.NewInvitationRepository()}
}

func (s *InvitationService) GetAllInvitations(params queryparams.ListParams) (*queryparams.PaginatedResult, error) {
	invitations, totalCount, err := s.repo.GetAllInvitations(params)
	if err != nil {
		logconfig.Log.Error("Davetiyeler alınamadı", zap.Error(err))
		return nil, errors.New("davetiyeler getirilirken bir hata oluştu")
	}
	result := &queryparams.PaginatedResult{
		Data: invitations,
		Meta: queryparams.PaginationMeta{
			CurrentPage: params.Page,
			PerPage:     params.PerPage,
			TotalItems:  totalCount,
			TotalPages:  queryparams.CalculateTotalPages(totalCount, params.PerPage),
		},
	}
	return result, nil
}

func (s *InvitationService) GetInvitationByID(id uint) (*models.Invitation, error) {
	invitation, err := s.repo.GetInvitationByID(id)
	if err != nil {
		logconfig.Log.Warn("Davetiye bulunamadı", zap.Uint("invitation_id", id), zap.Error(err))
		return nil, errors.New("davetiye bulunamadı")
	}
	return invitation, nil
}

func (s *InvitationService) GetInvitationByKey(key string) (*models.Invitation, error) {
	invitation, err := s.repo.GetInvitationByKey(key)
	if err != nil {
		logconfig.Log.Warn("Davetiye anahtarı ile bulunamadı", zap.String("invitation_key", key), zap.Error(err))
		return nil, errors.New("davetiye bulunamadı")
	}
	return invitation, nil
}

func (s *InvitationService) CreateInvitation(ctx context.Context, invitation *models.Invitation) error {
	uniqueKey, err := random.GenerateUniqueString(12)
	if err != nil {
		logconfig.Log.Error("Benzersiz davetiye anahtarı oluşturulamadı", zap.Error(err))
		return errors.New("davetiye anahtarı oluşturulamadı")
	}
	invitation.InvitationKey = uniqueKey

	return s.repo.CreateWithRelations(ctx, invitation)
}

func (s *InvitationService) UpdateInvitation(ctx context.Context, invitation *models.Invitation) error {
	return s.repo.UpdateWithRelations(ctx, invitation)
}

func (s *InvitationService) DeleteInvitation(ctx context.Context, id uint) error {
	invitation, err := s.GetInvitationByID(id)
	if err != nil {
		return errors.New("silinecek davetiye bulunamadı")
	}
	
	if err := s.repo.DeleteWithRelations(ctx, id); err != nil {
		return err
	}
	
	if invitation.Image != "" {
		filemanager.DeleteFile("invitations", invitation.Image)
	}

	return nil
}

var _ IInvitationService = (*InvitationService)(nil)