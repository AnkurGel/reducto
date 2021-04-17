package models

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// URL data model
type URL struct {
	gorm.Model
	Short    string `gorm:"type:VARCHAR(15);unique"`
	Original string `gorm:"type:varchar(2048);index:orig"`
	Retries  uint8
	Visits   []Visit
}

// ShortURL gives http URL string form for a short slug
func (u *URL) ShortURL() string {
	return viper.GetString("BaseURL") + u.Short
}
