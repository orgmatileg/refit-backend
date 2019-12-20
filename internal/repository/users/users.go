package users

import "context"

type IUsers interface {
	Create(ctx context.Context)
	FindOne(ctx context.Context)
	FindAll(ctx context.Context)
	Update(ctx context.Context)
	Delete(ctx context.Context)
	Count(ctx context.Context)
}

type users struct{}

// New Repository Users
func New() IUsers {
	return &users{}
}

func (u users) Create(ctx context.Context)  {}
func (u users) FindOne(ctx context.Context) {}
func (u users) FindAll(ctx context.Context) {}
func (u users) Update(ctx context.Context)  {}
func (u users) Delete(ctx context.Context)  {}
func (u users) Count(ctx context.Context)   {}
