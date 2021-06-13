package functions

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/url"
)

const genesis = 1613865600000
const binance_api_key = "QsoqKGt4xku8xoUCqSZG7YYWeOqjjyUSIdbbLQlGwpZQSvJK6m9bqNkSnOlgbHvj"
const binance_api_secret = "m4h1k5jNDregkjqa0lW5KmoOKFhK0Hklm6cVJeLNO7B7OcmSSQVnhfgBQmj59O4R"

var client = &http.Client{}

func hmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func addSignature(query url.Values) {
	query.Add("signature", hmacSha256(query.Encode(), binance_api_secret))
}

func addAPIKey(header http.Header) {
	header.Set("X-MBX-APIKEY", binance_api_key)
}

func doGet(url string, params map[string]string) ([]byte, error) {

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	addAPIKey(req.Header)

	query := req.URL.Query()
	for k, v := range params {
		query.Add(k, v)
	}

	addSignature(query)
	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
