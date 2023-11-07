package account

import (
	"context"
	"github.com/marmotedu/errors"
	v1 "go-socialapp/internal/socialserver/model/v1"
	"go-socialapp/internal/socialserver/store"
	transcationalDB "go-socialapp/pkg/db"
)

type AccountSrv interface {
	CreateBatch(ctx context.Context, accounts []v1.Account) error
	//Login(phone string) (string, error)
	Logout(phone string) error
	CreateOrUpdate(phone string, device string) error
	GetAllAccount(ctx context.Context) ([]v1.Account, error)
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

func (a *accountService) Logout(phone string) error {
	return nil
}

func (a *accountService) CreateOrUpdate(phone string, device string) error {
	return a.StartTransaction(context.Background(), func(ctx context.Context) error {
		//查询
		account, err := a.store.Accounts().GetByPhone(ctx, phone)
		if err != nil {
			return errors.Wrap(err, "")
		}
		if account != nil {
			return a.store.Accounts().UpdateDevice(ctx, phone, device)
		}
		return a.store.Accounts().Create(ctx, v1.Account{
			PhoneNumber: phone,
			Device:      device,
		})
	})
}

func (a *accountService) GetAllAccount(ctx context.Context) ([]v1.Account, error) {
	return a.store.Accounts().GetAllAccount(ctx)
}
