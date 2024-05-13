package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

var Config *Conf

type Conf struct {
	Metrics      MetricsConf `mapstructure:"metrics"`
	HttpPort     int         `mapstructure:"http_port"`
	GrpcPort     int         `mapstructure:"grpc_port"`
	DB           DBConf      `mapstructure:"db"`
	StreamServer StreamConf  `mapstructure:"streamserver"`
}

type MetricsConf struct {
	Enabled     bool   `mapstructure:"enabled" default:"true"`
	UseHttpConf bool   `mapstructure:"use_http_conf" default:"true"`
	BindAddress string `mapstructure:"bind_address" default:"127.0.0.1:9090"`
}

type DBConf struct {
	Driver   string `mapstructure:"driver"`
	DSN      string `mapstructure:"dsn"`
	FilePath string `mapstructure:"filepath"` // for sqlite only
}

type StreamConf struct {
	Secret string `mapstructure:"secret"`
}

func ParseConfig() (*Conf, error) {
	v := viper.GetViper()

	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.AutomaticEnv()
	v.SetEnvPrefix("STREAM_API")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	keys := getKeysFromStructType[Conf](".")
	for _, k := range keys {
		// fmt.Println("env bind: ", k)
		v.BindEnv(k)
	}

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	c := &Conf{}
	defaults.MustSet(c)
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

func getKeysFromStructType[S any](delim string) []string {
	var s = reflect.TypeFor[S]()
	if s.Kind() != reflect.Struct {
		panic("must input a struct type")
	}
	var keys = []string{}
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		keys = append(keys, getKeysFromAny(extractNameFromStructField(field), field.Type, delim)...)
	}
	return keys
}

func getKeysFromAny(prefix string, e reflect.Type, delim string) []string {
	if prefix == "" {
		return nil
	}
	switch e.Kind() {
	case reflect.Struct:
		keys := []string{}
		for i := 0; i < e.NumField(); i++ {
			field := e.Field(i)
			name := extractNameFromStructField(field)
			if name == "" {
				continue
			}
			keys = append(keys, getKeysFromAny(prefix+delim+name, field.Type, delim)...)
		}
		return keys
	case reflect.Pointer:
		return getKeysFromAny(prefix, e.Elem(), delim)
	default:
		return []string{prefix}
	}
}

func extractNameFromStructField(field reflect.StructField) string {
	name := strings.ToLower(field.Name)
	mapstructureTag := field.Tag.Get("mapstructure")
	if mapstructureTag != "" {
		tags := strings.SplitN(mapstructureTag, ",", 2)
		if len(tags) >= 1 {
			mName := strings.TrimSpace(tags[0])
			if mName == "-" {
				return ""
			} else if mName != "" {
				return mName
			}
		}

	}
	return name
}
