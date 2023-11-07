package notlogin

import (
	"context"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go.mau.fi/whatsmeow"
)

func (t tmpWaClientCache) GetQrCodeByNewWaClient() (<-chan whatsmeow.QRChannelItem, error) {
	if len(t.tmpWaClients) > 50 {
		return nil, errors.Errorf("tmpWaClientCache size is %d,allow max is 20", len(t.tmpWaClients))
	}
	client, err := whatsapp.NewClientWithNoDevice()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	ch, err := client.WaCli.GetQRChannel(context.Background())
	if err != nil {
		if errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
			_ = client.WaCli.Connect() // just connect to websocket
			if client.WaCli.IsLoggedIn() {
				return nil, errors.WithCode(code.ErrAlreadyLoggedIn, err.Error())
			}
			return nil, errors.WithCode(code.ErrSessionSaved, err.Error())
		} else {
			return nil, errors.WithCode(code.ErrQrChannel, err.Error())
		}
	}
	err = client.WaCli.Connect()
	if err != nil {
		client.WaCli.Disconnect()
		return nil, errors.WithCode(code.ErrReconnect, err.Error())
	}
	err = t.put(client)
	if err != nil {
		client.WaCli.Disconnect()
		t.Del(client)
		return nil, errors.Wrap(err, "")
	}
	return ch, nil

}
