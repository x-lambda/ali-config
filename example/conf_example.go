package main

import (
	"github.com/x-lambda/ali-config/conf"
)

func main() {
	// 先配置好环境变量
	conf.Get("EXAMPLE_A")        // 获取配置 EXAMPLE_A(string) 的值
	conf.GetInt32("EXAMPLE_B")   // 获取配置 EXAMPLE_B(int32) 的值
	conf.GetStrings("EXAMPLE_C") // 获取配置 EXAMPLE_C([]string) 的值
}
