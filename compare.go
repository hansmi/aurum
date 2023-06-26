package aurum

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

var ErrValueDifference = errors.New("values are not equal")

// Comparer is an interface implemented by types used for checking whether two
// values are equal or equivalent (the distinction is up to the
// implementation).
type Comparer interface {
	Equal(want, got any) error
}

// Cmp compares values using [cmp.Diff].
type Cmp struct {
	// Options for comparing values, e.g. [cmpopts.EqualEmpty]. If one of the
	// compared values is a [proto.Message] then the [protocmp.Transform]
	// option is automatically added.
	Options cmp.Options
}

var _ Comparer = (*Cmp)(nil)

func (c Cmp) Equal(want, got any) (err error) {
	// cmp never returns errors and panics instead.
	defer func() {
		if r := recover(); r != nil {
			rerr, ok := r.(error)
			if !ok || rerr == nil {
				rerr = errors.New(fmt.Sprint(r))
			}
			multierr.AppendInto(&err, rerr)
		}
	}()

	opts := c.Options

	for _, v := range []any{want, got} {
		if _, ok := v.(proto.Message); ok {
			opts = append(cmp.Options{protocmp.Transform()}, opts...)
			break
		}
	}

	if diff := cmp.Diff(want, got, opts...); diff != "" {
		return fmt.Errorf("%w (-want +got):\n%s", ErrValueDifference, diff)
	}

	return nil
}
