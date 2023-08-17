package navigator

import "fmt"

const (
	MinNameValue = 3
	MaxNameValue = 256
)

// Name represents a name used in the navigator.
type Name struct {
	value string
}

// ParseName parses the given string and returns a Name instance.
// An error is returned if the name is invalid.
func ParseName(name string) (Name, error) {
	typ := Name{value: name}
	if err := typ.validate(); err != nil {
		return Name{}, err
	}
	return typ, nil
}

// IsEmpty checks if the Name instance is empty.
func (t Name) IsEmpty() bool {
	return len(t.value) == 0
}

// IsEmpty checks if the Name instance is empty.
func (t Name) String() string {
	return t.value
}

func (t Name) validate() error {
	if len(t.value) < MinNameValue {
		return fmt.Errorf("name too short")
	}
	if len(t.value) > MaxNameValue {
		return fmt.Errorf("name too long")
	}
	return nil
}
