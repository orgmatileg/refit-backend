package bodyweight

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"refit_backend/internal/constants"
	"refit_backend/internal/helpers"
	"refit_backend/internal/infrastructures/s3"
	"refit_backend/internal/logger"
	"refit_backend/internal/repository"
	"refit_backend/models"
	"regexp"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/minio/minio-go"
)

var (
	regexNumberOnly = regexp.MustCompile("^[0-9]*$")
)

// IBodyWeight interface
type IBodyWeight interface {
	Create(ctx context.Context, weight, date, userID string, fh *multipart.FileHeader) (bodyweightID uint, err error)
	FindOneByID(ctx context.Context, bodyWeightID string) (m *models.BodyWeight, err error)
	FindAll(ctx context.Context, limit, offset, order, userID string) (lm []*models.BodyWeight, count uint, err error)
	UpdateByID(ctx context.Context, rm *models.BodyWeight, bodyweightID string) (err error)
	DeleteByID(ctx context.Context, bodyWeightID string) (err error)
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
	unixTime := time.Now().Unix()
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
			fmt.Sprintf("refit/users/%s/bodyweights/%d.%s", userID, unixTime, helpers.GetExtensionFile(ft)),
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
	m.Image = fmt.Sprintf("https://static.luqmanul.com/refit/users/%s/bodyweights/%d.%s", userID, unixTime, helpers.GetExtensionFile(fh.Header.Get("Content-Type")))
	m.CreatedAt = time.Now()

	bodyweightID, err = u.repository.BodyWeight().Create(ctx, &m)
	if err != nil {
		logger.Infof("could not create bodyweight repository: %s", err.Error())
		return 0, err
	}

	return bodyweightID, nil
}

func (u bodyweight) FindOneByID(ctx context.Context, bodyWeightID string) (m *models.BodyWeight, err error) {
	err = validation.Validate(bodyWeightID, validation.Match(regexNumberOnly))
	if err != nil {
		logger.Infof("could not validate: %s", err.Error())
		return nil, errors.New("invalid bodyweight_id param, should be number only")
	}
	m, err = u.repository.BodyWeight().FindOneByID(ctx, bodyWeightID)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			logger.Infof("could not find bodyweight by id: %s", err.Error())
			return nil, errors.New("bodyweight_id not exists")
		default:
			logger.Infof("could not find bodyweight by id: %s", err.Error())
			return nil, err
		}
	}

	return m, nil
}

func (u bodyweight) FindAll(ctx context.Context, limit, offset, order, userID string) (lm []*models.BodyWeight, count uint, err error) {
	err = helpers.ValidationQueryParamFindAll(limit, offset, order)
	if err != nil {
		logger.Infof("could not validate: %s", err.Error())
		return nil, 0, err
	}

	lm, err = u.repository.BodyWeight().FindAll(ctx, limit, offset, order, userID)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			logger.Infof("could not find all BodyWeight: %s", err.Error())
			return nil, 0, errors.New("no row exists")
		default:
			logger.Infof("could not find all BodyWeight: %s", err.Error())
			return nil, 0, err
		}
	}

	count, err = u.repository.BodyWeight().Count(ctx, userID)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			logger.Infof("could not find all body_weight: %s", err.Error())
			return nil, 0, errors.New("no row exists")
		default:
			logger.Infof("could not find all body_weight: %s", err.Error())
			return nil, 0, err
		}
	}

	return lm, count, nil
}

func (u bodyweight) UpdateByID(ctx context.Context, rm *models.BodyWeight, bodyweightID string) (err error) {

	mb, err := u.repository.BodyWeight().FindOneByID(ctx, bodyweightID)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			logger.Infof("could not find bodyweight by id: %s", err.Error())
			return errors.New("bodyweight_id not exists")
		default:
			logger.Infof("could not find bodyweight by id: %s", err.Error())
			return err
		}
	}

	rm.UserID = mb.UserID
	rm.CreatedAt = mb.CreatedAt

	_, err = u.repository.BodyWeight().UpdateByID(ctx, rm, bodyweightID)
	if err != nil {
		logger.Infof("could not update bodyweight by id: %s", err.Error())
		return err
	}
	return nil
}

func (u bodyweight) DeleteByID(ctx context.Context, bodyweightID string) (err error) {

	err = validation.Validate(bodyweightID, validation.Match(regexNumberOnly))
	if err != nil {
		logger.Infof("could not validate: %s", err.Error())
		return errors.New("invalid bodyweight_id param, should be number only")
	}
	_, err = u.repository.BodyWeight().DeleteByID(ctx, bodyweightID)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			logger.Infof("could not delete bodyweight by id: %s", err.Error())
			return errors.New("bodyweight_id not exists")
		default:
			logger.Infof("could not delete bodyweight by id: %s", err.Error())
			return err
		}
	}
	return nil
}
