package aurum

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestCmp(t *testing.T) {
	errTest := errors.New("test error")

	for _, tc := range []struct {
		name      string
		compare   Cmp
		valueWant any
		valueGot  any
		wantErr   error
	}{
		{name: "nil"},
		{
			name:      "strings equal",
			valueWant: "foo",
			valueGot:  "foo",
		},
		{
			name:      "strings difference",
			valueWant: "foo",
			valueGot:  "bar",
			wantErr:   ErrValueDifference,
		},
		{
			name:      "proto equal",
			valueWant: &wrapperspb.Int64Value{Value: 4321},
			valueGot:  &wrapperspb.Int64Value{Value: 4321},
		},
		{
			name:      "proto difference",
			valueWant: &wrapperspb.Int64Value{Value: 1},
			valueGot:  &wrapperspb.StringValue{Value: "bar"},
			wantErr:   ErrValueDifference,
		},
		{
			name: "panic with error value",
			compare: Cmp{
				Options: cmp.Options{
					cmp.Comparer(func(a, b int) bool {
						if a != b {
							panic(errTest)
						}
						return a == b
					}),
				},
			},
			valueWant: 0,
			valueGot:  1,
			wantErr:   errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.compare.Equal(tc.valueWant, tc.valueGot)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}
		})
	}
}
