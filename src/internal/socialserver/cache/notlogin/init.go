package notlogin

import (
	"github.com/marmotedu/iam/pkg/log"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var TmpWaClientCache *tmpWaClientCache

func InitTmpWaClientCache() {
	TmpWaClientCache = &tmpWaClientCache{
		tmpWaClients: make(map[string]*whatsapp.WaClient, 5),
		lock:         &sync.RWMutex{},
	}
	files, err := ioutil.ReadDir(whatsapp.PathSessionsTmp)
	if err != nil {
		log.Fatalf("get files in dir %s err %s", whatsapp.PathSessionsTmp, err.Error())
	}

	for _, file := range files {
		path := filepath.Join(whatsapp.PathSessionsTmp, file.Name())
		err = os.Remove(path)
		if err != nil {
			log.Fatalf("delete file %s err %s", path, err.Error())
		}
	}
	go TmpWaClientCache.scan()
}
