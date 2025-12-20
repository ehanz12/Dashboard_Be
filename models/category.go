package models

import "time"

type Category struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	UserID    uint   `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	//relasi untuk references
	User User `json:"-" gorm:"foreginKey:userID"`
	Expenses []Expense `json:"expenses,omitempty"`
}