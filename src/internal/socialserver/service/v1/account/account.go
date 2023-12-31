package account

import (
	"context"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/pkg/log"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/socialserver/cache/loggedin"
	"go-socialapp/internal/socialserver/model/network"
	v1 "go-socialapp/internal/socialserver/model/v1"
	"go-socialapp/internal/socialserver/store"
	transcationalDB "go-socialapp/pkg/db"
	"os"
	"path/filepath"
	"strings"
)

type AccountSrv interface {
	//Login(phone string) (string, error)
	Logout(phone string) error
	CreateOrUpdate(phone string, device string) error
	GetAllAccount(ctx context.Context) ([]v1.Account, error)
	IsOnWhatsApp(ctx context.Context, owner string, phones []string) (*network.IsOnWhatAppRes, error)
	DelByPhone(phone string) error
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

func (a *accountService) IsOnWhatsApp(ctx context.Context, owner string, phones []string) (*network.IsOnWhatAppRes, error) {
	waApi, err := loggedin.WaApiCache.Get(owner)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(phones) <= 0 {
		return nil, errors.New("members cannot be 0 when create group ")
	}
	err = waApi.GetClient().WaitLogin()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	result, err := waApi.General().IsOnWhatsApp(phones)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(result) <= 0 {
		return &network.IsOnWhatAppRes{
			Total: len(phones),
		}, nil
	}
	validMembers := make([]string, 0, len(result))
	for _, one := range result {
		validMembers = append(validMembers, one.Query)
	}
	isIn := &network.IsOnWhatAppRes{
		Total:   len(phones),
		Valid:   len(result),
		Members: strings.Join(validMembers, ","),
	}

	return isIn, nil
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

func (a *accountService) DelByPhone(phone string) error {
	loggedin.WaApiCache.Del(phone)
	path := filepath.Join(whatsapp.PathSessions, phone+".db")
	_, err := os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			log.Errorf("delete file %s err %s", path, err.Error())
		}
	}
	err = a.store.Accounts().DelByPhone(context.Background(), phone)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
