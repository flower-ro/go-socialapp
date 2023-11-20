package main

import (
	watest "go-socialapp/internal/pkg/third-party/test/client"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	phone := "62838505287608"
	priv := "wOCzfbzMOfWiEJwYjA7jQAZENnIPw41GJbC4in7Yf00="
	pub := "9ybEH7eD3agLTDo0BDySJt57o8pKR5/9EwiRPyhUjkw="
	idEcode := "NjI4Mzg1MjAwNDg3NDUjsUdxHp0pLsnkp7KanY8gBqxxnxw="
	watest.MainForTest(priv, pub, phone, 23, idEcode)
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

}
