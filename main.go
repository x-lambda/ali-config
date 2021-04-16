package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	endpoint := os.Getenv("ali_config_endpoint")
	namespaceId := os.Getenv("ali_config_namespace_id")
	accessKey := os.Getenv("ali_config_access_key")
	secretKey := os.Getenv("ali_config_secret_key")
	dataId := os.Getenv("ali_config_data_id")
	group := os.Getenv("ali_config_group")

	sc := []constant.ServerConfig{
		{
			IpAddr: "console.nacos.io",
			Port:   80,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         namespaceId,
		TimeoutMs:           100000,
		NotLoadCacheAtStart: true,
		RotateTime:          "1h",
		//OpenKMS:             false,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Endpoint:  endpoint + ":8080",
		LogDir:    "",
		CacheDir:  "",
	}

	// a more graceful way to create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	content, err := client.GetConfig(vo.ConfigParam{
		// DataId: "cipher-olympus",
		DataId: dataId,
		Group:  group,
	})
	fmt.Println("GetConfig,config :" + content)

	go func() {
		err = client.ListenConfig(vo.ConfigParam{
			DataId: dataId,
			Group:  group,
			OnChange: func(namespace, group, dataId, data string) {
				fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			},
		})
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	for {
		select {
		case <-stop:
			os.Exit(0)
		}
	}
}
