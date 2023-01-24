package usecase

import (
	"parkDB/models"
	"parkDB/repository"
)

type ServiceUsecaseInterface interface {
	Clear() (err error)
	GetStatus() (status *models.Status, err error)
}

type ServiceUsecase struct {
	serviceRepository repository.ServiceRepositoryInterface
}

func NewServiceUsecase(serviceRepository repository.ServiceRepositoryInterface) ServiceUsecaseInterface {
	return &ServiceUsecase{serviceRepository: serviceRepository}
}

func (serviceUseCase *ServiceUsecase) Clear() (err error) {
	return serviceUseCase.serviceRepository.Clear()
}

func (serviceUseCase *ServiceUsecase) GetStatus() (status *models.Status, err error) {
	return serviceUseCase.serviceRepository.GetStatus()
}
