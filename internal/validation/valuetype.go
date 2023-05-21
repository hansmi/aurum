package validation

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
)

var protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()

func CheckValueType(v any) error {
	if v == nil {
		return multierr.Append(os.ErrInvalid, errors.New("value must not be nil"))
	}

	if valueType := reflect.TypeOf(v); valueType.Kind() == reflect.Struct && reflect.PointerTo(valueType).Implements(protoMessageType) {
		return multierr.Append(os.ErrInvalid, fmt.Errorf("protobuf messages must be pointers, got %s", valueType.String()))
	}

	return nil
}
