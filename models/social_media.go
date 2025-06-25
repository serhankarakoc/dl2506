package models

type SocialMedia struct {
	BaseModel
	IsActive bool   `gorm:"default:true;index"`
	Icon     string `gorm:"size:50;not null"`  // Font Awesome icon class name
	Name     string `gorm:"size:255;not null;index"`
}

// TableName returns the table name for the SocialMedia model
func (SocialMedia) TableName() string {
	return "social_media"
}
