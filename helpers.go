package goelasticmigrator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func buildUrl(addr string, index string, api string) string {
	switch {
	case api != "":
		return fmt.Sprintf("%v/%v/%v?pretty", addr, index, api)

	default:
		return fmt.Sprintf("%v/%v?pretty", addr, index)
	}
}

func buildRequest(method string, username string, password string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(
		username,
		password,
	)

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func marshal(data map[string]interface{}) []byte {
	dataByte, _ := json.Marshal(data)

	return dataByte
}

func decode(r io.ReadCloser) (map[string]interface{}, error) {
	var d map[string]interface{}

	err := json.NewDecoder(r).Decode(&d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
