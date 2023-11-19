package main

import (
	watest "go-socialapp/internal/pkg/third-party/test/client"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	phone := "62838386510404"
	priv := "kH64bBtlfWRZpbnGFHZY1KqzCuj8LVAi9vBk28FiIVE="
	watest.MainForTest(priv, phone, 0)
	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

}
