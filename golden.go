package aurum

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/hansmi/aurum/internal/codecutil"
	"go.uber.org/multierr"
)

type logFunc func(string, ...any)

// TB is the subset of [testing.TB] used for golden tests.
type TB interface {
	Helper()
	Errorf(format string, args ...any)
	Logf(format string, args ...any)
}

// Codec is the interface implemented by types used for marshalling and
// unmarshalling values.
type Codec = codecutil.Codec

var errGoldenMissing = errors.New("golden file is missing")
var errGoldenUnmarshalFailed = errors.New("unmarshalling golden value failed")
var errUpdateNotSupported = errors.New("updating files is not supported")

type Golden struct {
	// Directory for storing golden files. Only used if [FS] is not set.
	Dir string

	// Filesystem for accessing golden files. Updates are only possible if the
	// file system implements [WriteFileFS].
	//
	// Defaults to [os.DirFS] for [Dir].
	FS fs.FS

	// Codec for marshalling and unmarshalling values.
	//
	// Defaults to [JSONCodec].
	Codec Codec

	// Defaults to [Cmp].
	Comparer Comparer

	// Options for the default [Cmp] comparer.
	CmpOptions cmp.Options

	g *globalOptions
}

func (o *Golden) applyDefaults() {
	if o.g == nil {
		o.g = global
	}

	if o.FS == nil {
		o.FS = newWritableDirFS(o.Dir)
	}

	if o.Codec == nil {
		o.Codec = &JSONCodec{}
	}

	if o.Comparer == nil {
		o.Comparer = &Cmp{
			Options: o.CmpOptions,
		}
	}
}

func (o *Golden) unmarshal(data []byte, valueType reflect.Type) (any, error) {
	return codecutil.Unmarshal(o.Codec, data, valueType)
}

func (o *Golden) verifiedMarshal(value any, valueType reflect.Type) ([]byte, error) {
	gotBytes, err := o.Codec.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshalling value: %w", err)
	}

	if restored, err := o.unmarshal(gotBytes, valueType); err != nil {
		return nil, fmt.Errorf("unmarshalling previously marshalled value: %w", err)
	} else if err := o.Comparer.Equal(value, restored); err != nil {
		return nil, fmt.Errorf("value differs after marshalling and unmarshalling: %w", err)
	}

	return gotBytes, nil
}

func (o *Golden) readGolden(path string, t reflect.Type) (any, error) {
	wantBytes, err := fs.ReadFile(o.FS, path)
	if err != nil {
		if os.IsNotExist(err) {
			err = multierr.Combine(errGoldenMissing, err)
		} else {
			err = fmt.Errorf("reading golden file: %w", err)
		}

		return nil, err
	}

	value, err := o.unmarshal(wantBytes, t)
	if err != nil {
		err = multierr.Append(errGoldenUnmarshalFailed, err)
	}

	return value, err
}

func (o Golden) assert(name string, value any, logf logFunc) error {
	o.applyDefaults()

	if err := codecutil.CheckValueType(value); err != nil {
		return err
	}

	value, valueType := codecutil.NormalizeValue(value)

	valueBytes, err := o.verifiedMarshal(value, valueType)
	if err != nil {
		return err
	}

	filename := url.PathEscape(name)
	updatesEnabled := o.g.checkUpdatesEnabled()

	var considerWrite bool
	var diffErr error

	if want, err := o.readGolden(filename, valueType); err == nil {
		diffErr = o.Comparer.Equal(want, value)
		considerWrite = diffErr != nil
	} else if updatesEnabled && (errors.Is(err, errGoldenMissing) || errors.Is(err, errGoldenUnmarshalFailed)) {
		considerWrite = true
	} else {
		return err
	}

	if updatesEnabled && considerWrite {
		if diffErr != nil {
			logf("%v", diffErr)
		}

		if wffs, ok := o.FS.(WriteFileFS); !ok || wffs == nil {
			return fmt.Errorf("%w: %#v", errUpdateNotSupported, o.FS)
		} else if err := wffs.WriteFile(filename, valueBytes, 0o644); err != nil {
			return fmt.Errorf("writing golden file: %w", err)
		}

		logf("Wrote %d bytes to golden file %q.", len(valueBytes), filename)
	} else if diffErr != nil {
		return diffErr
	}

	return nil
}

// Assert checks whether the value matches the stored golden value read from
// a file.
//
// If enabled via a flag (see [Init]) golden files are updated if they're
// missing or differences in values are detected. The name is URL-escaped
// before being used as a filename and should be of a reasonable length (the
// exact limits depend on the underlying filesystem).
func (o *Golden) Assert(tb TB, name string, value any) {
	tb.Helper()

	if err := o.assert(name, value, tb.Logf); err != nil {
		tb.Errorf("%s", err.Error())
	}
}
