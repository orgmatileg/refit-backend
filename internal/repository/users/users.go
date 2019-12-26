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

// IUsers repository interface
type IUsers interface {
	Create(ctx context.Context, mu *models.User) (userID uint, err error)
	FindOneByEmail(ctx context.Context, email string) (*models.User, error)
	FindOneByID(ctx context.Context, userID string) (*models.User, error)
	FindAll(ctx context.Context, limit, offset, order string) (lmu []*models.User, err error)
	UpdateByID(ctx context.Context, m *models.User, userID string) (rowUpdated int64, err error)
	DeleteByID(ctx context.Context, userID string) (rowDeleted int64, err error)
	Count(ctx context.Context) (count uint, err error)
}

type users struct{}

// New Repository Users
func New() IUsers {
	return &users{}
}

// FindOneByID repository users
func (u users) FindOneByID(ctx context.Context, userID string) (mu *models.User, err error) {
	q := `
		SELECT user.id, full_name, email, password, gender, country_code, country.name as country_name, role_id, role.name as role_name, user_image.image, user.created_at, updated_at 
		FROM user
		INNER JOIN user_image on user.id = user_image.user_id
		INNER JOIN role on user.role_id = role.id
		INNER JOIN country on user.country_code = country.code
		WHERE user.id = ?
	`

	mu = new(models.User)
	err = mysql.GetDB().QueryRowContext(ctx, q, userID).Scan(
		&mu.ID,
		&mu.FullName,
		&mu.Email,
		&mu.Password,
		&mu.Gender,
		&mu.CountryCode,
		&mu.CountryName,
		&mu.RoleID,
		&mu.RoleName,
		&mu.Image,
		&mu.CreatedAt,
		&mu.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return mu, nil
}

// Create repository users
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
		var errTx error
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			logger.Errorf("got panic when transaction: %s", p)
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			errTx = tx.Rollback()
			if errTx != nil {
				logger.Errorf("could not rollback transaction query: %s", err.Error())
			}
		} else {
			// all good, commit
			errTx = tx.Commit()
			if errTx != nil {
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

	q = `
			INSERT INTO user_image
			(image, user_id, created_at)
			VALUES (?, ?, ?)
	`

	_, err = tx.ExecContext(ctx, q, gravatar.Avatar(m.Email, 256), lastInsertedID, time.Now())
	if err != nil {
		logger.Infof("could not exec query: %s", err.Error())
		return 0, err
	}

	return uint(lastInsertedID), nil
}

// FindOneByEmail repository users
func (u users) FindOneByEmail(ctx context.Context, email string) (*models.User, error) {
	q := `
		SELECT id, full_name, email, password, gender, country_code, role_id, created_at, updated_at
		FROM user 
		WHERE email = ?
	`
	var mu models.User
	err := mysql.GetDB().QueryRowContext(ctx, q, email).Scan(
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

// FindAll repository users
func (u users) FindAll(ctx context.Context, limit, offset, order string) (lmu []*models.User, err error) {
	q := fmt.Sprintf(`
	SELECT user.id, full_name, email, password, gender, country_code, country.name as country_name, role_id, role.name as role_name, user_image.image, user.created_at, updated_at 
	FROM user
	INNER JOIN user_image on user.id = user_image.user_id
	INNER JOIN role on user.role_id = role.id
	INNER JOIN country on user.country_code = country.code
	ORDER BY user.created_at %s LIMIT ? OFFSET ?
`, order)

	row, err := mysql.GetDB().QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := row.Close(); e != nil {
			logger.Errorf("could not close row: %s", e.Error())
		}
	}()

	for row.Next() {
		var tmu models.User
		err = row.Scan(
			&tmu.ID,
			&tmu.FullName,
			&tmu.Email,
			&tmu.Password,
			&tmu.Gender,
			&tmu.CountryCode,
			&tmu.CountryName,
			&tmu.RoleID,
			&tmu.RoleName,
			&tmu.Image,
			&tmu.CreatedAt,
			&tmu.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		lmu = append(lmu, &tmu)
	}

	return lmu, nil
}

// Update repository users
func (u users) UpdateByID(ctx context.Context, m *models.User, userID string) (rowUpdated int64, err error) {
	q := `
		UPDATE user
		SET full_name=?, gender=?, country_code=?
		WHERE id=? 
	`
	res, err := mysql.GetDB().ExecContext(ctx, q,
		m.FullName,
		m.Gender,
		m.CountryCode,
		userID,
	)
	if err != nil {
		return -1, err
	}
	rowUpdated, err = res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return rowUpdated, nil
}

// Delete repository users
func (u users) DeleteByID(ctx context.Context, userID string) (rowDeleted int64, err error) {
	q := `
		DELETE FROM user
		WHERE id=? 
	`
	res, err := mysql.GetDB().ExecContext(ctx, q, userID)
	if err != nil {
		return -1, err
	}
	rowDeleted, err = res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return rowDeleted, nil
}

// Count repository users
func (u users) Count(ctx context.Context) (count uint, err error) {
	q := `
		SELECT COUNT(*) 
		FROM user 
	`
	err = mysql.GetDB().QueryRowContext(ctx, q).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
