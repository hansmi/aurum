package aurum

import (
	"bytes"
	"os"
	"testing"

	"github.com/hansmi/aurum/internal/codectest"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestTextProtoCodec(t *testing.T) {
	tests := []codectest.Case{
		{
			Name:           "string",
			Value:          "hello world",
			WantMarshalErr: os.ErrInvalid,
		},
		{
			Name:           "empty struct",
			Value:          struct{}{},
			WantMarshalErr: os.ErrInvalid,
		},
		{
			Name:  "proto int64value",
			Value: wrapperspb.Int64Value{Value: 1},
		},
		{
			Name:  "proto int64value pointer",
			Value: &wrapperspb.Int64Value{Value: 4321},
		},
		{
			Name: "proto structpb",
			Value: func() *structpb.Struct {
				s, err := structpb.NewStruct(map[string]any{
					"hello": true,
					"world": "false str",
				})
				if err != nil {
					t.Fatal(err)
				}
				return s
			}(),
		},
	}

	codectest.AssertAll(t, &TextProtoCodec{}, tests)
}

func TestTextProtoCodecMarshalEndsWithNewline(t *testing.T) {
	for _, value := range []any{&wrapperspb.Int64Value{Value: 4321}} {
		var c TextProtoCodec

		data, err := c.Marshal(value)
		if err != nil {
			t.Errorf("Marshal(%v) failed: %v", value, err)
		}

		if !(len(data) == 0 || bytes.HasSuffix(data, []byte{'\n'})) {
			t.Errorf("Marshal(%v) return value does not end in newline: %q", value, data)
		}
	}
}
