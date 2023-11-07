package loggedin

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/marmotedu/errors"
	whatsappBase "go-socialapp/internal/pkg/third-party/whatsapp"
	whatsappClient "go-socialapp/internal/socialserver/client/whatsapp"
	"sync"
	"time"
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
	if old, ok := w.cache[phone]; ok {
		old.GetClient().Disconnect()
	}
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
	spew.Dump("-------------------1==err====", newClient.WaCli.Connect())
	var i int
	for {

		if i > 3 {
			break
		}

		time.Sleep(30 * time.Second)

		spew.Dump("-------------------1==connted====", newClient.WaCli.IsConnected())
		spew.Dump("-------------------1==IsLoggedIn====", newClient.WaCli.IsLoggedIn())
		i++

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
