package notlogin

import (
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"sync"
)

type tmpWaClientCache struct {
	tmpWaClients map[string]*whatsapp.WaClient
	lock         *sync.RWMutex
}

func (t tmpWaClientCache) put(client *whatsapp.WaClient) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if _, ok := t.tmpWaClients[client.FileName]; ok {
		return errors.Errorf("existed same name %s when create tmp waClient", client.FileName)
	}
	t.tmpWaClients[client.FileName] = client
	return nil
}

func (t tmpWaClientCache) Del(client *whatsapp.WaClient) {
	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.tmpWaClients, client.FileName)
}
