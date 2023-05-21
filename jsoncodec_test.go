package aurum

import (
	"bytes"
	"testing"

	"github.com/hansmi/aurum/internal/codectest"
	"github.com/hansmi/aurum/internal/ref"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestJSONCodec(t *testing.T) {
	tests := []codectest.Case{
		{
			Name:  "string",
			Value: "hello world",
		},
		{
			Name:  "string pointer",
			Value: ref.Ref("hello world"),
		},
		{
			Name:  "int slice",
			Value: []int{0, 1, 2, 3},
		},
		{
			Name:  "int slice pointer",
			Value: ref.Ref([]int{0, 1, 2, 3}),
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
			Name:  "empty struct",
			Value: struct{}{},
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

	codectest.AssertAll(t, &JSONCodec{}, tests)
}

func TestJSONCodecMarshalEndsWithNewline(t *testing.T) {
	for _, value := range []any{0, "", &wrapperspb.Int64Value{Value: 4321}} {
		var c JSONCodec

		data, err := c.Marshal(value)
		if err != nil {
			t.Errorf("Marshal(%v) failed: %v", value, err)
		}

		if !(len(data) == 0 || bytes.HasSuffix(data, []byte{'\n'})) {
			t.Errorf("Marshal(%v) return value does not end in newline: %q", value, data)
		}
	}
}
