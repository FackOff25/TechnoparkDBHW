package usecase

import (
	userRep "github.com/FackOff25/TechnoparkDBHW/src/internal/user/repository"
	"github.com/FackOff25/TechnoparkDBHW/src/models"
)

type UseCaseInterface interface {
	CreateUser(user *models.User) ([]*models.User, error)
	SelectUser(nickname string) (*models.User, error)
	UpdateUser(user *models.User) error
}

type useCase struct {
	userRepository userRep.RepositoryInterface
}

func New(userRepository userRep.RepositoryInterface) UseCaseInterface {
	return &useCase{
		userRepository: userRepository,
	}
}

func (uc *useCase) CreateUser(user *models.User) ([]*models.User, error) {
	existUsers, err := uc.userRepository.SelectUsersByNickNameOrEmail(user.NickName, user.Email)
	if err != nil {
		return nil, err
	} else if len(existUsers) > 0 {
		return existUsers, models.ErrConflict
	}

	err = uc.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (uc *useCase) SelectUser(nickname string) (*models.User, error) {
	user, err := uc.userRepository.SelectUserByNickName(nickname)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *useCase) UpdateUser(user *models.User) error {
	selectedUser, err := uc.userRepository.SelectUserByNickName(user.NickName)
	if err != nil {
		return err
	}

	if user.FullName == "" && user.Email == "" && user.About == "" {
		user.About = selectedUser.About
		user.Email = selectedUser.Email
		user.FullName = selectedUser.FullName
		return nil
	}

	_, err = uc.userRepository.SelectUserByEmail(user.Email)
	if err != models.ErrNotFound && err != nil {
		return err
	} else if err == nil {
		return models.ErrConflict
	}

	err = uc.userRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}
