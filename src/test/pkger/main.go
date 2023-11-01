package main

import (
	"context"
	"github.com/markbates/pkger"
	"github.com/rookie-ninja/rk-boot"
	// Must be present in order to make pkger load embedded files into memory.
)

func init() {
	// This is used while running pkger CLI
	//表示 把哪个目录下的文件加载到pkger来
	pkger.Include("./")
}

// Application entrance.
func main() {
	// Create a new boot instance.
	boot := rkboot.NewBoot()

	// Bootstrap
	boot.Bootstrap(context.Background())

	// Wait for shutdown sig
	boot.WaitForShutdownSig(context.Background())
}
