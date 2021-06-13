package functions

import (
    "fmt"
    "net/http"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "time"
    "net/url"
)


func listBalance(w http.ResponseWriter, r *http.Request) {
	const host = "https://api.binance.com"
  	const path = "/api/v3/account"

    const binance_api_key = "QsoqKGt4xku8xoUCqSZG7YYWeOqjjyUSIdbbLQlGwpZQSvJK6m9bqNkSnOlgbHvj";
    const binance_api_secret = "m4h1k5jNDregkjqa0lW5KmoOKFhK0Hklm6cVJeLNO7B7OcmSSQVnhfgBQmj59O4R";

    params := url.Values{}
    params.Set("timestamp", time.Now().String())
    params.Set("signature", hmacSha256(params.Encode(), binance_api_secret))

    result, err := http.Get(host+path+"?"+params.Encode())
    if err != nil { fmt.Println(err) }
    fmt.Fprint(w, result)
}

func hmacSha256(data string, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}