package store

import transcationalDB "go-socialapp/pkg/db"

var store Factory

// Factory defines the tg task server platform storage interface.
type Factory interface {
	Accounts() AccountStore

	GetTxGenerate() transcationalDB.TxGenerate
	Close() error
}

// Client return the store client instance.
func Store() Factory {
	return store
}

// SetClient set the iam store client.
func SetClient(factory Factory) {
	store = factory
}
