# ali-config

<br>

使用阿里云ACM作为配置中心
> `https://help.aliyun.com/document_detail/59956.html`

<br>

基于官方提供的`go-sdk`
> `https://github.com/nacos-group/nacos-sdk-go`

<br>

#### 无加密版`https://github.com/x-lambda/ali-config/tree/main/conf`
1. 现在阿里云上创建对应的资源
![img](./imgs/WX20210416-112736@2x.png)

2. 将对应的资源配置写到环境变量中
```shell
$ export ALI_CONFIG_ENDPOINT=""
$ export ALI_CONFIG_NAMESPACE_ID=""
$ export ALI_CONFIG_DATA_ID=""
$ export ALI_CONFIG_GROUP=""
$ export ALI_CONFIG_ACCESS_KEY=""
$ export ALI_CONFIG_SECRET_KEY=""
```
3. 直接在代码中引用即可
```go
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
```

#### 加密版`https://github.com/x-lambda/ali-config/tree/main/confx`
TODO: 解密

### 使用
```shell
$ go get -u github.com/x-lambda/ali-config
```

refer:

    https://github.com/xiaojiaoyu100/aliyun-acm



