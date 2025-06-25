package models

type InvitationParticipant struct {
	BaseModel
	
	// Zorunlu Alanlar
	Title        string `gorm:"type:varchar(255);not null"`
	PhoneNumber  string `gorm:"type:varchar(20);not null"`
	GuestCount   int    `gorm:"not null;default:1"`
	InvitationID uint   `gorm:"index;not null"`
	
	// İlişki Tanımı
	Invitation Invitation `gorm:"foreignKey:InvitationID"`
}

func (InvitationParticipant) TableName() string {
	return "invitation_participants"
}