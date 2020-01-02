package bodyweight

import (
	"context"
	"refit_backend/internal/infrastructures/mysql"
	"refit_backend/internal/logger"
	"refit_backend/models"
)

// IBodyWeight interface
type IBodyWeight interface {
	Create(ctx context.Context, m *models.BodyWeight) (bodyweightID uint, err error)
	FindOne(ctx context.Context)
	FindAll(ctx context.Context)
	Update(ctx context.Context)
	Delete(ctx context.Context)
	Count(ctx context.Context)
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

func (u bodyweight) FindOne(ctx context.Context) {}

func (u bodyweight) FindAll(ctx context.Context) {}

func (u bodyweight) Update(ctx context.Context) {}

func (u bodyweight) Delete(ctx context.Context) {}

func (u bodyweight) Count(ctx context.Context) {}
