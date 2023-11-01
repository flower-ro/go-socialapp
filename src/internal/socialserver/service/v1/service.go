package v1

import (
	"go-socialapp/internal/socialserver/service/v1/account"
	"go-socialapp/internal/socialserver/store"
)

// Service defines functions used to return resource interface.
type Service interface {
	Accounts() account.AccountSrv
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
	return NewService(store.Client())
}

func (s *service) Accounts() account.AccountSrv {
	return account.GetAccount(s.store)
}
