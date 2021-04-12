package main

import (
	"github.com/ankurgel/reducto/internal/redisdb"
	"github.com/ankurgel/reducto/internal/router"
	"github.com/ankurgel/reducto/internal/store"
	"github.com/ankurgel/reducto/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	log.Info("Starting reducto-server...")
	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, os.Interrupt, os.Kill)
	util.InitLogger()
	util.ReadConfigs()
	config := viper.GetStringMap("Redis")
	redisClient, err := redisdb.New(config["address"].(string), config["db"].(int))
	if err != nil {
		panic(err)
	}
	s := store.InitStoreWithCache(redisClient)
	r := router.InitRouter(s)
	r.Engine.Run(":8081")
}
