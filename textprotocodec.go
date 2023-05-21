package aurum

import (
	"fmt"
	"os"

	"github.com/hansmi/aurum/internal/validation"
	"github.com/protocolbuffers/txtpbfmt/parser"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// TextProtoCodec stores values using the textproto format. Only protocol
// buffer messages are supported.
//
// [txtpbfmt] is used to format the resulting data as [prototext] produces
// unstable output by design.
type TextProtoCodec struct {
	ProtoMarshalOptions   prototext.MarshalOptions
	ProtoUnmarshalOptions prototext.UnmarshalOptions
}

var _ Codec = (*TextProtoCodec)(nil)

func toProtoMessage(v any) (proto.Message, error) {
	if err := validation.CheckValueType(v); err != nil {
		return nil, err
	}

	m, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("%w: only protobuf messages are supported, got %T", os.ErrInvalid, v)
	}

	return m, nil
}

func (c *TextProtoCodec) Marshal(v any) ([]byte, error) {
	m, err := toProtoMessage(v)
	if err != nil {
		return nil, err
	}

	data, err := c.ProtoMarshalOptions.Marshal(m)
	if err != nil {
		return nil, err
	}

	return parser.FormatWithConfig(data, parser.Config{
		ExpandAllChildren:        true,
		SkipAllColons:            true,
		WrapStringsAfterNewlines: true,
	})
}

func (c *TextProtoCodec) Unmarshal(data []byte, v any) error {
	m, err := toProtoMessage(v)
	if err != nil {
		return err
	}

	return c.ProtoUnmarshalOptions.Unmarshal(data, m)
}
