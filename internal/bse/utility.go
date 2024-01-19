package bse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sanksons/stock-price-sdk/internal/common"
)

func GetBseUrlForStock(name string) string {
	code := GetStockBseCode(name)
	return fmt.Sprintf(BSE_URL, code)
}

func CreateRequest(name string) (*http.Request, error) {
	req, err := http.NewRequest("GET", GetBseUrlForStock(name), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", common.USER_AGENT)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Referer", "https://www.bseindia.com/")
	return req, nil
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

	r := map[string]string{}
	if err := json.Unmarshal(b, &r); err != nil {
		return ret, err
	}

	// parse date
	dateStr, exists := r["CurrDate"]
	if !exists {
		return ret, fmt.Errorf("index `CurrDate` does not exists in response.")
	}
	date, err := time.Parse("Mon Jan 02 2006 15:04:05", dateStr)
	if err != nil {
		return ret, fmt.Errorf("could not parse date, error: %s", err.Error())
	}
	ret.CurrDate = date

	// parse nav
	navStr, exists := r["CurrVal"]
	if !exists {
		return ret, fmt.Errorf("index `CurrVal` does not exists in response.")
	}
	nav, err := strconv.ParseFloat(navStr, 64)
	if err != nil {
		return ret, fmt.Errorf("could not parse CurrVal, error: %s", err.Error())
	}
	ret.CurrVal = nav
	return ret, nil
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
	req, err := CreateRequest(name)
	if err != nil {
		return nil, err
	}
	resp, err := common.FireRequest(req)
	if err != nil {
		return nil, err
	}
	presp, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}
	return &stockNavImpl{
		nav:  presp.CurrVal,
		date: presp.CurrDate,
	}, nil
}
