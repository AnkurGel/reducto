package main

import (
	"github.com/ankurgel/reducto/internal/keygen"
	"github.com/ankurgel/reducto/internal/redisdb"
	"github.com/ankurgel/reducto/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Info("Initializing reducto-keygen...")
	util.InitLogger()
	util.ReadConfigs()
	config := viper.GetStringMap("Redis")
	redisClient, err := redisdb.New(config["address"].(string), config["db"].(int))
	if err != nil {
		panic(err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	keysCountToAdd := make(chan int64)
	resumePoller := make(chan bool, 1)

	resumePoller <- true
	go addKeysInRedis(redisClient, keysCountToAdd, resumePoller)
	go checkPool(redisClient, keysCountToAdd, resumePoller)
	<-quit
	log.Info("Exiting from reducto-keygen")
}

func addKeysInRedis(redisClient *redisdb.Redis, keysCountToAdd chan int64, resumePoller chan bool) {
	for {
		var i int64 = 0
		keysToAdd := <-keysCountToAdd
		for ; i < keysToAdd; i++ {
			err := redisClient.SaveKey(keygen.GenerateKey())
			if err != nil {
				panic(err)
			}
		}
		resumePoller <- true
	}
}

func checkPool(redisClient *redisdb.Redis, keysCountToAdd chan int64, resumePoller chan bool) {
	for {
		<-resumePoller
		length, err := redisClient.KeyPoolSize()
		if err != nil {
			panic(err)
		}
		if length < 5000 {
			log.Info("Sending signal to add ", 5000-length, " values")
			keysCountToAdd <- 5000 - length
		} else {
			keysCountToAdd <- 0
		}
		time.Sleep(time.Second * 1)
	}
}
