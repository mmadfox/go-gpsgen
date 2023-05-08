package types

import (
	"errors"

	"github.com/mmadfox/go-gpsgen/random"
)

const (
	maxModelLen = 64
	minModelLen = 1
)

var (
	ErrModelTooShort = errors.New("types/model: model value too short")
	ErrModelTooLong  = errors.New("types/model: model value too long")
)

type Model struct {
	val string
}

func (m Model) String() string {
	return m.val
}

func NewModel(value string) (Model, error) {
	if len(value) < minModelLen {
		return Model{}, ErrModelTooShort
	}
	if len(value) > maxModelLen {
		return Model{}, ErrModelTooLong
	}
	return Model{val: value}, nil
}

func RandomModel() Model {
	return Model{val: "RT-" + random.String(6)}
}

func (m Model) IsEmpty() bool {
	return len(m.val) == 0
}
