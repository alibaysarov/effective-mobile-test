package service

import (
	"effective-mobile/internal/dto"
	"effective-mobile/internal/model"
	"effective-mobile/internal/storage/repository"
	"log/slog"
)

type SubscribeService interface {
	GetAll(filter dto.Filter) ([]model.Subscription, error)
	GetSum(filter dto.Filter) (float64, error)
	Create(data dto.SubscriptionDto) (*model.Subscription, error)
	Update(data dto.SubscriptionDto) (*model.Subscription, error)
	GetOne(id string) (*model.Subscription, error)
	Delete(id string) error
}
type subscribeService struct {
	subscribeRepository repository.SubscribeRepository
	logger              *slog.Logger
}

func NewSubscribeService(subscribeRepository repository.SubscribeRepository, logger *slog.Logger) SubscribeService {
	return &subscribeService{
		subscribeRepository: subscribeRepository,
		logger:              logger,
	}
}

func (s *subscribeService) GetAll(filter dto.Filter) ([]model.Subscription, error) {
	result, err := s.subscribeRepository.GetAll(filter)

	if err != nil {
		return []model.Subscription{}, err
	}
	if result != nil {
		return *result, nil
	}
	return nil, nil
}

func (s *subscribeService) GetSum(filter dto.Filter) (float64, error) {
	result, err := s.subscribeRepository.GetSum(filter)
	if err != nil {
		return 0, nil
	}
	return result, nil
}

func (s *subscribeService) Create(data dto.SubscriptionDto) (*model.Subscription, error) {
	response, err := s.subscribeRepository.Create(data)
	if err != nil {
		s.logger.Error("error while creating subscription" + err.Error())
		return nil, err
	}
	s.logger.Info("Subscription created")
	return response, nil
}
func (s *subscribeService) Update(data dto.SubscriptionDto) (*model.Subscription, error) {
	response, err := s.subscribeRepository.Update(data)
	if err != nil {
		s.logger.Error("error while updating subscription" + err.Error())
		return nil, err
	}
	s.logger.Info("Subscription updated")
	return response, nil
}
func (s *subscribeService) GetOne(id string) (*model.Subscription, error) {
	firstSub, err := s.subscribeRepository.GetOne(id)
	if err != nil {
		return nil, err
	}
	return firstSub, nil
}
func (s *subscribeService) Delete(id string) error {
	err := s.subscribeRepository.Delete(id)
	if err != nil {
		s.logger.Error("error while deleting subscription: " + err.Error())
		return err
	}
	s.logger.Info("Subscription deleted")
	return nil
}
