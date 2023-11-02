package whatsapp

import (
	"fmt"
	"github.com/marmotedu/iam/pkg/log"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"time"
)

var (
	cli           *whatsmeow.Client
	historySyncID int32
	startupTime   = time.Now().Unix()
)

func InitWaDB() *sqlstore.Container {
	// Running Whatsapp
	storeContainer, err := sqlstore.New("sqlite3",
		fmt.Sprintf("file:%s/%s?_foreign_keys=off", PathStorages, DBName), nil)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return storeContainer
}

func InitWaCLI(storeContainer *sqlstore.Container) *whatsmeow.Client {
	device, err := storeContainer.GetFirstDevice()
	if err != nil {
		log.Errorf("Failed to get device: %v", err)
		panic(err)
	}

	osName := fmt.Sprintf("%s %s", AppOs, AppVersion)
	store.DeviceProps.PlatformType = &AppPlatform
	store.DeviceProps.Os = &osName
	cli = whatsmeow.NewClient(device, nil)
	cli.EnableAutoReconnect = true
	cli.AutoTrustIdentity = true
	cli.AddEventHandler(handler)

	return cli
}
