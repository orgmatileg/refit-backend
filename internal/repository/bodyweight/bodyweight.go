package bodyweight

import "context"

type IBodyWeight interface {
	Create(ctx context.Context)
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

func (u bodyweight) Create(ctx context.Context)  {}
func (u bodyweight) FindOne(ctx context.Context) {}
func (u bodyweight) FindAll(ctx context.Context) {}
func (u bodyweight) Update(ctx context.Context)  {}
func (u bodyweight) Delete(ctx context.Context)  {}
func (u bodyweight) Count(ctx context.Context)   {}
