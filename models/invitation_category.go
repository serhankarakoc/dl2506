package models

type InvitationCategory struct {
	BaseModel
	
	// Zorunlu Alanlar
	IsActive bool   `gorm:"not null;default:true;index"`
	Template string `gorm:"type:varchar(255);not null"`
	Name     string `gorm:"type:varchar(255);not null;index"`
	Icon     string `gorm:"type:varchar(50);not null"`
	
	// İlişki Tanımı
	Invitations []Invitation `gorm:"foreignKey:CategoryID"`
}

func (InvitationCategory) TableName() string {
	return "invitation_categories"
}