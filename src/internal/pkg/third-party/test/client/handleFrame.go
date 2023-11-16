package watest

import (
	"encoding/hex"
	"github.com/davecgh/go-spew/spew"
	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/socket"
)

func handleFrame(data []byte) {
	decompressed, err := waBinary.Unpack(data)
	if err != nil {
		log.Warnf("Failed to decompress frame: %v", err)
		log.Debugf("Errored frame hex: %s", hex.EncodeToString(data))
		return
	}
	node, err := waBinary.Unmarshal(decompressed)
	if err != nil {
		log.Warnf("Failed to decode node in frame: %v", err)
		log.Debugf("Errored frame hex: %s", hex.EncodeToString(decompressed))
		return
	}
	log.Debugf("%s", node.XMLString())
	spew.Dump("--------------node.tag---", node.Tag)
	spew.Dump("--------------node---", node)
	//if node.Tag == "xmlstreamend" {
	//	if !cli.isExpectedDisconnect() {
	//		log.Warnf("Received stream end frame")
	//	}
	//	// TODO should we do something else?
	//} else if cli.receiveResponse(node) {
	//	// handled
	//} else if _, ok := cli.nodeHandlers[node.Tag]; ok {
	//	select {
	//	case cli.handlerQueue <- node:
	//	default:
	//		log.Warnf("Handler queue is full, message ordering is no longer guaranteed")
	//		go func() {
	//			cli.handlerQueue <- node
	//		}()
	//	}
	//} else if node.Tag != "ack" {
	//	cli.Log.Debugf("Didn't handle WhatsApp node %s", node.Tag)
	//}
}

func onDisconnect(ns *socket.NoiseSocket, remote bool) {

	log.Debugf("OnDisconnect() called ")
	//ns.Stop(false)
	//cli.socketLock.Lock()
	//defer cli.socketLock.Unlock()
	//if cli.socket == ns {
	//	cli.socket = nil
	//	cli.clearResponseWaiters(xmlStreamEndNode)
	//	if !cli.isExpectedDisconnect() && remote {
	//		cli.Log.Debugf("Emitting Disconnected event")
	//		go cli.dispatchEvent(&events.Disconnected{})
	//		go cli.autoReconnect()
	//	} else if remote {
	//		cli.Log.Debugf("OnDisconnect() called, but it was expected, so not emitting event")
	//	} else {
	//		cli.Log.Debugf("OnDisconnect() called after manual disconnection")
	//	}
	//} else {
	//	cli.Log.Debugf("Ignoring OnDisconnect on different socket")
	//}
}
