package functions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hcchang0701/portfolio-tracker/src/utils"
)

type Account struct {
	Balances []struct {
		Asset  string `json:"asset"`
		Free   string `json:"free"`
		Locked string `json:"locked"`
	} `json:"balances"`
}

func ListBalance(w http.ResponseWriter, r *http.Request) {
	const url = "https://api.binance.com/api/v3/account"
	params := map[string]string{
		"timestamp":  fmt.Sprint(time.Now().UnixNano() / int64(time.Millisecond)),
		"recvWindow": "1000",
	}

	res, err := utils.DoGet(url, params)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, string(res))
}
