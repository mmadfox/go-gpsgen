package properties

// Properties represents a collection of key-value pairs with flexible value types.
type Properties map[string]interface{}

// Make creates and returns an empty Properties instance.
func Make() Properties {
	return make(Properties)
}

// Merge merges the key-value pairs from the other Properties into the current Properties.
func (p Properties) Merge(other Properties) {
	for k, v := range other {
		p[k] = v
	}
}

// String retrieves the value associated with the given key as a string.
func (p Properties) String(key string) (string, bool) {
	s, ok := p[key].(string)
	if ok {
		return s, ok
	}
	return "", false
}

// Bool retrieves the value associated with the given key as a boolean.
func (p Properties) Bool(key string) (bool, bool) {
	b, ok := p[key].(bool)
	if ok {
		return b, ok
	}
	return false, false
}

// Bool retrieves the value associated with the given key as a boolean.
func (p Properties) Int(key string) (int, bool) {
	switch val := p[key].(type) {
	case int:
		return val, true
	case int32:
		return int(val), true
	case int64:
		return int(val), true
	default:
		return 0, false
	}
}

// Float64 retrieves the value associated with the given key as a float64.
func (p Properties) Float64(key string) (float64, bool) {
	i, ok := p[key].(float64)
	if ok {
		return i, ok
	}
	return 0, false
}

// Set sets the value associated with the given key.
func (p Properties) Set(key string, value interface{}) Properties {
	p[key] = value
	return p
}

// Remove removes the key-value pair associated with the given key.
func (p Properties) Remove(key string) {
	delete(p, key)
}

// Reset clears all key-value pairs in the Properties.
func (p *Properties) Reset() {
	*p = make(Properties)
}
