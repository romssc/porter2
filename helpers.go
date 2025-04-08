package elasticmigrator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// buildUrl() constructs the URL for an Elasticsearch API request.
//
// It includes the index name and an optional API endpoint for specific actions (e.g., _bulk, _delete_by_query).
// The 'pretty' query parameter is added to make the response more readable.
func buildUrl(addr string, index string, api string) string {
	switch {
	case api != "":
		return fmt.Sprintf("%v/%v/%v?pretty", addr, index, api)

	default:
		return fmt.Sprintf("%v/%v?pretty", addr, index)
	}
}

// buildRequest() creates a new HTTP request with the specified method, username, password, URL, and request body.
//
// It also sets the appropriate headers for basic authentication and JSON content type.
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

// marshal() serializes the provided data (map) into JSON format.
// This function is used to convert Go data structures into JSON before sending to Elasticsearch.
func marshal(data map[string]interface{}) []byte {
	dataByte, _ := json.Marshal(data)

	return dataByte
}

// decode() deserializes the JSON response body into a Go map (map[string]interface{}).
//
// It takes the response body (io.ReadCloser) and converts it into a structured map.
func decode(r io.ReadCloser) (map[string]interface{}, error) {
	var d map[string]interface{}

	err := json.NewDecoder(r).Decode(&d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
