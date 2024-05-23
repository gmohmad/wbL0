package orders

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validate(data []byte) (*OrderItem, error) {
	var ordItem OrderItem

	err := json.Unmarshal(data, &ordItem)

	if err != nil {
		return nil, fmt.Errorf("Order data decoding failed: %w", err)
	}

	validator := validator.New(validator.WithRequiredStructEnabled())

	err = validator.Struct(ordItem)

	if err != nil {
		return nil, fmt.Errorf("Order data validation failed: %w", err)
	}

	return &ordItem, nil
}
