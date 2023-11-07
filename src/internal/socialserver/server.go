package socialserver

import (
	"github.com/marmotedu/iam/pkg/log"
	genericoptions "go-socialapp/internal/pkg/options"
	"go-socialapp/internal/socialserver/cache/loggedin"
	"go-socialapp/internal/socialserver/cache/notlogin"
	"go-socialapp/internal/socialserver/config"
	"go-socialapp/internal/socialserver/enter/listen"
	"go-socialapp/internal/socialserver/store"
	"go-socialapp/internal/socialserver/store/db"
	"go-socialapp/internal/socialserver/ws"
	"go-socialapp/pkg/shutdown"
	"go-socialapp/pkg/shutdown/shutdownmanagers/posixsignal"

	//"github.com/marmotedu/iam/pkg/storage"
	genericapiserver "go-socialapp/internal/pkg/server"
)

type socialServer struct {
	gs               *shutdown.GracefulShutdown
	genericAPIServer *genericapiserver.GenericAPIServer
	dbOptions        *genericoptions.DBOptions
}

func createSocialServer(cfg *config.Config) (*socialServer, error) {
	//创建 shutdown 实例
	gs := shutdown.New()
	//添加监听的信号
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New() //创建api server并初始化
	if err != nil {
		return nil, err
	}

	server := &socialServer{
		gs:               gs,
		genericAPIServer: genericServer,
		dbOptions:        cfg.DBOptions,
	}

	return server, nil
}

type preparedTaskServer struct {
	*socialServer
}

func (s *socialServer) PrepareRun() preparedTaskServer {
	s.initDB()
	initRouter(s.genericAPIServer.Engine)
	ws.InitWsClientManager()
	loggedin.InitWaApiCache()
	notlogin.InitTmpWaClientCache()
	listen.InitWaListen()

	//	s.initRedisStore()
	//设置监听到指定信号时，需要执行的回调函数。这些回调函数可以执行一些清理工作。
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		mysqlStore, _ := db.GetPgFactoryOr(nil)
		if mysqlStore != nil {
			_ = mysqlStore.Close()
		}
		s.genericAPIServer.Close()

		return nil
	}))

	return preparedTaskServer{s}
}

func (s preparedTaskServer) Run() error {
	//法启动 shutdown 实例
	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}
	return s.genericAPIServer.Run()

}

func (s *socialServer) initDB() {
	storeIns, err := db.GetPgFactoryOr(s.dbOptions)
	if err != nil {
		log.Fatalf("err=%+v", err)
	}
	// storeIns, _ := etcd.GetEtcdFactoryOr(c.etcdOptions, nil)
	store.SetClient(storeIns)
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, lastErr error) {
	genericConfig = genericapiserver.NewConfig()
	//health mode middlewares
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	//EnableProfiling EnableMetrics
	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}
