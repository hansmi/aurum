package aurum

import (
	"bytes"
	"encoding/json"

	"github.com/hansmi/aurum/internal/codecutil"
	"google.golang.org/protobuf/encoding/protojson"
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
	rv, m, err := codecutil.PrepareMarshalValue(v)
	if err != nil {
		return nil, err
	}

	var data []byte

	if m != nil {
		data, err = c.ProtoMarshalOptions.Marshal(m)
	} else {
		data, err = json.Marshal(rv.Interface())
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
	rv, m, err := codecutil.PrepareUnmarshalDest(v)
	if err != nil {
		return err
	}

	if m != nil {
		return c.ProtoUnmarshalOptions.Unmarshal(data, m)
	}

	return json.Unmarshal(data, rv.Interface())
}
