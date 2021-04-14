package util

import (
	"fmt"
	"github.com/ankurgel/reducto/internal/store"
	"github.com/goware/urlx"
	"github.com/spf13/viper"
	"net/url"
	"strings"
)

// URLValidationError represents standard url validation exception
type URLValidationError struct {
	Reason string
}

func (e *URLValidationError) Error() string {
	return fmt.Sprintf("URLValidationError: %s", e.Reason)
}

// NormalizeURL normalizes the URL or throws error if URl is not parsable
func NormalizeURL(str string, s *store.Store) (string, error) {
	strings.Replace(str, " ", "", -1)
	if len(str) < 4 || len(str) > 2048 {
		return "", &URLValidationError{"URL has inadequate length"}
	}
	var val, baseUrl *url.URL
	var err error
	var isHostBanned bool
	val, err = urlx.Parse(str)
	if err != nil {
		return "", &URLValidationError{"Cannot parse. Invalid URL."}
	}
	baseUrl, err = urlx.Parse(viper.GetString("BaseURL"))
	if err != nil {
		return "", &URLValidationError{"Problem with parsing base URL."}
	}
	if baseUrl.Host == val.Host {
		return "", &URLValidationError{"This is already a shortened link."}
	}

	isHostBanned, err = s.IsHostBanned(val.Host)
	if err != nil {
		return "", &URLValidationError{err.Error()}
	}
	if isHostBanned {
		return "", &URLValidationError{"URL Domain is banned."}
	}

	normalize, err := urlx.Normalize(val)
	if err != nil {
		return "", &URLValidationError{"Cannot normalize URL"}
	}
	return normalize, nil
}
