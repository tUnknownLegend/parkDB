package usecase

import (
	conf "parkDB/config"
	"parkDB/models"
	"parkDB/repository"
)

type UserUsecaseInterface interface {
	CreateUser(user *models.User) (users *models.Users, err error)
	GetUser(nickname string) (user *models.User, err error)
	UpdateUser(user *models.User) (err error)
}

type UserUsecase struct {
	userRepository repository.UserRepositoryInterface
}

func NewUserUsecase(userRepository repository.UserRepositoryInterface) UserUsecaseInterface {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (userUseCase *UserUsecase) CreateUser(user *models.User) (users *models.Users, err error) {
	usersSlice, err := userUseCase.userRepository.GetMatchedUsers(user)
	if err != nil {
		err = conf.ConflictError
		return
	} else if len(*usersSlice) > 0 {
		users = new(models.Users)
		*users = *usersSlice
		err = conf.ConflictError
		return
	}

	err = userUseCase.userRepository.CreateUser(user)
	return
}

func (userUseCase *UserUsecase) GetUser(nickname string) (user *models.User, err error) {
	user, err = userUseCase.userRepository.GetByNickname(nickname)
	if err != nil {
		err = conf.NotFoundError
		return
	}
	return
}

func (userUseCase *UserUsecase) UpdateUser(user *models.User) (err error) {
	oldUser, err := userUseCase.userRepository.GetByNickname(user.Nickname)
	if oldUser.Nickname == "" || err != nil {
		err = conf.NotFoundError
		return
	}

	err = userUseCase.userRepository.UpdateUser(user)
	if err != nil {
		err = conf.ConflictError
		return
	}
	return
}
