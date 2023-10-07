package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidPriceFormat = errors.New("invalid price format")

type Price float64

func (r Price) MarshalJSON() ([]byte, error) {

	jsonValue := fmt.Sprintf("%d price", r)

	quotedJSONValue := strconv.Quote(jsonValue)

	// Convert the quoted string value to a byte slice and return it.
	return []byte(quotedJSONValue), nil
}

func (r *Price) UnmarshalJSON(jsonValue []byte) error {

	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidPriceFormat
	}
	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedJSONValue, " ")
	// Sanity check the parts of the string to make sure it was in the expected format. // If it isn't, we return the ErrInvalidRuntimeFormat error again.
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidPriceFormat
	}
	// Otherwise, parse the string containing the number into an int32. Again, if this // fails return the ErrInvalidRuntimeFormat error.
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidPriceFormat
	}
	// Convert the int32 to a Runtime type and assign this to the receiver. Note that we // use the * operator to deference the receiver (which is a pointer to a Runtime
	// type) in order to set the underlying value of the pointer.
	*r = Price(i)
	return nil
}
