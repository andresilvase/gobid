package jsonutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/andresilvase/gobid/internal/validator"
)

func EncodeJson[T any](w http.ResponseWriter, r *http.Request, statusCode int, data T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}

	return nil
}

func DecodeValidJson[T validator.Validator](r *http.Request) (T, map[string]string, error) {
	var data T

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError
		var errorMsg error

		if errors.As(err, &unmarshalTypeError) {
			errorMsg = fmt.Errorf("field: %s, expected: %s, got: %s", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value)
		} else {
			errorMsg = fmt.Errorf("invalid JSON format: %w", err)
		}

		return data, map[string]string{"error": errorMsg.Error()}, errorMsg
	}

	if problems := data.Valid(r.Context()); len(problems) > 0 {
		return data, problems, fmt.Errorf("invalid %T: %d problems", data, len(problems))
	}

	return data, nil, nil
}

func DecodeJson[T any](r *http.Request) (T, error) {
	var data T

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("failed to decode data: %w", err)
	}

	return data, nil
}
