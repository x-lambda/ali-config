package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/x-lambda/ali-config/conf"
)

func main() {
	// 先配置好环境变量
	fmt.Println("GetConfig config: ", conf.GetInt("a"))
	fmt.Println("GetConfig config: ", conf.GetInt("b"))
	fmt.Println("GetConfig config: ", conf.GetInt("c"))
	fmt.Println("GetConfig config: ", conf.GetInt("d"))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		select {
		case <-stop:
			os.Exit(0)
		}
	}
}
