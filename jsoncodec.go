package aurum

import (
	"bytes"
	"encoding/json"

	"github.com/hansmi/aurum/internal/validation"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JSONCodec stores values using the JSON format.
//
// Protocol buffer messages are detected and marshalled using [protojson]
// before re-formatting the resulting JSON data. Protojson produces unstable
// output by design.
type JSONCodec struct {
	ProtoMarshalOptions   protojson.MarshalOptions
	ProtoUnmarshalOptions protojson.UnmarshalOptions
}

var _ Codec = (*JSONCodec)(nil)

func (c *JSONCodec) Marshal(v any) ([]byte, error) {
	if err := validation.CheckValueType(v); err != nil {
		return nil, err
	}

	var err error
	var data []byte

	if m, ok := v.(proto.Message); ok {
		data, err = c.ProtoMarshalOptions.Marshal(m)
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := json.Indent(&buf, data, "", "  "); err != nil {
		return nil, err
	}

	if buf.Len() > 0 {
		// json.Indent doesn't write a terminating newline.
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

func (c *JSONCodec) Unmarshal(data []byte, v any) error {
	if err := validation.CheckValueType(v); err != nil {
		return err
	}

	if m, ok := v.(proto.Message); ok {
		return c.ProtoUnmarshalOptions.Unmarshal(data, m)
	}

	return json.Unmarshal(data, v)
}
