package usecase

import (
	serviceRep "github.com/FackOff25/TechnoparkDBHW/src/internal/service/repository"
	"github.com/FackOff25/TechnoparkDBHW/src/models"
)

type UseCaseInterface interface {
	ClearData() error
	SelectStatus() (*models.ServiceStatus, error)
}

type useCase struct {
	serviceRepository serviceRep.RepositoryInterface
}

func New(serviceRepository serviceRep.RepositoryInterface) UseCaseInterface {
	return &useCase{
		serviceRepository: serviceRepository,
	}
}

func (uc *useCase) ClearData() error {
	err := uc.serviceRepository.ClearData()
	return err
}

func (uc *useCase) SelectStatus() (*models.ServiceStatus, error) {
	status, err := uc.serviceRepository.SelectStatus()
	if err != nil {
		return nil, err
	}

	return status, nil
}
