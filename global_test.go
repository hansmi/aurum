package aurum

import (
	"flag"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGlobalOptionsInit(t *testing.T) {
	for _, tc := range []struct {
		name     string
		opts     []InitOption
		wantFlag string
		wantErr  error
	}{
		{name: "empty"},
		{
			name: "named flag",
			opts: []InitOption{
				WithFlagName(DefaultUpdateFlagName),
			},
			wantFlag: "update_golden_files",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := globalOptions{
				flagSet: flag.NewFlagSet("", flag.PanicOnError),
			}

			err := g.init(tc.opts)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Error diff (-want +got):\n%s", diff)
			}

			var flagNames []string

			g.flagSet.VisitAll(func(f *flag.Flag) {
				flagNames = append(flagNames, f.Name)
			})

			if tc.wantFlag == "" {
				if len(flagNames) != 0 {
					t.Errorf("Flag set is not empty: %s", flagNames)
				}
			} else if got := g.flagSet.Lookup(tc.wantFlag); got == nil {
				t.Errorf("Flag %q not set", tc.wantFlag)
			}
		})
	}
}
