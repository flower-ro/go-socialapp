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
		old.GetClient().WaCli.Disconnect()
	}
	w.cache[phone] = factory
	return nil
}

func (w *waApiCache) Get(phone string) (whatsappApi.Factory, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	waApi, ok := w.cache[phone]
	if ok {
		waApi.UpdateLastOperationTime()
		return waApi, nil
	}
	newClient, err := whatsappBase.NewWaClientWithDevice(phone)
	if err != nil {
		return waApi, errors.Wrap(err, " ")
	}
	waApi = whatsappApi.NewFactory(newClient)
	w.cache[phone] = waApi
	return waApi, nil
}

// // 可以做优化，达到最大数量，或者闲置时间太长都可以 删除该缓存
func (w *waApiCache) Del(phone string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	waApi, ok := w.cache[phone]
	if ok {
		waApi.GetClient().WaCli.Disconnect()
		delete(w.cache, phone)
	}

}

//
//func (w *waClientCache) Size() int {
//	w.lock.RLock()
//	defer w.lock.RUnlock()
//	return len(w.cache)
//}
