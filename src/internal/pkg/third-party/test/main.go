package main

import (
	watest "go-socialapp/internal/pkg/third-party/test/client"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	phone := "62838520048745"
	priv := "6P/aUIWsxIn/+MUq1BTLrex8S4PG/MrugtI/IszcbWk="
	pub := "udD63BQitDhUP00YB2KVdGrgP8jneCUmg7mkkGcEDUw="
	idEcode := "NjI4Mzg1MjAwNDg3NDUjsUdxHp0pLsnkp7KanY8gBqxxnxw="
	watest.MainForTest(priv, pub, phone, 12, idEcode)
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

}
