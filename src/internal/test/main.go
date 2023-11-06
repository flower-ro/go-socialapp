package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/marmotedu/iam/pkg/log"
	whatsappbase "go-socialapp/internal/socialserver/enter/ws"
	"os"
	"path/filepath"
	"time"
)

func main1() {
	c := make(chan int, 3)
	c <- 10
	c <- 20
	//close(c)
	for i := 0; i < 5; i++ {
		x, ok := <-c //ok是 true表示 通道里还有数据，通道可能关闭或者未关闭都有可能
		// ok 是false,表示通道里没有数据了，且此时通道是关闭的，返回的x是零值
		//如果通道未关闭且没有数据了就会堵塞
		fmt.Println(i, ":", ok, x)
	}
}

type a struct {
	closeFlag bool
}

func (s *a) cicle() {
	for {
		if s.closeFlag == true {
			break
		}
		fmt.Println("111")
		time.Sleep(2 * time.Second)
	}
}
func main2() {
	ss := &a{}

	go func() {
		ss.cicle()
		fmt.Println("结束")
	}()

	time.Sleep(4 * time.Second)
	ss.closeFlag = true
	select {}
}

func main3() {
	defer func() {
		if err := recover(); err != nil {
			log.Infof("fas err %v", err)
		}
	}()
	send := make(chan whatsappbase.BroadcastMessage, 10)
	close(send)
	close(send)
}

func main() {
	pathQrCode := "E:/software/sqlite3/dbx"
	//CreateFolder(pathQrCode)

	spew.Dump(os.Stat(pathQrCode))

}
func CreateFolder(folderPath ...string) error {
	for _, folder := range folderPath {
		newFolder := filepath.Join(folder)
		spew.Dump(newFolder)
		_, err := os.Create(newFolder)
		if err != nil {
			return err
		}
	}
	return nil
}
