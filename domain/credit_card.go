package domain

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type CreditCard struct {
	Id              string
	Name            string
	Number          string
	ExpirationMonth int32
	ExpirationYear  int32
	Cvv             int32
	Balance         float64
	Limit           float64
	CreatedAt       time.Time
}

func NewCreditCard() *CreditCard {
	c := &CreditCard{}
	c.Id = uuid.NewV4().String()
	c.CreatedAt = time.Now()
	return c
}
