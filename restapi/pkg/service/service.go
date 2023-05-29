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
	GetUser(userid uint16) (domian.ApiResponse, error)
	GetPageQty(limit int) (int, error)
	PostUser(nickname, name, lastname, password string) (domian.ApiResponse, error)
	DeleteUser(userid uint16) error
	PutUser(upUser domian.UpdateUser) (domian.ApiResponse, error)
}

// interface which collect verification methods
type ServiceValidation interface {
	Authentication(nickname, password string, userid uint16) (bool, error)
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
func (s *service) GetUser(userid uint16) (domian.ApiResponse, error) {
	return s.store.GetUser(userid)
}
func (s *service) GetPageQty(limit int) (int, error) {
	return s.store.GetPageQty(limit)
}
func (s *service) PostUser(nickname, name, lastname, password string) (domian.ApiResponse, error) {
	return s.store.PostUser(nickname, name, lastname, password)
}
func (s *service) DeleteUser(userid uint16) error {
	return s.store.DeleteUser(userid)
}
func (s *service) PutUser(upUser domian.UpdateUser) (domian.ApiResponse, error) {
	return s.store.PutUser(upUser)
}
func (s *service) Authentication(nickname, password string, userid uint16) (bool, error) {
	return s.store.Authentication(nickname, password, userid)
}
func (s *service) PasswordHashing(password string) (string, error) {
	return s.store.PasswordHashing(password)
}
