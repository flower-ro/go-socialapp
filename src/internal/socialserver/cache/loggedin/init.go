package loggedin

//
//import (
//	"github.com/marmotedu/iam/pkg/log"
//	"go-socialapp/internal/pkg/third-party/whatsapp"
//	whatsappClient "go-socialapp/internal/socialserver/client/whatsapp"
//	"io/ioutil"
//	"path/filepath"
//	"strings"
//)
//
//func (w *loggedin.waClientCache) initSessionFiles(sessionPath string) {
//	dir, err := ioutil.ReadDir(sessionPath)
//	var files []string
//	if err != nil {
//		log.Fatalf("read fail for session path %s", sessionPath)
//	}
//	for _, fi := range dir {
//		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".db") { // 目录, 递归遍历
//			files = append(files, fi.Name())
//		}
//	}
//	if len(files) <= 0 {
//		return
//	}
//	for _, file := range files {
//		names := strings.Split(file, ".")
//		sessionFile := filepath.Join(sessionPath, file)
//		db, err := whatsapp.InitWaDB(sessionFile)
//		if err != nil {
//			log.Fatalf("%+v", err)
//		}
//		cli, err := whatsapp.InitWaCLI(db)
//		if err != nil {
//			log.Fatalf("%+v", err)
//		}
//		w.cache[names[0]] = whatsappClient.NewFactory(cli, db)
//	}
//
//	log.Infof("初始化 wa session 数量 %d", len(w.cache))
//}
