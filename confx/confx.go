package confx

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/nacos-group/nacos-sdk-go/clients"
	client "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// PublicAddr 公网（测试）
	PublicAddr = "acm.aliyun.com"
	// HZAddr 华东 1（杭州）
	HZAddr = "addr-hz-internal.edas.aliyun.com"
	// QDAddr 华北 1（青岛）
	QDAddr = "addr-qd-internal.edas.aliyun.com"
	// SHAddr 华东 2（上海）
	SHAddr = "addr-sh-internal.edas.aliyun.com"
	// BJAddr 华北 2（北京）
	BJAddr = "addr-bj-internal.edas.aliyun.com"
	// SZAddr 华南 1（深圳）
	SZAddr = "addr-sz-internal.edas.aliyun.com"
	// HKAddr 香港
	HKAddr = "addr-hk-internal.edas.aliyuncs.com"
	// SingaporeAddr 新加坡
	SingaporeAddr = "addr-singapore-internal.edas.aliyun.com"
	// ApAddr 澳大利亚（悉尼）
	ApAddr = "addr-ap-southeast-2-internal.edas.aliyun.com"
	// USWest1Addr 美国（硅谷）
	USWest1Addr = "addr-us-west-1-internal.acm.aliyun.com"
	// USEast1Addr 美国（弗吉尼亚）
	USEast1Addr = "addr-us-east-1-internal.acm.aliyun.com"
	// ShanghaiFinance1Addr 华东 2（上海）金融云
	ShanghaiFinance1Addr = "addr-cn-shanghai-finance-1-internal.edas.aliyun.com"
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

var endpoint string    // os.Getenv("ALI_CONFIG_ENDPOINT")
var namespaceId string // os.Getenv("ALI_CONFIG_NAMESPACE_ID")
var accessKey string   // os.Getenv("ALI_CONFIG_ACCESS_KEY")
var secretKey string   // os.Getenv("ALI_CONFIG_SECRET_KEY")
var dataID string      // os.Getenv("ALI_CONFIG_DATA_ID")
var group string       // os.Getenv("ALI_CONFIG_GROUP")

type Conf struct {
	viper *viper.Viper
}

var c Conf

func init() {
	getConfig()
}

// kmsDecrypt kms解密 content：密文内容
func kmsDecrypt(accessKey, secretKey, content string) (text string, err error) {
	request := kms.CreateDecryptRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.AcceptFormat = "json"
	request.CiphertextBlob = content
	request.ConnectTimeout = 1 * time.Minute

	// TODO regionId?
	kc, err := kms.NewClientWithAccessKey("acm", accessKey, secretKey)
	if err != nil {
		return
	}

	resp, err := kc.Decrypt(request)
	if err != nil {
		return
	}

	text = resp.Plaintext
	return
}

// sliceHeader is the runtime representation of a slice.
// It should be identical to reflect.sliceHeader
type sliceHeader struct {
	data     unsafe.Pointer
	sliceLen int
	sliceCap int
}

// stringHeader is the runtime representation of a string.
// It should be identical to reflect.StringHeader
type stringHeader struct {
	data      unsafe.Pointer
	stringLen int
}

// ByteToString unsafely converts b into a string.
// If you modify b, then s will also be modified. This violates the
// property that strings are immutable.
func ByteToString(b []byte) (s string) {
	sliceHeader := (*sliceHeader)(unsafe.Pointer(&b))
	stringHeader := (*stringHeader)(unsafe.Pointer(&s))
	stringHeader.data = sliceHeader.data
	stringHeader.stringLen = len(b)
	return s
}

func getConfig() (cli client.IConfigClient, content string) {
	endpoint = os.Getenv("ALI_CONFIG_ENDPOINT")
	if endpoint == "" {
		endpoint = PublicAddr // 公网测试
	}
	namespaceId = os.Getenv("ALI_CONFIG_NAMESPACE_ID")
	accessKey = os.Getenv("ALI_CONFIG_ACCESS_KEY")
	secretKey = os.Getenv("ALI_CONFIG_SECRET_KEY")
	dataID = os.Getenv("ALI_CONFIG_DATA_ID")
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

	content, err = cli.GetConfig(vo.ConfigParam{DataId: dataID, Group: group})
	if err != nil {
		panic(err)
	}
	fmt.Println(content)

	body := ByteToString([]byte(content))
	switch {
	case strings.HasPrefix(dataID, "cipher-kms-aes-128-"):
	case strings.HasPrefix(dataID, "cipher-"):
		cs, err := kmsDecrypt(accessKey, secretKey, body)
		if err != nil {
			panic(err)
		}
		fmt.Println("解密后的数据: ", cs)
	default:
	}
	return
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
