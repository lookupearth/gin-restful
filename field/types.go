package field

import "errors"

type IgnoreError struct{}

func (ie IgnoreError) Error() string {
	return "ignore"
}

var null = []byte("null")

var errEmptyInput = errors.New("empty input for UnmarshalJSON")
