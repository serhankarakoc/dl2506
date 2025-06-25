package models

import "time"

type Invitation struct {
	BaseModel
	
	// --- GÜNCELLENMİŞ ZORUNLU VE INDEXLİ ALANLAR ---
	InvitationKey  string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Image         string    `gorm:"type:varchar(255);not null"`
	UserID         uint      `gorm:"index;not null"`
	CategoryID     uint      `gorm:"index;not null"`
	Template       string    `gorm:"type:varchar(100);not null"`
	Type           string    `gorm:"type:varchar(50);not null;default:'basic'"`
	IsConfirmed    bool      `gorm:"not null;default:false;index"`
	IsParticipant  bool      `gorm:"not null;default:true"`
	
	// --- Opsiyonel Alanlar (Değişiklik Yok) ---
	Title         string    `gorm:"type:varchar(255)"`
	Description   string    `gorm:"type:text"`
	Venue         string    `gorm:"type:varchar(255)"`
	Address       string    `gorm:"type:varchar(255)"`
	Location      string    `gorm:"type:varchar(255)"`
	Link          string    `gorm:"type:varchar(255)"`
	Telephone     string    `gorm:"type:varchar(20)"`
	Note          string    `gorm:"type:text"`
	Date          time.Time `gorm:"index"`
	Time          time.Time
	
	// --- Relationships (İlişki tanımları daha sonra ayarlanacak) ---
	User               *User
	Category           *InvitationCategory
	InvitationDetail   *InvitationDetail
	Participants       []InvitationParticipant
}

func (Invitation) TableName() string {
	return "invitations"
}