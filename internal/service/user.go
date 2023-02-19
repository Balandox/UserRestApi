package service

import (
	"my-rest-api/internal/model"
	"my-rest-api/internal/repository"
	"my-rest-api/pkg/store"
)

type UserService interface {
	GetAllUsers() ([]*model.User, error)

	GetUserById(id int) (*model.User, error)

	AddUser(user *model.User) (*model.User, error)

	UpdateUser(user *model.User) error

	DeleteUser(id int) error

	FindByEmail(email string) (*model.User, error)
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(store *store.Store) UserService {
	return &UserServiceImpl{
		userRepository: repository.NewUserRepository(store),
	}
}

func (u *UserServiceImpl) GetAllUsers() ([]*model.User, error) {
	return u.userRepository.FindAll()
}

func (u *UserServiceImpl) GetUserById(id int) (*model.User, error) {
	return u.userRepository.FindById(id)
}

func (u *UserServiceImpl) AddUser(user *model.User) (*model.User, error) {
	userAfterInsert, err := u.userRepository.Create(user)
	if err != nil {
		return nil, err
	}
	user.Id = userAfterInsert.Id

	return user, nil
}

func (u *UserServiceImpl) UpdateUser(user *model.User) error {
	return u.userRepository.Update(user)
}

func (u *UserServiceImpl) DeleteUser(id int) error {
	return u.userRepository.Delete(id)
}

func (u *UserServiceImpl) FindByEmail(email string) (*model.User, error) {
	return u.userRepository.FindByEmail(email)
}
