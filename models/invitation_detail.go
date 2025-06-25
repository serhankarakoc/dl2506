package models

type InvitationDetail struct {
	BaseModel
	
	// Zorunlu alan (Birebir ilişki için)
	InvitationID uint `gorm:"uniqueIndex;not null"`

	// Opsiyonel Alanlar
	Title        string `gorm:"type:varchar(255)"`
	Person       string `gorm:"type:varchar(255)"`

	// Kişinin Ebeveynleri
	IsMotherLive    bool   `gorm:"not null;default:true"`
	MotherName      string `gorm:"type:varchar(100)"`
	MotherSurname   string `gorm:"type:varchar(100)"`
	IsFatherLive    bool   `gorm:"not null;default:true"`
	FatherName      string `gorm:"type:varchar(100)"`
	FatherSurname   string `gorm:"type:varchar(100)"`

	// Gelin Detayları
	BrideName           string `gorm:"type:varchar(100)"`
	BrideSurname        string `gorm:"type:varchar(100)"`
	IsBrideMotherLive   bool   `gorm:"not null;default:true"`
	BrideMotherName     string `gorm:"type:varchar(100)"`
	BrideMotherSurname  string `gorm:"type:varchar(100)"`
	IsBrideFatherLive   bool   `gorm:"not null;default:true"`
	BrideFatherName     string `gorm:"type:varchar(100)"`
	BrideFatherSurname  string `gorm:"type:varchar(100)"`

	// Damat Detayları
	GroomName           string `gorm:"type:varchar(100)"`
	GroomSurname        string `gorm:"type:varchar(100)"`
	IsGroomMotherLive   bool   `gorm:"not null;default:true"`
	GroomMotherName     string `gorm:"type:varchar(100)"`
	GroomMotherSurname  string `gorm:"type:varchar(100)"`
	IsGroomFatherLive   bool   `gorm:"not null;default:true"`
	GroomFatherName     string `gorm:"type:varchar(100)"`
	GroomFatherSurname  string `gorm:"type:varchar(100)"`

	// İlişki Tanımı
	Invitation *Invitation `gorm:"foreignKey:InvitationID"`
}

func (InvitationDetail) TableName() string {
	return "invitation_details"
}