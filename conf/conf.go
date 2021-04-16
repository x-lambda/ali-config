package conf

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	client "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// Hostname 主机名 服务树-service_name-rnd
	Hostname = "localhost"
	// AppID service name
	AppID = "localapp"
	// Env dev/uat/pre/prod
	Env = "dev"
	// Zone 区域
	Zone = "sh001"
)

// var path string // 配置文件所在目录

var endpoint string    // os.Getenv("ALI_CONFIG_ENDPOINT")
var namespaceId string // os.Getenv("ALI_CONFIG_NAMESPACE_ID")
var accessKey string   // os.Getenv("ALI_CONFIG_ACCESS_KEY")
var secretKey string   // os.Getenv("ALI_CONFIG_SECRET_KEY")
var dataId string      // os.Getenv("ALI_CONFIG_DATA_ID")
var group string       // os.Getenv("ALI_CONFIG_GROUP")

type Conf struct {
	viper *viper.Viper
}

var c Conf

func init() {
	Hostname, _ = os.Hostname()
	if appID := os.Getenv("APP_ID"); appID != "" {
		AppID = appID
	} else {
		logger().Warn("env APP_ID is empty")
	}

	if env := os.Getenv("DEPLOY_ENV"); env != "" {
		Env = env
	} else {
		logger().Warn("env DEPLOY_ENV is empty")
	}

	if zone := os.Getenv("ZONE"); zone != "" {
		Zone = zone
	} else {
		logger().Warn("env ZONE is empty")
	}

	//path = os.Getenv("CONF_PATH")
	//if path == "" {
	//	logger().Warn("env CONF_PATH is empty")
	//	var err error
	//	if path, err = os.Getwd(); err != nil {
	//		panic(err)
	//	}
	//	logger().WithField("path", path).Info("use default conf path")
	//}

	cli, content := getConfig()
	c.viper = viper.New()
	err := c.viper.ReadConfig(bytes.NewReader([]byte(content)))
	if err != nil {
		panic(err)
	}
	logger().Infof("init with config: %s", content)

	go onConfigChange(cli)
}

// getConfig 首次启动进程时，获取配置文件信息
func getConfig() (cli client.IConfigClient, content string) {
	// 阿里云acm配置
	endpoint = os.Getenv("ALI_CONFIG_ENDPOINT")
	namespaceId = os.Getenv("ALI_CONFIG_NAMESPACE_ID")
	accessKey = os.Getenv("ALI_CONFIG_ACCESS_KEY")
	secretKey = os.Getenv("ALI_CONFIG_SECRET_KEY")
	dataId = os.Getenv("ALI_CONFIG_DATA_ID")
	group = os.Getenv("ALI_CONFIG_GROUP")

	// 初始化配置
	sc := []constant.ServerConfig{{IpAddr: "console.nacos.io", Port: 80}}
	cc := constant.ClientConfig{
		NamespaceId:         namespaceId,
		TimeoutMs:           100000,
		NotLoadCacheAtStart: true,
		RotateTime:          "1h",
		AccessKey:           accessKey,
		SecretKey:           secretKey,
		Endpoint:            endpoint + ":8080",
		LogDir:              "",
		CacheDir:            "",
		//OpenKMS:             false,
	}

	var err error
	// a more graceful way to create config client
	cli, err = clients.NewConfigClient(vo.NacosClientParam{ClientConfig: &cc, ServerConfigs: sc})
	if err != nil {
		panic(err)
	}

	content, err = cli.GetConfig(vo.ConfigParam{DataId: dataId, Group: group})
	if err != nil {
		panic(err)
	}

	return
}

// onConfigChange 监听文件的变更
func onConfigChange(cli client.IConfigClient) {
	err := cli.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		// 监听文件变更
		OnChange: func(ns, g, did, data string) {
			// fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			logger().Infof("[conf][listen] namespace: %s, dataId: %s, group: %s, on change: %s",
				ns, did, g, data)

			// 修改配置
			err := c.viper.ReadConfig(bytes.NewReader([]byte(data)))
			if err != nil {
				logger().Errorf("[ali-config][ListenConfig] err: %+v", err)
			}
		},
	})

	if err != nil {
		panic(err)
	}
}

var levels = map[string]logrus.Level{
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
	"error": logrus.ErrorLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
}

func logger() *logrus.Entry {
	if level, ok := levels[os.Getenv("LOG_LEVEL")]; ok {
		logrus.SetLevel(level)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return logrus.WithFields(logrus.Fields{
		"app_id":      AppID,
		"instance_id": Hostname,
		"env":         Env,
	})
}

// GetFloat64 获取浮点数配置
func GetFloat64(key string) float64 {
	return c.viper.GetFloat64(key)
}

// Get 获取字符串配置
func Get(key string) string {
	return c.viper.GetString(key)
}

// GetStrings 获取字符串列表
func GetStrings(key string) (s []string) {
	value := Get(key)
	if value == "" {
		return
	}

	for _, v := range strings.Split(value, ",") {
		s = append(s, v)
	}
	return
}

// GetInt32s 获取数字列表
// 1,2,3 => []int32{1,2,3}
func GetInt32s(key string) (s []int32, err error) {
	s64, err := GetInt64s(key)
	for _, v := range s64 {
		s = append(s, int32(v))
	}
	return
}

// GetInt64s 获取数字列表
func GetInt64s(key string) (s []int64, err error) {
	value := Get(key)
	if value == "" {
		return
	}

	var i int64
	for _, v := range strings.Split(value, ",") {
		i, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return
		}
		s = append(s, i)
	}
	return
}

// GetInt 获取整数配置
func GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetInt32 获取 int32 配置
func GetInt32(key string) int32 {
	return c.viper.GetInt32(key)
}

// GetInt64 获取 int64 配置
func GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

// GetDuration 获取时间配置
func GetDuration(key string) time.Duration {
	return c.viper.GetDuration(key)
}

// GetTime 查询时间配置
// 默认时间格式为 "2006-01-02 15:04:05"，conf.GetTime("FOO_BEGIN")
// 如果需要指定时间格式，则可以多传一个参数，conf.GetString("FOO_BEGIN", "2006")
//
// 配置不存在或时间格式错误返回**空时间对象**
// 使用本地时区
func GetTime(key string, args ...string) time.Time {
	fmt := "2006-01-02 15:04:05"
	if len(args) == 1 {
		fmt = args[0]
	}

	t, _ := time.ParseInLocation(fmt, c.viper.GetString(key), time.Local)
	return t
}

// GetBool 获取配置布尔配置
func GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// Set 设置配置，仅用于测试
func Set(key string, value string) {
	c.viper.Set(key, value)
}
