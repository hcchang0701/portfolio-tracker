package functions

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"fmt"
	"github.com/hcchang0701/portfolio-tracker/src/utils"
	"net/http"
	//"sort"
	"time"
	"unicode"
	//"google.golang.org/api/option"
)

type Tx struct {
	Symbol          string `json:"symbol" bigquery:"symbol"`
	Id              int64  `json:"id" bigquery:"id"`
	OrderId         int64  `json:"orderId" bigquery:"orderId"`
	OrderListId     int64  `json:"orderListId" bigquery:"orderListId"`
	Price           string `json:"price" bigquery:"price"`
	Qty             string `json:"qty" bigquery:"qty"`
	QuoteQty        string `json:"quoteQty" bigquery:"quoteQty"`
	Commission      string `json:"commission" bigquery:"commission"`
	CommissionAsset string `json:"commissionAsset" bigquery:"commissionAsset"`
	Time            int64  `json:"time" bigquery:"time"`
	IsBuyer         bool   `json:"isBuyer" bigquery:"isBuyer"`
	IsMaker         bool   `json:"isMaker" bigquery:"isMaker"`
	IsBestMatch     bool   `json:"isBestMatch" bigquery:"isBestMatch"`
}

const genesis = 1613865600000

func UpdateTxHistory(w http.ResponseWriter, r *http.Request) {

	symbol := r.URL.Query().Get("symbol");
	if !isUpperCase(symbol) {
		fmt.Fprint(w, "Error: symbol must be uppercase")
		return
	}

	ts, err := getLastUpdateTimestamp(symbol)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	res, err := getTxHistory(symbol, ts+1)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	if err := utils.InsertData(symbol, res); err != nil {
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, fmt.Sprintf("%d of %s records updated", len(res), symbol))
}

func isUpperCase(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func getLastUpdateTimestamp(symbol string) (int64, error) {

	if !utils.BqTableExists(symbol) {

		schema, err := bigquery.InferSchema(Tx{})
		if err != nil {
			return 0, err
		}
		if err := utils.CreateBqTable(symbol, schema); err != nil {
			return 0, err
		}
		return genesis, nil
	}

	res, err := utils.QueryData(`
   		SELECT time
		FROM ` + fmt.Sprintf("`%s.%s`", utils.DatasetID, symbol) + `
		ORDER BY time DESC
		LIMIT 1
	`)
	if err != nil {
		return 0, err
	}

	return res[0].(int64), nil
}

func getTxHistory(symbol string, startTime int64) ([]Tx, error) {
	const url = "https://api.binance.com/api/v3/myTrades"
	params := map[string]string{
		"symbol":    symbol,
		"startTime": fmt.Sprint(startTime),
		"limit":     "1000",
		"timestamp": fmt.Sprint(time.Now().UnixNano() / int64(time.Millisecond)),
	}

	res, err := utils.DoGet(url, params)
	if err != nil {
		return nil, err
	}

	txs := []Tx{}
	if err := json.Unmarshal(res, &txs); err != nil {
		return nil, err
	}

	// sort descendingly
	// sort.Slice(txs, func(i, j int) bool {
	// 	return txs[i].Time >= txs[j].Time
	// })
	return txs, nil
}
