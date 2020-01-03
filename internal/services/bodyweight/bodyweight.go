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

	s := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAPoAAAD6CAIAAAAHjs1qAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAA2RpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuNS1jMDIxIDc5LjE1NTc3MiwgMjAxNC8wMS8xMy0xOTo0NDowMCAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wTU09Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9tbS8iIHhtbG5zOnN0UmVmPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvc1R5cGUvUmVzb3VyY2VSZWYjIiB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iIHhtcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD0ieG1wLmRpZDo0MjY4QTU4ODI2NkRFMjExQTA5QkU3NzkzRUE5Qjc5QiIgeG1wTU06RG9jdW1lbnRJRD0ieG1wLmRpZDo2MTg5MDlGMEZFRUUxMUU0QTc2QkNEMjY5OTFDODIxMSIgeG1wTU06SW5zdGFuY2VJRD0ieG1wLmlpZDo2MTg5MDlFRkZFRUUxMUU0QTc2QkNEMjY5OTFDODIxMSIgeG1wOkNyZWF0b3JUb29sPSJBZG9iZSBQaG90b3Nob3AgQ1M1IFdpbmRvd3MiPiA8eG1wTU06RGVyaXZlZEZyb20gc3RSZWY6aW5zdGFuY2VJRD0ieG1wLmlpZDpCQkY0NTk2QjEwQjRFMjExQkVBN0JFQTA3MEQ2OTg2QSIgc3RSZWY6ZG9jdW1lbnRJRD0ieG1wLmRpZDo0MjY4QTU4ODI2NkRFMjExQTA5QkU3NzkzRUE5Qjc5QiIvPiA8L3JkZjpEZXNjcmlwdGlv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"
	ss := strings.Split(s, ",")
	if err != nil {
		logger.Infof("could not decode base64 to byte: %s", err.Error())
		return 0, err
	}
	r := strings.NewReader(ss[1])

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
