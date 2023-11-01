package db

import (
	"context"
	"gorm.io/gorm"
)

type TxGenerate interface {
	StartTransaction(ctx context.Context, opts ...Option) error
}

type txGenerateTool struct {
	db *gorm.DB
}

func newTxGenerateTool(db *gorm.DB) *txGenerateTool {
	return &txGenerateTool{
		db: db,
	}
}

type Option func(ctx context.Context) error

func (t *txGenerateTool) StartTransaction(ctx context.Context, opts ...Option) error {

	return t.db.Transaction(func(tx *gorm.DB) error {
		nctx := t.setTxContext(ctx, tx)
		for _, o := range opts {
			err := o(nctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (t *txGenerateTool) setTxContext(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, DB, tx)
}
