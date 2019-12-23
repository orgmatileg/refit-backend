package delivery

import (
	"refit_backend/internal/delivery/http"
)

// IDelivery interface
type IDelivery interface {
	HTTP() http.IServerHTTP
}

type delivery struct{}

// New Delivery
func New() IDelivery {
	return &delivery{}
}

func (s delivery) HTTP() http.IServerHTTP {
	return http.NewServerHTTP()
}
