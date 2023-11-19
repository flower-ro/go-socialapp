package watest

import (
	"github.com/davecgh/go-spew/spew"
	"go.mau.fi/whatsmeow/socket"
	"go.mau.fi/whatsmeow/util/keys"
	waLog "go.mau.fi/whatsmeow/util/log"
	"net/http"
	"net/url"
)

// var cli *whatsmeow.Client
var log waLog.Logger

var logLevel = "INFO"

func MainForTest(priv string, user string, index int) {
	logLevel = "DEBUG"
	//store.DeviceProps.RequireFullSync = proto.Bool(true)
	log = waLog.Stdout("Main", logLevel, true)

	fs := socket.NewFrameSocket(log.Sub("Socket"), socket.WAConnHeader, proxy)
	//fs := socket.NewFrameSocket(log.Sub("Socket"), socket.WAConnHeader, nil)

	if err := fs.Connect(); err != nil {
		log.Errorf("Failed to Connect: %v", err)
		fs.Close(0)
		return
	} else if err = doHandshake(fs, *keys.NewKeyPair(), priv, user, index); err != nil {
		fs.Close(0)
		log.Errorf("Failed to noise handshake: %v", err)
	}
}

// (https://web.whatsapp.com/ws/chat)
func proxy(req *http.Request) (*url.URL, error) {
	//spew.Dump(req.URL)
	//// 创建代理URL
	//u, _ := url.Parse("http://signaler-pa.clients6.google.com")
	//preq := httputil.ProxyRequest{Out: req}
	//
	//preq.SetURL(u)
	//spew.Dump("preq.Out.URL=", preq.Out.URL)

	spew.Dump("----------"+
		"--request111=", req)
	return nil, nil

	//tran:=&http.Transport{
	//	Proxy: http.ProxyURL(proxyURL),
	//}

	//return nil, nil
}
