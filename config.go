package tradestation

import (
	"github.com/blbgo/general"
)

const configSection = "Tradestation"
const configBaseURL = "BaseURL"
const configClientID = "ClientID"
const configClientSecret = "ClientSecret"

func (r *tradestation) LoadConfig(config general.Config) error {
	var err error

	r.baseURL, err = config.Value(configSection, configBaseURL)
	if err != nil {
		return err
	}
	r.clientID, err = config.Value(configSection, configClientID)
	if err != nil {
		return err
	}
	r.clientSecret, err = config.Value(configSection, configClientSecret)
	return err
}
