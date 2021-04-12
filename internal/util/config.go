package util

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
)

func ReadConfigs() {
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	var configFilePath = viper.GetString("REDUCTO_CONFIG_PATH")

	yamlContent, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		log.Error(err)
		panic(fmt.Errorf("error in Parsing Configration(): %s", err))
	}
	if err = viper.ReadConfig(bytes.NewBuffer(yamlContent)); err != nil {
		log.Error(err)
		panic(fmt.Errorf("error in ReadConfigs(): %s", err))
	}
	log.Info(viper.GetString("Environment"), " configuration set successfully")
}
