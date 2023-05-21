package codecutil

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hansmi/aurum/internal/ref"
)

func TestNormalizeValue(t *testing.T) {
	for _, tc := range []struct {
		name     string
		value    any
		want     any
		wantType reflect.Type
	}{
		{name: "nil"},
		{
			name:     "string",
			value:    "test",
			want:     ref.Ref("test"),
			wantType: reflect.TypeOf((*string)(nil)),
		},
		{
			name:     "string pointer",
			value:    ref.Ref("another"),
			want:     ref.Ref("another"),
			wantType: reflect.TypeOf((*string)(nil)),
		},
		{
			name:     "pointer to pointer to pointer",
			value:    ref.Ref(ref.Ref(ref.Ref(123))),
			want:     ref.Ref(123),
			wantType: reflect.TypeOf((*int)(nil)),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, gotType := NormalizeValue(tc.value)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Value diff (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.wantType, gotType, cmp.Comparer(func(a, b reflect.Type) bool {
				// reflect.Type values are comparable. See also
				// <https://github.com/google/go-cmp/issues/80>.
				return a == b
			})); diff != "" {
				t.Errorf("Type diff (-want +got):\n%s", diff)
			}
		})
	}
}
