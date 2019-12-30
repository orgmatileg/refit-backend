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

// ValidateUpdate model user for update
// func (u BodyWeight) ValidateUpdate() error {
// 	return validation.ValidateStruct(&u,
// 		validation.Field(&u.FullName, validation.Required, validation.Length(5, 50)),
// 		validation.Field(&u.Email, validation.Required, is.Email),
// 		validation.Field(&u.Gender, validation.Required, validation.In("male", "female", "others", "wont_tell")),
// 		validation.Field(&u.RoleID, validation.In(1, 2, 3)),
// 	)
// }
