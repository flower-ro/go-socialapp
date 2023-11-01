package account

import (
	"github.com/marmotedu/errors"
	v1 "go-socialapp/internal/socialserver/model/v1"
	transcationalDB "go-socialapp/pkg/db"

	"context"
	"go-socialapp/internal/socialserver/store"
)

type AccountSrv interface {
	CreateBatch(ctx context.Context, accounts []v1.Account) error
}

type accountService struct {
	transcationalDB.TxGenerate
	store store.Factory
}

var _ AccountSrv = (*accountService)(nil)

var accountSrv *accountService

func GetAccount(store store.Factory) *accountService {
	if accountSrv != nil {
		return accountSrv
	}
	accountSrv = &accountService{
		TxGenerate: store.GetTxGenerate(),
		store:      store}
	return accountSrv
}

func (a *accountService) CreateBatch(ctx context.Context, accounts []v1.Account) error {
	return a.StartTransaction(ctx, func(ctx context.Context) error {
		err := a.store.Accounts().CreateBatch(ctx, accounts)
		if err != nil {
			return errors.Wrapf(err, "CreateBatch error")
		}
		return nil
	})
}
