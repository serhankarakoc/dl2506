package models

// CardSocialMedia represents the many-to-many relationship between Card and SocialMedia with additional URL field
type CardSocialMedia struct {
	BaseModel
	CardID        uint        `gorm:"index;not null"`
	SocialMediaID uint        `gorm:"index;not null"`
	URL           string      `gorm:"size:255;not null"`
	Card          Card        `gorm:"foreignKey:CardID"`
	SocialMedia   SocialMedia `gorm:"foreignKey:SocialMediaID"`
}

// TableName returns the table name for the CardSocialMedia model
func (CardSocialMedia) TableName() string {
	return "card_social_media"
}
