package main

import (
	"context"
	rkboot "github.com/rookie-ninja/rk-boot"
	_ "github.com/rookie-ninja/rk-boot/gf"
)

// Application entrance.
// 参考 https://zhuanlan.zhihu.com/p/441987879
// 访问 http://localhost:8080/rk/v1/static
// 对应boot.yaml必须放在src目录下吗
func main() {
	// Create a new boot instance.
	boot := rkboot.NewBoot()

	// Bootstrap
	boot.Bootstrap(context.Background())

	// Wait for shutdown sig
	boot.WaitForShutdownSig(context.Background())
}
