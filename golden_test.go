package aurum

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/aurum/internal/ref"
	"github.com/hansmi/aurum/internal/testutil"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestGoldenAssert(t *testing.T) {
	type person struct {
		Name string
	}
	type addressBook []person

	type test struct {
		name           string
		initialContent *string
		value          any
		wantErr        error
		wantUpdate     bool
	}

	for _, updatesEnabled := range []bool{false, true} {
		tests := []test{
			{
				name:  "missing initial value",
				value: 0,
				wantErr: map[bool]error{
					false: os.ErrNotExist,
				}[updatesEnabled],
				wantUpdate: updatesEnabled,
			},
			{
				name:           "empty",
				initialContent: ref.Ref("[]"),
				value:          addressBook{},
			},
			{
				name: "existing file matches",
				initialContent: ref.Ref(`[
					{ "name": "John" },
					{ "name": "Jane" }
				]`),
				value: &addressBook{
					{"John"},
					{"Jane"},
				},
			},
			{
				name:           "value mismatch",
				initialContent: ref.Ref(`[]`),
				value: addressBook{
					{"A"},
					{"B"},
				},
				wantErr: map[bool]error{
					false: ErrValueDifference,
				}[updatesEnabled],
				wantUpdate: updatesEnabled,
			},
			{
				name:           "bad json in golden file",
				initialContent: ref.Ref("> bad json"),
				value:          "hello world",
				wantErr: map[bool]error{
					false: errGoldenUnmarshalFailed,
				}[updatesEnabled],
				wantUpdate: updatesEnabled,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name+map[bool]string{true: " with update"}[updatesEnabled], func(t *testing.T) {
				o := &Golden{
					g: &globalOptions{
						updatesEnabled: updatesEnabled,
					},
					Dir: t.TempDir(),
					CmpOptions: cmp.Options{
						cmpopts.EquateEmpty(),
					},
				}

				path := filepath.Join(o.Dir, "file")

				var fiBefore os.FileInfo

				if tc.initialContent == nil {
					testutil.MustNotExist(t, path)
				} else {
					testutil.MustWriteFile(t, path, *tc.initialContent)

					fiBefore = testutil.MustLstat(t, path)
				}

				err := o.assert(filepath.Base(path), tc.value, t.Logf)

				if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
					t.Errorf("Error diff (-want +got):\n%s", diff)
				}

				if err == nil {
					if tc.initialContent == nil && !updatesEnabled {
						testutil.MustNotExist(t, path)
					} else {
						fiAfter := testutil.MustLstat(t, path)

						changed := !(os.SameFile(fiBefore, fiAfter) &&
							fiBefore.Size() == fiAfter.Size() &&
							fiBefore.ModTime().Equal(fiAfter.ModTime()))

						if diff := cmp.Diff(tc.wantUpdate && updatesEnabled, changed); diff != "" {
							t.Errorf("File modification (-want, +got):\n%s", diff)
						}

						o.Assert(t, filepath.Base(path), tc.value)
					}
				} else if fiBefore == nil {
					testutil.MustNotExist(t, path)
				} else {
					testutil.MustLstat(t, path)
				}
			})
		}
	}
}

func TestGoldenAssertSimple(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   any
		wantErr error
	}{
		{
			name:    "nil",
			wantErr: os.ErrInvalid,
		},
		{
			name:  "string",
			value: "test\nline",
		},
		{
			name:  "int pointer",
			value: ref.Ref(1234),
		},
		{
			name:  "proto",
			value: &wrapperspb.StringValue{Value: "test"},
		},
		{
			name:    "proto non-pointer",
			value:   wrapperspb.StringValue{},
			wantErr: os.ErrInvalid,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			o := &Golden{
				g: &globalOptions{
					updatesEnabled: true,
				},
				Dir: t.TempDir(),
			}

			err := o.assert(tc.name, tc.value, t.Logf)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				o.Assert(t, tc.name, tc.value)
			}
		})
	}
}

func TestGoldenAssertNameExists(t *testing.T) {
	for _, updatesEnabled := range []bool{false, true} {
		o := &Golden{
			g: &globalOptions{
				updatesEnabled: updatesEnabled,
			},
			Dir: t.TempDir(),
		}

		testutil.MustMkdir(t, filepath.Join(o.Dir, "directory"))

		err := o.assert("directory", 0, t.Logf)

		if diff := cmp.Diff(syscall.EISDIR, err, cmpopts.EquateErrors()); diff != "" {
			t.Errorf("Error diff (-want +got):\n%s", diff)
		}
	}
}

func TestGoldenAssertFS(t *testing.T) {
	o := &Golden{
		g:     &globalOptions{},
		Codec: &TextCodec{},
		FS: fstest.MapFS{
			"foobar": {
				Data: []byte("content"),
			},
		},
	}

	if err := o.assert("foobar", "content", t.Logf); err != nil {
		t.Errorf("assert() failed: %v", err)
	}

	if err := o.assert("foobar", "changed", t.Logf); !errors.Is(err, ErrValueDifference) {
		t.Errorf("assert() returned %v, want %v", err, ErrValueDifference)
	}

	o.g.updatesEnabled = true

	if err := o.assert("foobar", "changed", t.Logf); !errors.Is(err, errUpdateNotSupported) {
		t.Errorf("assert() returned %v, want %v", err, errUpdateNotSupported)
	}
}
