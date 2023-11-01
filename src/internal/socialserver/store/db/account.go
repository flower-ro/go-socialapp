package db

import (
	"context"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	v1 "go-socialapp/internal/socialserver/model/v1"
	transcationalDB "go-socialapp/pkg/db"
)

type accountStore struct {
	transcationalDB.TransactionHelper
}

func newAccountStore(db transcationalDB.TransactionHelper) *accountStore {
	return &accountStore{
		TransactionHelper: db,
	}
}

func (a *accountStore) CreateBatch(ctx context.Context, accounts []v1.Account) error {
	err := a.GetTxDB(ctx).CreateInBatches(accounts, 100).Error
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil

}
