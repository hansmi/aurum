package aurum

import (
	"os"
	"testing"
	"time"

	"github.com/hansmi/aurum/internal/codectest"
	"github.com/hansmi/aurum/internal/ref"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestTextCodec(t *testing.T) {
	genericTests := []codectest.Case{
		{
			Name:  "string",
			Value: "hello world",
		},
		{
			Name:  "string pointer",
			Value: ref.Ref("hello world"),
		},
		{
			Name:  "byte slice",
			Value: []byte("byte value"),
		},
		{
			Name:  "byte slice pointer",
			Value: ref.Ref([]byte("byte value")),
		},
		{
			Name:  "rune slice",
			Value: []rune("rune value"),
		},
		{
			Name:  "time",
			Value: time.Date(2000, time.January, 1, 0, 1, 2, 3, time.UTC),
		},
	}

	codectest.AssertAll(t, &TextCodec{}, append([]codectest.Case{
		{
			Name:           "empty struct",
			Value:          struct{}{},
			WantMarshalErr: os.ErrInvalid,
		},
		{
			Name:           "proto int64value",
			Value:          &wrapperspb.Int64Value{Value: 4321},
			WantMarshalErr: os.ErrInvalid,
		},
	}, genericTests...))

	t.Run("JSON fallback", func(t *testing.T) {
		codectest.AssertAll(t, &TextCodec{
			Fallback: &JSONCodec{},
		}, append([]codectest.Case{
			{
				Name:  "empty struct",
				Value: struct{}{},
			},
			{
				Name:  "proto int64value",
				Value: &wrapperspb.Int64Value{Value: 4321},
			},
		}, genericTests...))
	})
}
