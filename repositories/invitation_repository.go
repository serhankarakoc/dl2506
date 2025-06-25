package repositories

import (
	"context"

	"davet.link/configs/databaseconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"

	"gorm.io/gorm"
)

type IInvitationRepository interface {
	GetAllInvitations(params queryparams.ListParams) ([]models.Invitation, int64, error)
	GetInvitationByID(id uint) (*models.Invitation, error)
	CreateInvitation(ctx context.Context, invitation *models.Invitation) error
	UpdateInvitation(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error
	DeleteInvitation(ctx context.Context, id uint) error
	GetInvitationCount() (int64, error)
	GetParticipantsByInvitationID(invitationID uint) ([]models.InvitationParticipant, error)
	UpdateParticipant(id uint, participant *models.InvitationParticipant) error
	DeleteParticipant(id uint) error
}

type InvitationRepository struct {
	base IBaseRepository[models.Invitation]
	db   *gorm.DB
}

func NewInvitationRepository() IInvitationRepository {
	base := NewBaseRepository[models.Invitation](databaseconfig.GetDB())
	base.SetAllowedSortColumns([]string{"id", "invitation_key", "user_id", "category_id", "created_at"})
	base.SetPreloads(
		"User",
		"Category",
		"InvitationDetail",
		"Participants",
	)
	return &InvitationRepository{base: base, db: databaseconfig.GetDB()}
}

func (r *InvitationRepository) GetAllInvitations(params queryparams.ListParams) ([]models.Invitation, int64, error) {
	return r.base.GetAll(params)
}

func (r *InvitationRepository) GetInvitationByID(id uint) (*models.Invitation, error) {
	return r.base.GetByID(id)
}

func (r *InvitationRepository) CreateInvitation(ctx context.Context, invitation *models.Invitation) error {
	return r.base.Create(ctx, invitation)
}

func (r *InvitationRepository) UpdateInvitation(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error {
	return r.base.Update(ctx, id, data, updatedBy)
}

func (r *InvitationRepository) DeleteInvitation(ctx context.Context, id uint) error {
	return r.base.Delete(ctx, id)
}

func (r *InvitationRepository) GetInvitationCount() (int64, error) {
	return r.base.GetCount()
}

func (r *InvitationRepository) GetParticipantsByInvitationID(invitationID uint) ([]models.InvitationParticipant, error) {
	var participants []models.InvitationParticipant
	err := r.db.Where("invitation_id = ?", invitationID).Find(&participants).Error
	return participants, err
}

func (r *InvitationRepository) UpdateParticipant(id uint, participant *models.InvitationParticipant) error {
	return r.db.Model(&models.InvitationParticipant{}).Where("id = ?", id).Updates(participant).Error
}

func (r *InvitationRepository) DeleteParticipant(id uint) error {
	return r.db.Delete(&models.InvitationParticipant{}, id).Error
}

var _ IInvitationRepository = (*InvitationRepository)(nil)
var _ IBaseRepository[models.Invitation] = (*BaseRepository[models.Invitation])(nil)
