package types

import (
	"errors"

	"github.com/mmadfox/go-gpsgen/random"
)

const (
	MaxModelLen = 64
	MinModelLen = 3
)

var (
	// ErrModelTooShort represents an error condition where the model value is too short.
	ErrModelTooShort = errors.New("types/model: model value too short")
	// ErrModelTooLong represents an error condition where the model value is too long.
	ErrModelTooLong = errors.New("types/model: model value too long")
)

// Model represents the tracker device model.
type Model struct {
	val string
}

// String returns the string representation of the model.
func (m Model) String() string {
	return m.val
}

// NewModel creates a new Model instance with the provided value.
// It performs length validation on the value and returns an error
// if it's too short or too long.
func NewModel(value string) (Model, error) {
	if len(value) < MinModelLen {
		return Model{}, ErrModelTooShort
	}
	if len(value) > MaxModelLen {
		return Model{}, ErrModelTooLong
	}
	return Model{val: value}, nil
}

// RandomModel generates and returns a random Model instance.
// It creates a Model object with a value in the format "RT-" followed by a
// randomly generated string of length 6.
func RandomModel() Model {
	return Model{val: "RT-" + random.String(6)}
}

// IsEmtpy checks whether the model value is empty (i.e., has a length of 0).
// It returns true if the value is empty, and false otherwise.
func (m Model) IsEmpty() bool {
	return len(m.val) == 0
}
