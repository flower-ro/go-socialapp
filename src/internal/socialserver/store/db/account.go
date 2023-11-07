package db

import (
	"context"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	v1 "go-socialapp/internal/socialserver/model/v1"
	transcationalDB "go-socialapp/pkg/db"
	"gorm.io/gorm"
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

func (a *accountStore) GetAllAccount(ctx context.Context) ([]v1.Account, error) {
	accounts := []v1.Account{}
	err := a.GetTxDBOr(ctx).Where("id >0").Find(&accounts).Error
	if err != nil {
		return accounts, err
	}

	return accounts, nil
}

func (a *accountStore) GetByPhone(ctx context.Context, phone string) (*v1.Account, error) {
	account := &v1.Account{}
	err := a.GetTxDBOr(ctx).Model(v1.Account{}).Where("phone_number = ?", phone).Limit(1).Find(account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//return nil, errors.WithCode(code.ErrRecordNotExisted, err.Error())
			return nil, nil
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return account, nil
}
func (a *accountStore) UpdateDevice(ctx context.Context, phone string, device string) error {
	update := map[string]interface{}{"device": device}
	return a.GetTxDBOr(ctx).Model(&v1.Account{}).Where("phone_number = ?", phone).
		Select("device").Updates(update).Error
}
func (a *accountStore) Create(ctx context.Context, account v1.Account) error {
	return a.GetTxDBOr(ctx).Create(&account).Error
}
