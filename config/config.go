package config

import "github.com/spf13/viper"

var Config *Conf

type Conf struct {
	HttpPort     int        `yaml:"http_port"`
	GrpcPort     int        `yaml:"grpc_port"`
	DB           DBConf     `yaml:"db"`
	StreamServer StreamConf `yaml:"streamserver"`
}

type DBConf struct {
	Driver   string `yaml:"driver"`
	DSN      string `yaml:"dsn"`
	FilePath string `yaml:"filepath"` // for sqlite only
}

type StreamConf struct {
	Secret string `yaml:"secret"`
}

func InitConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("STREAM_API")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	Config = &Conf{}
	err = viper.Unmarshal(Config)
	if err != nil {
		return err
	}
	return nil
}
