package utils

import (
	"encoding/json"
	"io"
	"os"
)

type MapperFunc func() map[string]interface{}

func GetMap(funcs []MapperFunc) map[string]interface{} {
	result := map[string]interface{}{}

	for _, fn := range funcs {
		if fn == nil {
			continue
		}
		for k, v := range fn() {
			result[k] = v
		}
	}

	return result
}

func GetContents(path string) ([]byte, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func MarshalJSON(v any) []byte {
	m, _ := json.Marshal(v)

	return m
}

func ExtractError(body io.ReadCloser) (string, bool) {
	var r map[string]interface{}

	json.NewDecoder(body).Decode(&r)

	err, ok := r["error"].(map[string]interface{})
	if !ok {
		return "", false
	}

	reason, _ := err["reason"].(string)

	return reason, true
}

func ExtractBulkErrors(body io.ReadCloser) ([]string, bool) {
	var r map[string]interface{}

	json.NewDecoder(body).Decode(&r)

	ok, _ := r["errors"].(bool)
	if !ok {
		return nil, false
	}

	items, ok := r["items"].([]interface{})
	if !ok {
		return nil, true
	}

	var errors []string

	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		for _, v := range m {
			doc, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			err, ok := doc["error"].(map[string]interface{})
			if ok {
				errors = append(errors, err["reason"].(string))
			}
		}
	}

	return errors, true
}
