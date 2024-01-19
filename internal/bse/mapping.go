package bse

var isin2bseCode map[string]string = map[string]string{
	"INF200KA1FS1": "590138", // SBI Nifty 50
	"INF109K012R6": "537007", // ICICI Nifty 50
	"INF204KB14I2": "590103", // NIPPON Nifty bees
	"INF179KC1965": "539516", // HDFC Nifty 50
	"INF769K01HF4": "543291", // NYSE FANG
	"INF109KC1NT3": "533244", // ICICI GOLD ETF
	"INF247L01AP3": "533385", // NASDAQ-100
}

func GetStockBseCode(isin string) string {
	return isin2bseCode[isin]
}
