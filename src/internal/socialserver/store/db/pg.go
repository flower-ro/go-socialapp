package db

import (
	"fmt"
	"github.com/marmotedu/iam/pkg/log"
	logger "go-socialapp/internal/pkg/dblogger"
	genericoptions "go-socialapp/internal/pkg/options"
	"go-socialapp/internal/socialserver/store"
	transcationalDB "go-socialapp/pkg/db"
	"gorm.io/gorm"
	"sync"
)

type datastore struct {
	db transcationalDB.TransactionHelper
	// can include two database instance if needed
	// docker *grom.DB
	// db *gorm.DB
}

func (ds *datastore) Accounts() store.AccountStore {
	return newAccountStore(ds.db)
}

func (ds *datastore) GetTxGenerate() transcationalDB.TxGenerate {
	return ds.db.GetTxGenerate()
}

func (ds *datastore) Close() error {
	return ds.db.Close()
}

var (
	pgFactory store.Factory
	once      sync.Once
)

func GetPgFactoryOr(opts *genericoptions.DBOptions) (store.Factory, error) {
	if opts == nil && pgFactory == nil {
		return nil, fmt.Errorf("failed to get pg store fatory")
	}
	if pgFactory != nil {
		return pgFactory, nil
	}
	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		options := &transcationalDB.Options{
			Host:                  opts.Host,
			Port:                  opts.Port,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			LogLevel:              opts.LogLevel,
			Logger:                logger.New(opts.LogLevel),
		}
		dbIns, err = transcationalDB.NewPG(options)
		// uncomment the following line if you need auto migration the given models
		// not suggested in production environment.
		// migrateDatabase(dbIns)
		if err != nil {
			log.Fatal(err.Error())
		}
		transcationalDB.InitTransactionHelper(dbIns)
		pgFactory = &datastore{transcationalDB.GetTransactionHelper()}
	})

	if pgFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get pg store fatory, pgFactory: %+v, error: %w", pgFactory, err)
	}

	return pgFactory, nil
}
