package pkg

import (
	"fmt"
	"time"

	"github.com/sanksons/stock-price-sdk/internal/bse"
	"github.com/sanksons/stock-price-sdk/internal/common"
	"github.com/sanksons/stock-price-sdk/internal/nse"
)

const TYPE_NSE = "nse"
const TYPE_BSE = "bse"

type StockNav struct {
	Nav  float64
	Date time.Time
}

func GetStockNav(name string) (*StockNav, error) {
	var codeType string = TYPE_BSE
	var code string = bse.GetStockBseCode(name)
	if code == "" {
		code = nse.GetStockNseCode(name)
		codeType = TYPE_NSE
	}

	if code == "" {
		return nil, fmt.Errorf("No Code found for stock")
	}

	var stockI common.StockNav
	var err error
	if codeType == TYPE_BSE {
		stockI, err = bse.GetStockNav(name)
	} else {
		stockI, err = nse.GetStockNav(name)
	}

	if err != nil {
		return nil, err
	}
	return &StockNav{Nav: stockI.GetNav(), Date: stockI.GetNavDate()}, nil

}
