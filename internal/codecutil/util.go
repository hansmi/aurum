package codecutil

import (
	"fmt"
	"os"
	"reflect"

	"google.golang.org/protobuf/proto"
)

var rvZero reflect.Value
var protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()

type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

// Unmarshal decodes raw data using the given codec and returns a newly
// constructed value.
func Unmarshal(c Codec, data []byte, valueType reflect.Type) (any, error) {
	// The destination is a pointer (passed to Unmarshal()) to a pointer of the
	// underlying value type. This way JSON's "null" is unmarshalled as nil.
	dest := reflect.New(valueType)
	dest.Elem().Set(reflect.New(valueType.Elem()))

	if err := c.Unmarshal(data, dest.Interface()); err != nil {
		return nil, err
	}

	return dest.Elem().Interface(), nil
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

func checkProtoType(rv reflect.Value) error {
	if vt := rv.Type(); vt.Kind() == reflect.Struct && reflect.PointerTo(vt).Implements(protoMessageType) {
		return fmt.Errorf("%w: protobuf messages must be pointers, got %s", os.ErrInvalid, vt.String())
	}

	return nil
}

func CheckValueType(v any) error {
	if rv := reflect.ValueOf(v); rv.IsValid() {
		return checkProtoType(rv)
	}

	return nil
}

// Validate the input value for marshalling.
func PrepareMarshalValue(v any) (reflect.Value, proto.Message, error) {
	rv := reflect.ValueOf(v)

	if !rv.IsValid() {
		return rvZero, nil, fmt.Errorf("%w: invalid input value %#v", os.ErrInvalid, v)
	} else if err := checkProtoType(rv); err != nil {
		return rvZero, nil, err
	}

	m, _ := rv.Interface().(proto.Message)

	return rv, m, nil
}

// Validate the output destination for unmarshalling. It must be a non-nil
// pointer.
func PrepareUnmarshalDest(v any) (reflect.Value, proto.Message, error) {
	rv := reflect.ValueOf(v)

	if rv.IsNil() || rv.Kind() != reflect.Pointer {
		return rvZero, nil, fmt.Errorf("%w: v must be a non-nil pointer, got %#v", os.ErrInvalid, rv)
	} else if err := checkProtoType(rv.Elem()); err != nil {
		return rvZero, nil, err
	}

	m, _ := rv.Elem().Interface().(proto.Message)

	return rv, m, nil
}
