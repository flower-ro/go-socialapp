package v1

import (
	"go-socialapp/internal/socialserver/service/v1/account"
	"go-socialapp/internal/socialserver/service/v1/group"
	"go-socialapp/internal/socialserver/store"
)

// Service defines functions used to return resource interface.
type Service interface {
	Accounts() account.AccountSrv
	Groups() group.GroupSrv
}

type service struct {
	store store.Factory
}

var srv Service

// NewService returns Service interface.
func NewService(store store.Factory) Service {
	return &service{
		store: store,
	}
}

func GetService() Service {
	if srv != nil {
		return srv
	}
	srv = NewService(store.Store())
	return srv
}

func (s *service) Accounts() account.AccountSrv {
	return account.GetAccount(s.store)
}

func (s *service) Groups() group.GroupSrv {
	return group.GetGroup(s.store)
}
