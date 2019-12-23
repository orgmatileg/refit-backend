package models

import (
	"time"
)

// User Model Struct
// gender: 'male', 'female', 'others', 'wont_tell'
// role_id: '1:administrator', '2:normal user', '3:coach'
type User struct {
	ID          uint      `json:"id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Gender      string    `json:"gender"`
	CountryCode uint      `json:"country_code"`
	RoleID      uint      `json:"role_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
