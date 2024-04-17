package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var Config *Conf

type Conf struct {
	Metrics      MetricsConf `yaml:"metrics"`
	HttpPort     int         `yaml:"http_port"`
	GrpcPort     int         `yaml:"grpc_port"`
	DB           DBConf      `yaml:"db"`
	StreamServer StreamConf  `yaml:"streamserver"`
}

type MetricsConf struct {
	Enabled     bool   `mapstructure:"enabled" mapstructure_default:"true"`
	UseHttpConf bool   `mapstructure:"use_http_conf" default:"true"`
	BindAddress string `mapstructure:"bind_address" default:"127.0.0.1:9090"`
}

type DBConf struct {
	Driver   string `yaml:"driver"`
	DSN      string `yaml:"dsn"`
	FilePath string `yaml:"filepath"` // for sqlite only
}

type StreamConf struct {
	Secret string `yaml:"secret"`
}

func ParseConfig() (*Conf, error) {
	v := viper.GetViper()

	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.AutomaticEnv()
	v.SetEnvPrefix("STREAM_API")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	c := &Conf{}
	// defaults.MustSet(c)
	fmt.Println(v.Get("metrics.bind_address"))

	err = v.Unmarshal(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func InitConfig() error {
	c, err := ParseConfig()
	if err != nil {
		return err
	}
	Config = c
	return nil
}
