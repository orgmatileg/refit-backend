package bodyweight

import (
	"context"
	"fmt"
	"refit_backend/internal/constants"
	"refit_backend/internal/infrastructures/s3"
	"refit_backend/internal/logger"
	"refit_backend/internal/repository"
	"refit_backend/models"
	"strings"
	"time"

	"github.com/minio/minio-go"
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

	s3Client := s3.GetS3Client()

	s := "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"

	if err != nil {
		logger.Infof("could not decode base64 to byte: %s", err.Error())
		return 0, err
	}
	r := strings.NewReader(s)

	n, err := s3Client.PutObjectWithContext(
		ctx,
		"static-luqmanul",
		fmt.Sprintf("%d.png", time.Now().Unix()),
		r,
		r.Size(),
		minio.PutObjectOptions{
			ContentType: "image/png",
		},
	)
	if err != nil {
		logger.Infof("could not put object to spaces: %s", err.Error())
		return 0, err
	}

	fmt.Println(n)
	fmt.Println(n)
	fmt.Println(n)
	fmt.Println(n)
	fmt.Println(n)
	fmt.Println(n)

	return bodyweightID, nil
}

func (u bodyweight) FindOne(ctx context.Context) {}
func (u bodyweight) FindAll(ctx context.Context) {}
func (u bodyweight) Update(ctx context.Context)  {}
func (u bodyweight) Delete(ctx context.Context)  {}
func (u bodyweight) Count(ctx context.Context)   {}
