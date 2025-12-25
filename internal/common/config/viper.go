package config

import (
	"strings"

	"github.com/spf13/viper"
)

func NewViperConfig() error {
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	// 这个是要在其他服务下面使用，所以路径是这样的
	viper.AddConfigPath("../common/config")
	viper.EnvKeyReplacer(strings.NewReplacer("-", "_"))
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}
