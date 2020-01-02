package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v3"
)

// BodyWeight Model struct
type BodyWeight struct {
	ID        uint      `json:"id"`
	Weight    float32   `json:"weight"`
	Image     string    `json:"image"`
	Date      time.Time `json:"gender"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ValidateCreate model user for create
func (u BodyWeight) ValidateCreate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Weight, validation.Required, validation.Min(1), validation.Max(500)),
	)
}

// // ValidateUpdate model user for update
// func (u BodyWeight) ValidateUpdate() error {
// 	return validation.ValidateStruct(&u,
// 		validation.Field(&u.Weight, validation.Required, validation.Min(1), validation.Max(200)),
// 		validation.Field(&u.Image, validation.Required, is.Email),
// 	)
// }
