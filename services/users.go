package services

import (
	"github.com/frusdelion/zendesk-product_security_challenge/models"
	"github.com/frusdelion/zendesk-product_security_challenge/repositories"
)

type UserService interface {
	FindUserByEmail(email string) (*models.User, error)
	ActivateEmail(user *models.User) error
	UpdateUser(user *models.User, diffUpdate models.User) (*models.User, error)
}

func NewUserService(repository repositories.UserRepository) UserService {
	return &userService{repository: repository}
}

type userService struct {
	repository repositories.UserRepository
}

func (u userService) UpdateUser(user *models.User, diffUpdate models.User) (*models.User, error) {
	return u.repository.UpdateUser(user, diffUpdate)
}

func (u userService) ActivateEmail(user *models.User) error {
	_, err := u.repository.UpdateUser(user, models.User{EmailVerified: true})
	return err
}

func (u userService) FindUserByEmail(email string) (*models.User, error) {
	return u.repository.FindUserByEmail(email)
}
