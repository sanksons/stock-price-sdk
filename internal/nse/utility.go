package nse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sanksons/stock-price-sdk/internal/common"
)

const TYPE_FIRST_REQUEST = "first"
const TYPE_SECOND_REQUEST = "second"

func GetNseUrlForStock(name string, reqno string) string {

	code := GetStockNseCode(name)
	switch reqno {
	case TYPE_FIRST_REQUEST:
		return fmt.Sprintf(NSE1_URL, code)

	case TYPE_SECOND_REQUEST:
		return fmt.Sprintf(NSE2_URL, code)
	}
	return fmt.Sprintf(NSE1_URL, code)
}

func CreateFirstRequest(name string) (*http.Request, error) {
	req, err := http.NewRequest("GET", GetNseUrlForStock(name, TYPE_FIRST_REQUEST), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", common.USER_AGENT)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Referer", "https://www.nseindia.com/")
	return req, nil
}

func CreateSecondRequest(name string, cookies []*http.Cookie, referer string) (*http.Request, error) {
	req, err := http.NewRequest("GET", GetNseUrlForStock(name, TYPE_SECOND_REQUEST), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", common.USER_AGENT)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Referer", referer)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	return req, nil
}

type stockNavImpl struct {
	nav  float64
	date time.Time
}

func (this *stockNavImpl) GetNav() float64 {
	return this.nav
}
func (this *stockNavImpl) GetNavDate() time.Time {
	return this.date
}

func GetStockNav(name string) (common.StockNav, error) {
	req, err := CreateFirstRequest(name)
	if err != nil {
		return nil, err
	}

	resp, err := common.FireRequest(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	cookies := resp.Cookies()

	req2, err := CreateSecondRequest(name, cookies, req.RequestURI)
	if err != nil {
		return nil, err
	}
	resp2, err := common.FireRequest(req2)
	if err != nil {
		return nil, err
	}
	ret, err := ParseResponse(resp2)
	if err != nil {
		return nil, err
	}
	return &stockNavImpl{
		nav:  ret.CurrVal,
		date: ret.CurrDate,
	}, nil
}

type Response struct {
	CurrDate time.Time
	CurrVal  float64
}

func ParseResponse(resp *http.Response) (Response, error) {

	defer resp.Body.Close()

	var ret = Response{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}

	type Result struct {
		Metadata  map[string]string
		PriceInfo struct {
			LastPrice float64
		}
	}
	r := Result{}
	if err := json.Unmarshal(b, &r); err != nil {
		return ret, err
	}

	// parse date
	dateStr, exists := r.Metadata["lastUpdateTime"]
	if !exists {
		return ret, fmt.Errorf("index `lastUpdateTime` does not exists in response.")
	}

	// 02-Jan-2006
	date, err := time.Parse("02-Jan-2006 15:04:05", dateStr)
	if err != nil {
		return ret, fmt.Errorf("could not parse date, error: %s", err.Error())
	}
	ret.CurrDate = date

	// parse nav
	ret.CurrVal = r.PriceInfo.LastPrice

	return ret, nil
}
