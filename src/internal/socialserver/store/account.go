package store

import (
	"context"
	v1 "go-socialapp/internal/socialserver/model/v1"
)

// ChainStore defines the account storage interface.
type AccountStore interface {
	CreateBatch(ctx context.Context, accounts []v1.Account) error
}
