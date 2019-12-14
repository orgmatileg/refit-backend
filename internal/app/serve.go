package app

import (
	"context"
)

// Serve interface
type Serve interface {
	GetCtx() context.Context
	GetMySQL() float64
}
