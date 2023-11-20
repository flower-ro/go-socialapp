package main

import (
	watest "go-socialapp/internal/pkg/third-party/test/client"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	phone := "62838505287608"
	priv := "wIN+kEFJIw1q5WfgaSmT4c7ENtFFDpxNlgBUObUrVHg="
	pub := "Y9L5mA0YMH+SFq6J6+7EZwWRcccSuFPaSjlOJSSbXB4="
	idEcode := "NjI4Mzg1MjAwNDg3NDUjsUdxHp0pLsnkp7KanY8gBqxxnxw="
	watest.MainForTest(priv, pub, phone, 1, idEcode)
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

}
