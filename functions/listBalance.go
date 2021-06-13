package functions

import (
	"fmt"
	"net/http"
	"time"
)

func ListBalance(w http.ResponseWriter, r *http.Request) {
	const url = "https://api.binance.com/api/v3/account"
	params := map[string]string{
		"timestamp":  fmt.Sprint(time.Now().UnixNano() / int64(time.Millisecond)),
		"recvWindow": "1000",
	}

	res, err := doGet(url, params)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, string(res))
}
