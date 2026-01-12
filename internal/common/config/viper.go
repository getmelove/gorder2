package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

func init() {
	if err := NewViperConfig(); err != nil {
		panic(err)
	}
}

var once sync.Once

func NewViperConfig() error {
	once.Do(func() {
		err := newViperConfig()
		if err != nil {
			panic(err)
		}
	})
	return nil
}

func newViperConfig() error {
	relPath, err := getRelativePathFromCaller()
	if err != nil {
		return err
	}
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	// 这个是要在其他服务下面使用，所以路径是这样的
	viper.AddConfigPath(relPath)
	viper.EnvKeyReplacer(strings.NewReplacer("-", "_"))
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY")
	_ = viper.BindEnv("endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}

func getRelativePathFromCaller() (relPath string, err error) {
	callerPwd, err := os.Getwd()
	if err != nil {
		return
	}
	_, filename, _, _ := runtime.Caller(0)
	relPath, _ = filepath.Rel(callerPwd, filepath.Dir(filename))
	fmt.Printf("caller from: %s, here: %s\n", filename, callerPwd)
	return
}
