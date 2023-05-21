package aurum

import (
	"encoding"
	"fmt"
	"os"
	"reflect"
)

// TextCodec stores values in plain-text files. Supports strings, byte slices,
// rune slices and values implementing [encoding.TextUnmarshaler].
type TextCodec struct {
	// Codec to use for unsupported types.
	Fallback Codec
}

var _ Codec = (*TextCodec)(nil)

func (t TextCodec) Marshal(v any) ([]byte, error) {
	if v, ok := v.(encoding.TextMarshaler); ok {
		return v.MarshalText()
	}

	if v := reflect.Indirect(reflect.ValueOf(v)).Interface(); v != nil {
		switch v := v.(type) {
		case []byte:
			return v, nil
		case string:
			return []byte(v), nil
		case []rune:
			return []byte(string(v)), nil
		}
	}

	if t.Fallback == nil {
		return nil, fmt.Errorf("%w: marshalling %T as text is not supported", os.ErrInvalid, v)
	}

	return t.Fallback.Marshal(v)
}

func (t TextCodec) Unmarshal(data []byte, v any) error {
	switch v := v.(type) {
	case *[]byte:
		*v = append([]byte(nil), data...)
		return nil
	case *string:
		*v = string(data)
		return nil
	case *[]rune:
		*v = []rune(string(data))
		return nil
	case encoding.TextUnmarshaler:
		return v.UnmarshalText(data)
	}

	if t.Fallback == nil {
		return fmt.Errorf("%w: unmarshalling text into %T is not supported", os.ErrInvalid, v)
	}

	return t.Fallback.Unmarshal(data, v)
}
