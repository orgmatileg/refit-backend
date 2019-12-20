package todos

import "context"

type ITodos interface {
	Create(ctx context.Context)
	FindOne(ctx context.Context)
	FindAll(ctx context.Context)
	Update(ctx context.Context)
	Delete(ctx context.Context)
	Count(ctx context.Context)
}

type todos struct{}

// New Repository todos
func New() ITodos {
	return &todos{}
}

func (u todos) Create(ctx context.Context)  {}
func (u todos) FindOne(ctx context.Context) {}
func (u todos) FindAll(ctx context.Context) {}
func (u todos) Update(ctx context.Context)  {}
func (u todos) Delete(ctx context.Context)  {}
func (u todos) Count(ctx context.Context)   {}
