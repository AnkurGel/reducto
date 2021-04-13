package store

import (
	"errors"
	"fmt"
	"github.com/ankurgel/reducto/internal/models"
	"github.com/ankurgel/reducto/internal/redisdb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Store struct {
	Db          *gorm.DB
	redisClient *redisdb.Redis
}

// InitStore configures the store for connection, models,
// logging etc and returns instantiated store
func InitStoreWithCache(redisClient *redisdb.Redis) *Store {
	s := &Store{redisClient: redisClient}
	s.EstablishConnection()
	defer log.Info("Store configured successfully")
	s.SetupModels()
	return s
}

// GetDSN returns Data Source Name for sql configuration
func (s *Store) GetDSN() string {
	config := viper.GetStringMapString("Mysql")
	host, username, password, database := config["host"], config["username"], config["password"], config["database"]
	if viper.GetString("Environment") == "development" {
		return fmt.Sprintf("%s:%s@/%s?parseTime=true", username, password, database)
	}
	return fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true", username, password, host, database)
}

// EstablishConnection establishes connection of store with sql server
func (s *Store) EstablishConnection() {
	var err error
	s.Db, err = gorm.Open(mysql.Open(s.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect to DB: %s", err))
	}
}

// SetupModels setups and migrates all the models
func (s *Store) SetupModels() {
	s.Db.AutoMigrate(&models.URL{})
}

// CreateByLongURL interacts with database to create short URL
// and returns URL object or error
func (s *Store) CreateByLongURL(longURL string, customSlug string) (*models.URL, error) {
	var u models.URL
	var shortURL *models.URL
	var err error
	var shortHash string
	var retries uint8 = 0
	var lenCustomSlug int = len(customSlug)
	if lenCustomSlug > 0 {
		if lenCustomSlug > 15 {
			return nil, errors.New("length for custom URL cannot be more than 15 characters")
		}
		customExists := s.Db.Where("short = ?", customSlug).First(&u)
		if customExists.Error == nil {
			if u.Original == longURL {
				return &u, nil
			}
			return nil, fmt.Errorf("custom slug %s is already taken", customSlug)
		}
	}

	if result := s.Db.Where("original = ?", longURL).First(&u); result.Error != nil || lenCustomSlug > 0 {
		log.Error(result.Error)
		if lenCustomSlug > 0 {
			shortHash, err = customSlug, nil
			retries = uint8(viper.GetUint32("MaxRetries"))
		} else {
			shortHash, err = s.redisClient.GetKey()
		}

		if err != nil {
			return nil, err
		}

		shortURL, err = s.FindByShortURL(shortHash)
		for err == nil && retries < uint8(viper.GetUint32("MaxRetries")) {
			retries++
			shortHash, _ = s.redisClient.GetKey()
			shortURL, err = s.FindByShortURL(shortHash)
		}
		if shortURL == nil {
			short := models.URL{
				Short:    shortHash,
				Original: longURL,
				Retries:  retries,
			}
			if result := s.Db.Create(&short); result.Error != nil {
				return nil, errors.New("couldn't shorten. Something went wrong")
			}
			return &short, nil
		}
		return nil, errors.New("couldn't shorten. Out of lives")

	}
	return &u, nil
}

// FindByShortURL looks-up the store for given short url
// and returns URL object or error
func (s *Store) FindByShortURL(shortURL string) (*models.URL, error) {
	var u models.URL
	if result := s.Db.Where("short = ?", shortURL).First(&u); result.Error != nil {
		return nil, result.Error
	}
	return &u, nil
}
