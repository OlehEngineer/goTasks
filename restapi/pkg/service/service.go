package service

import (
	"github.com/OlehEngineer/goTasks/goTasks/restapi/pkg/domian"
	"github.com/OlehEngineer/goTasks/goTasks/restapi/pkg/usecases"
)

type Service interface {
	ServiceStore
	ServiceValidation
}

// interface which collect user processing methods
type ServiceStore interface {
	GetUsersPage(page, limit int) ([]domian.ApiResponse, error)
	GetUser(userId uint16) (domian.ApiResponse, error)
	GetPageQty(limit int) (int, error)
	CreateUser(newUser domian.UserSignUp, password string) (domian.ApiResponse, error)
	DeleteUser(userId uint16) error
	UpdateUser(upUser domian.UpdateUser) (domian.ApiResponse, error)
}

// interface which collect verification methods
type ServiceValidation interface {
	Authentication(nickname, password string, userId uint16) (bool, error)
	PasswordHashing(password string) (string, error)
}
type service struct {
	store usecases.Store
}

func New(store usecases.Store) *service {
	return &service{
		store: store,
	}
}

// implement methods as for usecases.Store from domian package
func (s *service) GetUsersPage(page, limit int) ([]domian.ApiResponse, error) {
	return s.store.GetUsersPage(page, limit)
}
func (s *service) GetUser(userId uint16) (domian.ApiResponse, error) {
	return s.store.GetUser(userId)
}
func (s *service) GetPageQty(limit int) (int, error) {
	return s.store.GetPageQty(limit)
}
func (s *service) CreateUser(newUser domian.UserSignUp, password string) (domian.ApiResponse, error) {
	return s.store.CreateUser(newUser, password)
}
func (s *service) DeleteUser(userId uint16) error {
	return s.store.DeleteUser(userId)
}
func (s *service) UpdateUser(upUser domian.UpdateUser) (domian.ApiResponse, error) {
	return s.store.UpdateUser(upUser)
}
