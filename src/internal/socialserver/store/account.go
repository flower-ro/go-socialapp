package store

import (
	"context"
	v1 "go-socialapp/internal/socialserver/model/v1"
)

// ChainStore defines the account storage interface.
type AccountStore interface {
	CreateBatch(ctx context.Context, accounts []v1.Account) error
	GetAllAccount(ctx context.Context) ([]v1.Account, error)
	GetByPhone(ctx context.Context, phone string) (*v1.Account, error)
	UpdateDevice(ctx context.Context, phone string, device string) error
	Create(ctx context.Context, account v1.Account) error
	DelByPhone(ctx context.Context, phone string) error
}
