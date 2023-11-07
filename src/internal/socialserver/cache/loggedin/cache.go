package loggedin

import (
	"github.com/marmotedu/errors"
	whatsappBase "go-socialapp/internal/pkg/third-party/whatsapp"
	whatsappApi "go-socialapp/internal/socialserver/client/whatsapp"
	"sync"
)

var WaApiCache *waApiCache

type waApiCache struct {
	cache map[string]whatsappApi.Factory
	lock  *sync.RWMutex
}

func InitWaApiCache() {
	WaApiCache = newWaApi()
	//WaClientCache.initSessionFiles(whatsappBase.PathSessions)
}

func newWaApi() *waApiCache {
	return &waApiCache{
		cache: make(map[string]whatsappApi.Factory, 10),
		lock:  &sync.RWMutex{},
	}
}

func (w *waApiCache) Put(phone string, factory whatsappApi.Factory) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if old, ok := w.cache[phone]; ok {
		old.GetClient().Disconnect()
	}
	w.cache[phone] = factory
	return nil
}

func (w *waApiCache) Get(phone string) (whatsappApi.Factory, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()
	waApi, ok := w.cache[phone]
	if ok {
		waApi.UpdateLastOperationTime()

		err := whatsappBase.WaitLogin(waApi.GetClient())
		if err != nil {
			return waApi, errors.Wrap(err, " ")
		}

		return waApi, nil
	}
	newClient, err := whatsappBase.NewWaClientWithDevice(phone)
	err = whatsappBase.WaitLogin(newClient.WaCli)
	if err != nil {
		return waApi, errors.Wrap(err, " ")
	}
	waApi = whatsappApi.NewFactory(newClient.WaCli, newClient.Db)
	w.cache[phone] = waApi
	return waApi, nil
}

//
//// 可以做优化，达到最大数量，或者闲置时间太长都可以 删除该缓存
//func (w *waClientCache) Del(phone string) {
//	w.lock.Lock()
//	defer w.lock.Unlock()
//	delete(w.cache, phone)
//}
//
//func (w *waClientCache) Size() int {
//	w.lock.RLock()
//	defer w.lock.RUnlock()
//	return len(w.cache)
//}
