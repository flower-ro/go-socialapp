package cache

import (
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/socialserver/client/whatsapp"
	"sync"
)

var WaClientCache *waClientCache

type waClientCache struct {
	cache map[string]whatsapp.Factory
	lock  *sync.RWMutex
}

func InitWaClientCache() {
	WaClientCache = newWaClient()
}

func newWaClient() *waClientCache {
	return &waClientCache{
		cache: make(map[string]whatsapp.Factory, 10),
		lock:  &sync.RWMutex{},
	}
}

func (w *waClientCache) Put(phone string, factory whatsapp.Factory) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if _, ok := w.cache[phone]; ok {
		return errors.WithCode(code.WaClientExistedInCache, "")
	}
	w.cache[phone] = factory
	return nil
}

func (w *waClientCache) Get(phone string) whatsapp.Factory {
	w.lock.RLock()
	defer w.lock.RUnlock()
	client, _ := w.cache[phone]
	return client
}

func (w *waClientCache) Del(phone string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	delete(w.cache, phone)
}

func (w *waClientCache) Size() int {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return len(w.cache)
}
