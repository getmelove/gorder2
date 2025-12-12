package main

import (
	"log"

	"github.com/getmelove/gorder2/common/config"
	"github.com/spf13/viper"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func main() {
	log.Printf("%v", viper.Get("order"))
}
