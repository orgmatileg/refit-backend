package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/go-ozzo/ozzo-validation/v3/is"
)

// User Model Struct
// gender: 'male', 'female', 'others', 'wont_tell'
// role_id: '1:administrator', '2:normal user', '3:coach'
type User struct {
	ID                 uint      `json:"id"`
	FullName           string    `json:"full_name"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	PasswordValidation string    `json:"password_validation"`
	Image              string    `json:"image"`
	Gender             string    `json:"gender"`
	RoleID             uint      `json:"role_id"`
	RoleName           string    `json:"role_name"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ValidateCreate model user for create
func (u User) ValidateCreate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FullName, validation.Required, validation.Length(5, 50)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(5, 50)),
		validation.Field(&u.PasswordValidation, validation.Required, validation.In(u.Password)),
		validation.Field(&u.Gender, validation.Required, validation.In("male", "female", "others", "wont_tell")),
		validation.Field(&u.RoleID, validation.In(1, 2, 3)),
	)
}

// ValidateUpdate model user for update
func (u User) ValidateUpdate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FullName, validation.Required, validation.Length(5, 50)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Gender, validation.Required, validation.In("male", "female", "others", "wont_tell")),
		validation.Field(&u.RoleID, validation.In(1, 2, 3)),
	)
}
