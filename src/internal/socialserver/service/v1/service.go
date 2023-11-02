package v1

import (
	whatsappClient "go-socialapp/internal/socialserver/client/whatsapp"
	"go-socialapp/internal/socialserver/service/v1/account"
	"go-socialapp/internal/socialserver/store"
)

// Service defines functions used to return resource interface.
type Service interface {
	Accounts() account.AccountSrv
}

type service struct {
	store    store.Factory
	waClient whatsappClient.Factory
}

var srv Service

// NewService returns Service interface.
func NewService(store store.Factory, waClient whatsappClient.Factory) Service {
	return &service{
		store:    store,
		waClient: waClient,
	}
}

func GetService() Service {
	if srv != nil {
		return srv
	}
	srv = NewService(store.Store(), whatsappClient.Client())
	return srv
}

func (s *service) Accounts() account.AccountSrv {
	return account.GetAccount(s.store, s.waClient)
}