package repository

import (
	"effective-mobile/internal/dto"
	"effective-mobile/internal/model"
)

type SubscribeRepository interface {
	GetAll(filter dto.Filter) (*[]model.Subscription, error)
	GetSum(filter dto.Filter) (float64, error)
	Create(data dto.SubscriptionDto) (*model.Subscription, error)
	Update(data dto.SubscriptionDto) (*model.Subscription, error)
	GetOne(id string) (*model.Subscription, error)
	Delete(id string) error
}
