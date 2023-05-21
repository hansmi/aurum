package codectest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/aurum/internal/codecutil"
	"google.golang.org/protobuf/testing/protocmp"
)

type Codec = codecutil.Codec

type Case struct {
	Name             string
	Value            any
	WantMarshalErr   error
	WantUnmarshalErr error
}

func AssertAll(t *testing.T, c Codec, tests []Case) {
	t.Helper()

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			Assert(t, c, tc)
		})
	}
}

func Assert(t testing.TB, c Codec, tc Case) {
	t.Helper()

	value, valueType := codecutil.NormalizeValue(tc.Value)

	valueBytes, marshalErr := c.Marshal(value)

	if diff := cmp.Diff(tc.WantMarshalErr, marshalErr, cmpopts.EquateErrors()); diff != "" {
		t.Errorf("Marshal error diff (-want +got):\n%s", diff)
	}

	if marshalErr == nil {
		restored, unmarshalErr := codecutil.Unmarshal(c, valueBytes, valueType)

		if diff := cmp.Diff(tc.WantUnmarshalErr, unmarshalErr, cmpopts.EquateErrors()); diff != "" {
			t.Errorf("Unmarshal error diff (-want +got):\n%s", diff)
		}

		if unmarshalErr == nil {
			if diff := cmp.Diff(value, restored, protocmp.Transform(), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Restored value differs (-want +got):\n%s", diff)
			}
		}
	}
}
