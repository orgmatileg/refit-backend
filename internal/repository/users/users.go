package users

import (
	"context"
	"fmt"
	"refit_backend/internal/infrastructures/mysql"
	"refit_backend/internal/logger"
	"refit_backend/models"
	"time"

	"github.com/zoonman/gravatar"
)

type IUsers interface {
	Create(ctx context.Context, mu *models.User) (userID uint, err error)
	FindOneByEmail(ctx context.Context, email string) (*models.User, error)
	FindAll(ctx context.Context)
	Update(ctx context.Context)
	Delete(ctx context.Context)
	Count(ctx context.Context)
}

type users struct{}

// New Repository Users
func New() IUsers {
	return &users{}
}

// Create repository
func (u users) Create(ctx context.Context, m *models.User) (userID uint, err error) {

	q := `
		INSERT INTO user
		(full_name, email, password, gender, country_code, role_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	db := mysql.GetDB()

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			logger.Errorf("got panic when transaction: %s", p)
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			err = tx.Rollback()
			if err != nil {
				logger.Errorf("could not rollback transaction query: %s", err.Error())
			}
		} else {
			// all good, commit
			err = tx.Commit()
			if err != nil {
				logger.Errorf("could not commit transaction query: %s", err.Error())
			}
		}
	}()

	res, err := tx.ExecContext(ctx, q,
		m.FullName,
		m.Email,
		m.Password,
		m.Gender,
		m.CountryCode,
		m.RoleID,
		m.CreatedAt,
		m.UpdatedAt,
	)

	if err != nil {
		logger.Infof("could not exec query: %s", err.Error())
		return 0, err
	}

	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		logger.Infof("could not get result query last insert id: %s", err.Error())
		return 0, err
	}

	q = fmt.Sprintf(`
			INSERT INTO user_image
			(image, user_id, created_at)
			VALUES (%s, %d, %s)
		`, gravatar.Avatar(m.Email, 256), lastInsertedID, time.Now())

	_, err = tx.ExecContext(ctx, q)
	if err != nil {
		logger.Infof("could not exec query: %s", err.Error())
		return 0, err
	}

	return uint(lastInsertedID), nil
}

// FindOneByEmail repository
func (u users) FindOneByEmail(ctx context.Context, email string) (*models.User, error) {
	q := `
		SELECT id, full_name, email, password, gender, country_code, role_id, created_at, updated_at
		FROM user 
		WHERE email = ?
	`
	var mu models.User
	db := mysql.GetDB()
	err := db.QueryRowContext(ctx, q, email).Scan(
		&mu.ID,
		&mu.FullName,
		&mu.Email,
		&mu.Password,
		&mu.Gender,
		&mu.CountryCode,
		&mu.RoleID,
		&mu.CreatedAt,
		&mu.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &mu, nil
}
func (u users) FindAll(ctx context.Context) {}
func (u users) Update(ctx context.Context)  {}
func (u users) Delete(ctx context.Context)  {}
func (u users) Count(ctx context.Context)   {}
