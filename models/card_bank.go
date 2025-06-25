package models

// CardBank represents the many-to-many relationship between Card and Bank with additional IBAN field
type CardBank struct {
	BaseModel
	CardID uint   `gorm:"index;not null"`
	BankID uint   `gorm:"index;not null"`
	IBAN   string `gorm:"size:50;not null"`
	Card   Card   `gorm:"foreignKey:CardID"`
	Bank   Bank   `gorm:"foreignKey:BankID"`
}

// TableName returns the table name for the CardBank model
func (CardBank) TableName() string {
	return "card_banks"
}
