package db

import (
	"context"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"
)

const DB = "DB"

type TransactionHelper interface {
	TxGenerate
	GetTxDB(ctx context.Context) *gorm.DB
	GetDB() *gorm.DB
	Close() error
	GetTxGenerate() TxGenerate
}

type transactionTool struct {
	*txGenerateTool
}

var transactionHelper TransactionHelper

func InitTransactionHelper(db *gorm.DB) {
	transactionHelper = newTransactionTool(db)
}

func GetTransactionHelper() TransactionHelper {
	return transactionHelper
}

func newTransactionTool(db *gorm.DB) *transactionTool {
	return &transactionTool{
		txGenerateTool: newTxGenerateTool(db),
	}
}

func (t *transactionTool) GetDB() *gorm.DB {
	return t.db
}

func (t *transactionTool) GetTxDB(ctx context.Context) *gorm.DB {
	if ctx.Value(DB) == nil {
		return nil
	}
	return ctx.Value(DB).(*gorm.DB)
}

func (t *transactionTool) Close() error {
	db, err := t.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

func (t *transactionTool) GetTxGenerate() TxGenerate {
	return t.txGenerateTool
}
