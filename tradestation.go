package tradestation

import (
	"errors"
	"net/http"
	"time"

	"github.com/blbgo/general"
)

// Tradestation provides access to the Tradestation API
type Tradestation interface {
	AccessTokenStatus() string
	AccessTokenOK() error
	StartAuth(redirectURL string) (string, error)
	FinishAuth(code string) error

	SymbolInfo(symbol string) (*SymbolInfoModel, error)

	DailyStartingOn(symbol string, start time.Time) error
	//MakeGetRequest(url string, form url.Values, result interface{}) error
	//MakePostRequest(url string, form url.Values, body interface{}, result interface{}) error

	//EtradeTimeAsGoTime(etradeTime int64) time.Time

	//Accounts() ([]Account, error)
	//Transactions(accountIDKey string, start time.Time) ([]Transaction, error)
}

type tradestation struct {
	general.LoggerFactory
	general.DumperFactory
	http.RoundTripper
	general.PersistentState

	baseURL      string
	clientID     string
	clientSecret string

	authState
}

type authState struct {
	RefreshToken   string
	AccessToken    string
	Expires        time.Time
	RedredirectURI string
}

// ErrNoAccessToken no access token
var ErrNoAccessToken = errors.New("No access token")

// ErrStatusNotOK response status is not ok
var ErrStatusNotOK = errors.New("Response status is not ok")

const apiVersion = "20160101"
const dateFormat = "01-02-2006"

// New provides an implementation of the Tradestation interface
func New(
	config general.Config,
	loggerFactory general.LoggerFactory,
	dumperFactory general.DumperFactory,
	roundTripper http.RoundTripper,
	persistentState general.PersistentState,
) (Tradestation, error) {
	r := &tradestation{
		LoggerFactory:   loggerFactory,
		DumperFactory:   dumperFactory,
		RoundTripper:    roundTripper,
		PersistentState: persistentState,
	}

	err := r.LoadConfig(config)
	if err != nil {
		return nil, err
	}

	return r, nil
}

//const timeFormatMMDDYYYY = "01022006"
//func (r *tradestation) EtradeTimeAsGoTime(etradeTime int64) time.Time {
//	return time.Unix(etradeTime/1000, 0).UTC()
//}
