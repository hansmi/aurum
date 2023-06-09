package example_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/hansmi/aurum"
	"google.golang.org/protobuf/types/known/structpb"
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

func TestTextProto(t *testing.T) {
	g := aurum.Golden{
		Dir:   "./testdata",
		Codec: &aurum.TextProtoCodec{},
	}

	g.Assert(t, "struct.textproto", func() *structpb.Struct {
		s, err := structpb.NewStruct(map[string]any{
			"hello":  true,
			"world":  "text",
			"number": 123,
		})
		if err != nil {
			t.Fatal(err)
		}
		return s
	}())
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<html><body>Hello World!</body></html>\n")
}

func TestHTTPHandler(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/example", nil)
	if err != nil {
		log.Fatal(err)
	}

	httpHandler(recorder, req)

	g := aurum.Golden{
		Dir:   "./testdata",
		Codec: &aurum.TextCodec{},
	}
	g.Assert(t, "http_handler_body", recorder.Body.String())
}
