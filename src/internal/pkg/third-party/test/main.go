package main

import (
	watest "go-socialapp/internal/pkg/third-party/test/client"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	phone := "62838520048745"
	priv := "Vtzw+sul0qUAGXI32sWnrbcIBbX3uSFB5AsS/YMtdD4="
	pub := "wKR5KXJG7kDRjqI1F1a2xfCDXqeRsyaJK71Z9Xzyn3o="
	idEcode := "NjI4Mzg1MjAwNDg3NDUjsUdxHp0pLsnkp7KanY8gBqxxnxw="
	watest.MainForTest(priv, pub, phone, 12, idEcode)
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

}
