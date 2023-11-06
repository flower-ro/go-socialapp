package notlogin

import (
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"sync"
)

var TmpWaClientCache *tmpWaClientCache

func init() {
	TmpWaClientCache = &tmpWaClientCache{
		tmpWaClients: make(map[string]*whatsapp.WaClient, 5),
		lock:         &sync.RWMutex{},
	}
	go TmpWaClientCache.scan()
}
