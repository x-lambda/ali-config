package conf

import (
	"bytes"
	"testing"

	"bou.ke/monkey"
	client "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	// TODO patch
	monkey.Patch(getConfig, func() (cli client.IConfigClient, content string) {
		return
	})
	defer monkey.Unpatch(getConfig)

	monkey.Patch(onConfigChange, func(cli client.IConfigClient) {})
	defer monkey.Unpatch(onConfigChange)
}

func TestViperConf(t *testing.T) {
	viper.SetConfigType("toml")
	data := `
a = 1
b = 2
c = 3
d = "123"
ff = 999
`
	r := bytes.NewReader([]byte(data))
	err := viper.ReadConfig(r)
	assert.Nil(t, err)

	assert.Equal(t, viper.GetInt("a"), 1)
	assert.Equal(t, viper.GetInt("b"), 2)
	assert.Equal(t, viper.GetInt("c"), 3)
	assert.Equal(t, viper.GetString("d"), "123")
	assert.Equal(t, viper.GetInt("ff"), 999)

	data = `
a = 100
b = 200
c = 300
d = "***"
`
	err = viper.ReadConfig(bytes.NewReader([]byte(data)))
	assert.Nil(t, err)

	assert.Equal(t, viper.GetInt("a"), 100)
	assert.Equal(t, viper.GetInt("b"), 200)
	assert.Equal(t, viper.GetInt("c"), 300)
	assert.Equal(t, viper.GetString("d"), "***")
	assert.Equal(t, viper.GetInt("ff"), 0)
}

func TestConf(t *testing.T) {
	// mock
}
