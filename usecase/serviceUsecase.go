package usecase

import (
	"parkDB/models"
	"parkDB/repository"
)

type ServiceUsecaseInterface interface {
	ClearDB() (err error)
	GetStatusOfDB() (status *models.Status, err error)
}

type ServiceUsecase struct {
	serviceRepository repository.ServiceRepositoryInterface
}

func NewServiceUsecase(serviceRepository repository.ServiceRepositoryInterface) ServiceUsecaseInterface {
	return &ServiceUsecase{serviceRepository: serviceRepository}
}

func (serviceUseCase *ServiceUsecase) ClearDB() (err error) {
	return serviceUseCase.serviceRepository.ClearDB()
}

func (serviceUseCase *ServiceUsecase) GetStatusOfDB() (status *models.Status, err error) {
	return serviceUseCase.serviceRepository.GetStatusOfDB()
}
