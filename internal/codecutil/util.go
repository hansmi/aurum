package codecutil

import "reflect"

type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

func Unmarshal(c Codec, data []byte, valueType reflect.Type) (any, error) {
	dest := reflect.New(valueType.Elem())

	if err := c.Unmarshal(data, dest.Interface()); err != nil {
		return nil, err
	}

	return dest.Interface(), nil
}

// Internally values are always pointers to a non-pointer type.
func NormalizeValue(value any) (any, reflect.Type) {
	v := reflect.ValueOf(value)

	if !v.IsValid() {
		return nil, nil
	}

	if v.Kind() == reflect.Pointer {
		// Find last pointer in chain (if any).
		for v.Kind() == reflect.Pointer && v.Elem().Kind() == reflect.Pointer {
			v = v.Elem()
		}
	} else {
		// Get pointer to a copy
		v = reflect.New(v.Type())
		v.Elem().Set(reflect.ValueOf(value))
	}

	return v.Interface(), v.Type()
}
