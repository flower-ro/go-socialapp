package loggedin

import (
	"github.com/marmotedu/errors"
	whatsappBase "go-socialapp/internal/pkg/third-party/whatsapp"
	whatsappClient "go-socialapp/internal/socialserver/client/whatsapp"
	"sync"
)

var WaClientCache *waClientCache

type waClientCache struct {
	cache map[string]whatsappClient.Factory
	lock  *sync.RWMutex
}

func InitWaClientCache() {
	WaClientCache = newWaClient()
	//WaClientCache.initSessionFiles(whatsappBase.PathSessions)
}

func newWaClient() *waClientCache {
	return &waClientCache{
		cache: make(map[string]whatsappClient.Factory, 10),
		lock:  &sync.RWMutex{},
	}
}

func (w *waClientCache) Put(phone string, factory whatsappClient.Factory) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	//if _, ok := w.cache[phone]; ok {
	//	return errors.WithCode(code.WaClientExistedInCache, "")
	//}
	w.cache[phone] = factory
	return nil
}

func (w *waClientCache) Get(phone string) (whatsappClient.Factory, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()
	client, ok := w.cache[phone]
	if ok {
		client.UpdateLastOperationTime()
		return client, nil
	}
	newClient, err := whatsappBase.NewWaClientWithDevice(phone)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	client = whatsappClient.NewFactory(newClient.WaCli, newClient.Db)
	w.cache[phone] = client
	return client, nil
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
