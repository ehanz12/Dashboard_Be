package models

import "time"

type Expense struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`
	CategoryID *uint     `json:"category_id"`
	Amount     float64   `json:"amount"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	//relasi untuk category
	User     User     `json:"-" gorm:"foreginKey:UserID"`
	Category Category `json:"category" gorm:"foreginKey:CategoryID"`
}
