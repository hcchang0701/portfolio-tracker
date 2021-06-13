package functions

import (
    "fmt"
    "net/http"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "time"
    "io/ioutil"
)


func ListBalance(w http.ResponseWriter, r *http.Request) {
	const url = "https://api.binance.com/api/v3/account"
    const binance_api_key = "QsoqKGt4xku8xoUCqSZG7YYWeOqjjyUSIdbbLQlGwpZQSvJK6m9bqNkSnOlgbHvj";
    const binance_api_secret = "m4h1k5jNDregkjqa0lW5KmoOKFhK0Hklm6cVJeLNO7B7OcmSSQVnhfgBQmj59O4R";

    client := &http.Client{}
    req, _ := http.NewRequest(http.MethodGet, url, nil)
    req.Header.Set("X-MBX-APIKEY", binance_api_key)

    params := req.URL.Query()
    params.Add("timestamp", fmt.Sprint(time.Now().UnixNano()/int64(time.Millisecond)))
    params.Add("recvWindow", "5000")
    params.Add("signature", hmacSha256(params.Encode(), binance_api_secret))
    req.URL.RawQuery = params.Encode()

    fmt.Println(req.URL)
    res, _ := client.Do(req)
    defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

    fmt.Fprint(w, string(body))
}

func hmacSha256(data string, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}