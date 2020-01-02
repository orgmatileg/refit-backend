package bodyweight

import (
	"context"
	"refit_backend/internal/constants"
	"refit_backend/internal/logger"
	"refit_backend/internal/repository"
	"refit_backend/models"
	"time"
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

type bodyweight struct {
	repository repository.IRepository
}

// New Repository todos
func New() IBodyWeight {
	return &bodyweight{
		repository: repository.New(),
	}
}

func (u bodyweight) Create(ctx context.Context, m *models.BodyWeight) (bodyweightID uint, err error) {

	err = m.ValidateCreate()
	if err != nil {
		logger.Infof("could not validate: %s", err.Error())
		return 0, err
	}

	if m.Image == "" {
		m.Image = constants.ImageDefault
	}

	m.CreatedAt = time.Now()

	bodyweightID, err = u.repository.BodyWeight().Create(ctx, m)
	if err != nil {
		logger.Infof("could not create bodyweight repository: %s", err.Error())
		return 0, err
	}

	return bodyweightID, nil
}

func (u bodyweight) FindOne(ctx context.Context) {}
func (u bodyweight) FindAll(ctx context.Context) {}
func (u bodyweight) Update(ctx context.Context)  {}
func (u bodyweight) Delete(ctx context.Context)  {}
func (u bodyweight) Count(ctx context.Context)   {}
