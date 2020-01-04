package bodyweight

import (
	"context"
	"fmt"
	"mime/multipart"
	"refit_backend/internal/constants"
	"refit_backend/internal/helpers"
	"refit_backend/internal/infrastructures/s3"
	"refit_backend/internal/logger"
	"refit_backend/internal/repository"
	"refit_backend/models"
	"strconv"
	"time"

	"github.com/minio/minio-go"
)

// IBodyWeight interface
type IBodyWeight interface {
	Create(ctx context.Context, weight, date, userID string, fh *multipart.FileHeader) (bodyweightID uint, err error)
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

func (u bodyweight) Create(ctx context.Context, weight, date, userID string, fh *multipart.FileHeader) (bodyweightID uint, err error) {

	m := models.BodyWeight{}

	timeDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		logger.Infof("could not parse time: %s", err.Error())
		return 0, err
	}
	m.Date = timeDate

	weightInt, err := strconv.ParseFloat(weight, 64)
	if err != nil {
		logger.Infof("could not parse weight string to float64: %s", err.Error())
		return 0, err
	}
	m.Weight = float64(weightInt)

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		logger.Infof("could not parse weight string to float64: %s", err.Error())
		return 0, err
	}
	m.UserID = uint(userIDInt)

	if fh == nil {
		m.Image = constants.ImageDefault
	} else {

		f, err := fh.Open()
		if err != nil {
			logger.Infof("could not open file header %s", err.Error())
			return 0, err
		}
		defer f.Close()

		ft := fh.Header.Get("Content-Type")

		_, err = s3.GetS3Client().PutObjectWithContext(
			ctx,
			"static-luqmanul",
			fmt.Sprintf("refit/users/%s/bodyweights/%d.%s", userID, time.Now().Unix(), helpers.GetExtensionFile(ft)),
			f,
			fh.Size,
			minio.PutObjectOptions{
				ContentType: ft,
			},
		)
		if err != nil {
			logger.Infof("could not put object to spaces: %s", err.Error())
			return 0, err
		}

	}
	m.CreatedAt = time.Now()

	bodyweightID, err = u.repository.BodyWeight().Create(ctx, &m)
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
