package notlogin

import (
	"github.com/marmotedu/iam/pkg/log"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"os"
	"sync"
)

var TmpWaClientCache *tmpWaClientCache

func InitTmpWaClientCache() {
	TmpWaClientCache = &tmpWaClientCache{
		tmpWaClients: make(map[string]*whatsapp.WaClient, 5),
		lock:         &sync.RWMutex{},
	}
	// 删除指定的文件夹及其所有子文件夹和文件
	err := os.RemoveAll(whatsapp.PathSessionsTmp)
	if err != nil {
		log.Fatalf("delete dir %s err %s", whatsapp.PathSessionsTmp, err.Error())
	}
	err = os.Mkdir(whatsapp.PathSessionsTmp, 0755) // 创建 1 级目录
	if err != nil {
		log.Fatalf("create dir %s err %s", whatsapp.PathSessionsTmp, err.Error())
	}
	go TmpWaClientCache.scan()
}
