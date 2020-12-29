package tradestation

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/blbgo/general"
)

// DailyStartingOnModel result item from Stream BarChart - Starting on Date, format from doc:
// {
// 	"Close": 0,
// 	"DownTicks": 0,
// 	"DownVolume": 0,
// 	"High": 0,
// 	"Low": 0,
// 	"Open": 0,
// 	"OpenInterest": 0,
// 	"Status": {
// 	  "bit0": 0,
// 	  "bit1": 0,
// 	  "bit2": 0,
// 	  "bit3": 0,
// 	  "bit4": 0,
// 	  "bit5": 0,
// 	  "bit6": 0,
// 	  "bit7": 0,
// 	  "bit8": 0,
// 	  "bit19": 0,
// 	  "bit23": 0,
// 	  "bit24": 0,
// 	  "bit25": 0,
// 	  "bit26": 0,
// 	  "bit27": 0,
// 	  "bit28": 0,
// 	  "bit29": 0
// 	},
// 	"TimeStamp": "string",
// 	"TotalTicks": 0,
// 	"TotalVolume": 0,
// 	"UnchangedTicks": 0,
// 	"UnchangedVolume": 0,
// 	"UpTicks": 0,
// 	"UpVolume": 0
// }
type DailyStartingOnModel struct {
	Close           float32
	DownTicks       uint64
	DownVolume      uint64
	High            float32
	Low             float32
	OpenInterest    uint64
	TimeStamp       string
	TotalTicks      uint64
	TotalVolume     uint64
	UnchangedTicks  uint64
	UnchangedVolume uint64
	UpTicks         uint64
	UpVolume        uint64
}

// /stream/barchart/{symbol}/{interval}/{unit}/{startDate}
const dailyStartingOnPathFormat = "/stream/barchart/%v/1/Daily/%v?%v"

func (r *tradestation) DailyStartingOn(symbol string, start time.Time) error {
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return fmt.Errorf("Argument empty: %v", "symbol")
	}
	err := r.AccessTokenOK()
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("access_token", r.AccessToken)
	data.Set("APIVersion", apiVersion)

	path := fmt.Sprintf(dailyStartingOnPathFormat, symbol, start.Format(dateFormat), data.Encode())
	req, err := http.NewRequest(http.MethodGet, r.baseURL+path, nil)
	if err != nil {
		return err
	}
	resp, err := r.RoundTripper.RoundTrip(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrStatusNotOK
	}

	dumper, err := r.DumperFactory.New("DailyStartingOn")
	if err != nil {
		return err
	}
	defer closeIfCloser(dumper)
	go io.Copy(dumpWriter{Dumper: dumper}, resp.Body)

	time.Sleep(time.Second * 10)

	return nil
}

type dumpWriter struct {
	general.Dumper
}

func (r dumpWriter) Write(p []byte) (int, error) {
	err := r.Dumper.Dump(p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func closeIfCloser(maybeCloser interface{}) {
	closer, ok := maybeCloser.(io.Closer)
	if ok {
		closer.Close()
	}
}
