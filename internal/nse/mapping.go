package nse

var isin2nseCode map[string]string = map[string]string{
	"INF204KC1402": "SILVERBEES", // silver bees
	"INF204KB15V2": "ITBEES",     // it bees
	"INF204KB17I5": "GOLDBEES",   // gold bees
}

func GetStockNseCode(isin string) string {
	return isin2nseCode[isin]
}
