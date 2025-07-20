package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/meh-hackathon/meh/apperror"
)

var (
	ErrParseBody = apperror.Define("server:parse_body", "Failed to parse request body").WithStatus(http.StatusBadRequest)
)

func ParseReqBody[T any](r *http.Request) (T, error) {
	var zeroValue T

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return zeroValue, ErrParseBody.WithCause(err).WithApiMessage("failed to read body")
	}
	defer func() { writeReqBodyBack(r, body) }()

	return ParseData[T](body)
}

func ParseRespBody[T any](resp *http.Response) (T, error) {
	var zeroValue T

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return zeroValue, ErrParseBody.WithCause(err).WithApiMessage("failed to read response body")
	}

	defer func() { writeRespBodyBack(resp, body) }()

	return ParseData[T](body)
}

func ParseData[T any](data []byte) (T, error) {
	var zeroValue T

	if len(data) == 0 {
		return zeroValue, ErrParseBody.WithMessage("body is empty")
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&zeroValue)
	if err != nil {
		return zeroValue, ErrParseBody.WithCause(err).WithApiMessage(err.Error())
	}

	err = checkFieldsAreSpecified(zeroValue, data)
	if err != nil {
		return zeroValue, ErrParseBody.WithCause(err).WithApiMessage(err.Error())
	}

	return zeroValue, nil
}

func checkFieldsAreSpecified[T any](v T, data []byte) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var jsonData map[string]json.RawMessage
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("failed to parse json: %w", err)
	}

	return validateStructFields(val, jsonData)
}

func validateStructFields(val reflect.Value, jsonData map[string]json.RawMessage) error {
	if !val.IsValid() {
		return fmt.Errorf("invalid value provided")
	}

	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return fmt.Errorf("nil pointer or interface provided")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %s", val.Kind())
	}

	if jsonData == nil {
		return fmt.Errorf("jsonData cannot be nil")
	}

	t := val.Type()
	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		if fieldType.Anonymous {
			// Recursively validate embedded structs
			if err := validateStructFields(field, jsonData); err != nil {
				return err
			}
			continue
		}

		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		tagParts := strings.Split(jsonTag, ",")
		if len(tagParts) == 0 {
			continue
		}

		// Check for omitempty option
		if slices.Contains(tagParts, "omitempty") {
			continue
		}

		fieldName := strings.TrimSpace(tagParts[0])
		if fieldName == "" {
			continue
		}

		_, specified := jsonData[fieldName]
		if !specified {
			return fmt.Errorf("required field %s is not set", fieldName)
		}
	}

	return nil
}

func writeReqBodyBack(r *http.Request, body []byte) {
	if r == nil || r.Body == nil {
		return
	}
	r.Body.Close() // Close original body

	r.Body = io.NopCloser(bytes.NewReader(body))
	r.ContentLength = int64(len(body))
	r.Header.Set("Content-Length", strconv.Itoa(len(body)))
}

func writeRespBodyBack(resp *http.Response, body []byte) {
	if resp == nil || resp.Body == nil {
		return
	}
	resp.Body.Close() // Always close the original body

	resp.Body = io.NopCloser(bytes.NewReader(body))
	resp.ContentLength = int64(len(body))
	resp.Header.Set("Content-Length", strconv.Itoa(len(body)))
}
