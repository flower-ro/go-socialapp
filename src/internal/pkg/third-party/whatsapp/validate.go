package whatsapp

import (
	"fmt"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	utils "go-socialapp/internal/pkg/util"
	"go.mau.fi/whatsmeow/types"
	"strings"
	"time"
)

func (w *WaClient) ValidateJidWithLogin(jid string) (types.JID, error) {
	err := w.MustLogin()
	if err != nil {
		return types.JID{}, err
	}
	return ParseJID(jid)
}

func (w *WaClient) MustLogin() error {
	if w.WaCli == nil {
		return errors.WithCode(code.ClientNotInitialized, "Whatsapp client is not initialized")
	}
	if !w.WaCli.IsConnected() {
		return errors.WithCode(code.NotConnectServer, "you are not connect to services server, please reconnect")
	} else if !w.WaCli.IsLoggedIn() {
		return errors.WithCode(code.NotLoginServer, "you are not login to services server, please login")
	}
	return nil
}

func (w *WaClient) WaitLogin() error {
	if w.WaCli == nil {
		return errors.WithCode(code.ClientNotInitialized, "Whatsapp client is not initialized")
	}
	var now = utils.GetCurrentTime()
	var defaultInterval = 5 * time.Minute
	expectExpireTime := now.Add(defaultInterval)
	for {
		now = utils.GetCurrentTime()
		if now.After(expectExpireTime) || now.Equal(expectExpireTime) {
			break
		}
		if !w.WaCli.IsConnected() {
			w.WaCli.Connect()
		}

		if w.WaCli.IsLoggedIn() {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if !w.WaCli.IsLoggedIn() {
		return errors.New("login fail")
	}

	return nil
}

func ParseJID(arg string) (types.JID, error) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), nil
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			fmt.Printf("invalid JID %s: %v", arg, err)
			return recipient, errors.WithCode(code.ErrInvalidJID, "")
		} else if recipient.User == "" {
			fmt.Printf("invalid JID %v: no server specified", arg)
			return recipient, errors.WithCode(code.ErrInvalidJID, "")
		}
		return recipient, nil
	}
}
