package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	viper.AddConfigPath("..")
	t.Run("Parse config from env", func(t *testing.T) {
		envs := [][2]string{
			// {"STREAM_API_METRICS_ENABLED", "false"},
			{"STREAM_API_METRICS_USE_HTTP_CONF", "false"},
			{"STREAM_API_METRICS_BIND_ADDRESS", "0.0.0.0:9099"},
		}

		for _, env := range envs {
			os.Setenv(env[0], env[1])
		}
		defer func() {
			for _, env := range envs {
				os.Unsetenv(env[0])
			}
		}()

		cfg, err := ParseConfig()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, true, cfg.Metrics.Enabled)
		assert.Equal(t, false, cfg.Metrics.UseHttpConf)
		assert.Equal(t, "0.0.0.0:9099", cfg.Metrics.BindAddress)

	})
}
