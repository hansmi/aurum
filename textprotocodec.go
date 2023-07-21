package aurum

import (
	"fmt"
	"os"

	"github.com/hansmi/aurum/internal/codecutil"
	"github.com/protocolbuffers/txtpbfmt/parser"
	"google.golang.org/protobuf/encoding/prototext"
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

func (c *TextProtoCodec) Marshal(v any) ([]byte, error) {
	_, m, err := codecutil.PrepareMarshalValue(v)
	if err != nil {
		return nil, err
	}

	if m == nil {
		return nil, fmt.Errorf("%w: only protobuf messages can be marshalled, got %T", os.ErrInvalid, v)
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
	_, m, err := codecutil.PrepareUnmarshalDest(v)
	if err != nil {
		return err
	}

	if m == nil {
		return fmt.Errorf("%w: only protobuf messages can be unmarshalled, got %T", os.ErrInvalid, v)
	}

	return c.ProtoUnmarshalOptions.Unmarshal(data, m)
}
