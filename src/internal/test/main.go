package main

import "fmt"

func main() {
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
