package example_test

import (
	"sort"
	"testing"
	"time"

	"github.com/hansmi/aurum"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func init() {
	aurum.Init()
}

func TestSortStrings(t *testing.T) {
	g := aurum.Golden{
		Dir: "./testdata",
	}

	for _, tc := range []struct {
		name   string
		values []string
	}{
		{name: "empty"},
		{
			name:   "names",
			values: []string{"liam", "noah", "oliver", "emma", "olivia", "amelia"},
		},
		{
			name:   "numbers",
			values: []string{"one", "two", "three", "four", "five", "six", "eight", "nine", "ten"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			values := append([]string{}, tc.values...)

			sort.Strings(values)

			g.Assert(t, tc.name, values)
		})
	}
}

func TestProtoMessage(t *testing.T) {
	g := aurum.Golden{
		Dir: "./testdata",
	}
	g.Assert(t, "proto", timestamppb.New(time.Date(2000, time.January, 1, 3, 2, 1, 0, time.UTC)))
}
