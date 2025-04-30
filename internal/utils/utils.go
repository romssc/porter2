package utils

import (
	"encoding/json"
	"io"
)

/*

This file provides utility functions for handling JSON data and extracting error messages
from HTTP response bodies.

These functions provide basic functionality for working with JSON data and error handling in the context
of HTTP interactions or other JSON-based workflows.

*/

// MarshalJSON() marshals an object into its JSON representation as a byte slice.
func MarshalJSON(v any) []byte {
	m, _ := json.Marshal(v)

	return m
}

// ExtractError() attempts to extract the error message from an HTTP response body.
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
