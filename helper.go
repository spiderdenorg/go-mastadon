package masta

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func addParamString(params *url.Values) func(key, value string) {
	return func(key, value string) {
		if value != "" {
			params.Add(key, value)
		}
	}
}

func addParamBool(params *url.Values) func(key string, value bool) {
	return func(key string, value bool) {
		if value {
			params.Add(key, "true")
		}
	}
}

func addParamFuncs(params *url.Values) (func(key, value string), func(key string, value bool)) {
	return addParamString(params), addParamBool(params)
}

// Base64EncodeFileName returns the base64 data URI format string of the file with the file name.
func Base64EncodeFileName(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return Base64Encode(file)
}

// Base64Encode returns the base64 data URI format string of the file.
func Base64Encode(file *os.File) (string, error) {
	fi, err := file.Stat()
	if err != nil {
		return "", err
	}

	d := make([]byte, fi.Size())
	_, err = file.Read(d)
	if err != nil {
		return "", err
	}

	return "data:" + http.DetectContentType(d) +
		";base64," + base64.StdEncoding.EncodeToString(d), nil
}

// String is a helper function to get the pointer value of a string.
func String(v string) *string { return &v }

func parseAPIError(prefix string, resp *http.Response) error {
	errMsg := fmt.Sprintf("%s: %s", prefix, resp.Status)
	var e struct {
		Error string `json:"error"`
	}

	json.NewDecoder(resp.Body).Decode(&e)
	if e.Error != "" {
		errMsg = fmt.Sprintf("%s: %s", errMsg, e.Error)
	}

	return &APIError{
		resp.StatusCode,
		errors.New(errMsg),
	}
}
