package tradestation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// SymbolInfoModel result from Get Symbol Info call, format from doc:
// {
// 	"Category": "string",
// 	"Country": "US",
// 	"Currency": "USD",
// 	"Description": "string",
// 	"DisplayType": 0,
// 	"Error": "string",
// 	"Exchange": "NYSE",
// 	"ExchangeID": 0,
// 	"ExpirationDate": "string",
// 	"ExpirationType": "string",
// 	"FutureType": "string",
// 	"IndustryCode": "string",
// 	"IndustryName": "string",
// 	"LotSize": 0,
// 	"MinMove": 0,
// 	"Name": "string",
// 	"OptionType": "string",
// 	"PointValue": 0,
// 	"Root": "string",
// 	"SectorName": "string",
// 	"StrikePrice": 0,
// 	"Underlying": "string"
// }
type SymbolInfoModel struct {
	Category       string
	Country        string
	Currency       string
	Description    string
	DisplayType    int
	Error          string
	Exchange       string
	ExchangeID     int
	ExpirationDate string
	ExpirationType string
	FutureType     string
	IndustryCode   string
	IndustryName   string
	LotSize        int
	MinMove        float32
	Name           string
	OptionType     string
	PointValue     float32
	Root           string
	SectorName     string
	StrikePrice    float32
	Underlying     string
}

const symbolInfoPathFormat = "/data/symbol/%v?%v"

func (r *tradestation) SymbolInfo(symbol string) (*SymbolInfoModel, error) {
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return nil, fmt.Errorf("Argument empty: %v", "symbol")
	}
	err := r.AccessTokenOK()
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Set("access_token", r.AccessToken)
	data.Set("APIVersion", apiVersion)

	req, err := http.NewRequest(
		http.MethodGet,
		r.baseURL+fmt.Sprintf(symbolInfoPathFormat, symbol, data.Encode()),
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := r.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, ErrStatusNotOK
	}

	// remove this section and use commented out NewDecoder line once some example formats are
	// gathered.
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = r.DumperFactory.Dump("SymbolInfo"+symbol, bodyBytes)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))

	//decoder := json.NewDecoder(resp.Body)
	symbolInfoModel := &SymbolInfoModel{}
	err = decoder.Decode(symbolInfoModel)
	if err != nil {
		return nil, err
	}

	return symbolInfoModel, nil
}
