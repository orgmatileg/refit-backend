package bodyweight

import (
	"context"
	"fmt"
	"refit_backend/internal/infrastructures/mysql"
	"refit_backend/internal/logger"
	"refit_backend/models"
)

// IBodyWeight interface
type IBodyWeight interface {
	Create(ctx context.Context, m *models.BodyWeight) (bodyweightID uint, err error)
	FindOneByID(ctx context.Context, bodyWeightID string) (m *models.BodyWeight, err error)
	FindAll(ctx context.Context, limit, offset, order, userID string) (lm []*models.BodyWeight, err error)
	UpdateByID(ctx context.Context, m *models.BodyWeight, bodyWeightID string) (rowUpdated int64, err error)
	DeleteByID(ctx context.Context, bodyWeightID string) (rowDeleted int64, err error)
	Count(ctx context.Context, userID string) (count uint, err error)
}

type bodyweight struct{}

// New Repository todos
func New() IBodyWeight {
	return &bodyweight{}
}

func (u bodyweight) Create(ctx context.Context, m *models.BodyWeight) (bodyweightID uint, err error) {

	q := `
		INSERT INTO body_weight
		(weight, image, date, created_at, user_id)
		VALUES (?, ?, ?, ?, ?)
	`

	res, err := mysql.GetDB().ExecContext(ctx, q,
		m.Weight,
		m.Image,
		m.Date,
		m.CreatedAt,
		m.UserID,
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

	return uint(lastInsertedID), nil
}

func (u bodyweight) FindOneByID(ctx context.Context, bodyWeightID string) (m *models.BodyWeight, err error) {
	q := `
		SELECT id, weight, image, date, user_id, created_at
		FROM body_weight
		WHERE id = ?
	`
	err = mysql.GetDB().QueryRowContext(ctx, q, bodyWeightID).Scan(
		&m.ID,
		&m.Weight,
		&m.Image,
		&m.Date,
		&m.UserID,
		&m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (u bodyweight) FindAll(ctx context.Context, limit, offset, order, userID string) (lm []*models.BodyWeight, err error) {
	q := fmt.Sprintf(`
		SELECT id, weight, image, date, user_id, created_at
		FROM body_weight
		WHERE user_id = ?
		ORDER BY created_at %s LIMIT ? OFFSET ?
`, order)

	row, err := mysql.GetDB().QueryContext(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := row.Close(); e != nil {
			logger.Errorf("could not close row: %s", e.Error())
		}
	}()

	for row.Next() {
		var tm models.BodyWeight
		err = row.Scan(
			&tm.ID,
			&tm.Weight,
			&tm.Image,
			&tm.Date,
			&tm.UserID,
			&tm.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		lm = append(lm, &tm)
	}

	return lm, nil
}

func (u bodyweight) UpdateByID(ctx context.Context, m *models.BodyWeight, bodyWeightID string) (rowUpdated int64, err error) {
	q := `
	UPDATE body_weight
	SET weight=?, date=?, image=?
	WHERE id=? 
`
	res, err := mysql.GetDB().ExecContext(ctx, q,
		m.Weight,
		m.Date,
		m.Image,
		bodyWeightID,
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

func (u bodyweight) DeleteByID(ctx context.Context, bodyWeightID string) (rowDeleted int64, err error) {
	q := `
		DELETE FROM body_weight
		WHERE id = ?
	`
	res, err := mysql.GetDB().ExecContext(ctx, q, bodyWeightID)
	if err != nil {
		return -1, err
	}
	rowDeleted, err = res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return rowDeleted, nil
}

func (u bodyweight) Count(ctx context.Context, userID string) (count uint, err error) {

	q := `
		SELECT COUNT(*) 
		FROM body_weight 
		WHERE user_id = ?
	`
	err = mysql.GetDB().QueryRowContext(ctx, q, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
