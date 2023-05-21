# Golden tests in Go

[![Latest release](https://img.shields.io/github/v/release/hansmi/aurum)][releases]
[![CI workflow](https://github.com/hansmi/aurum/actions/workflows/ci.yaml/badge.svg)](https://github.com/hansmi/aurum/actions/workflows/ci.yaml)
[![Go reference](https://pkg.go.dev/badge/github.com/hansmi/aurum.svg)](https://pkg.go.dev/github.com/hansmi/aurum)

The aurum[^name-explanation] package implements golden tests for use in [Go
unit tests](https://pkg.go.dev/testing). Values expected from a computation are
stored in a file termed "golden file" and differences are reported as test
errors. Golden files are only written when requested via a command line flag
and if a logical change is detected. Version control is used to track and
review changes at the file level.

By default a generic JSON codec is used for storing expected values. The
implementation includes special handling of [Protocol
Buffers](https://protobuf.dev/) which need to be serialized through the
[`protojson`
package](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson).

String-like data may be stored in plain-text files (`TextCodec`). When only
protocol buffers are compared the textproto codec improves readability over
JSON (`TextProtoCodec`).

[^name-explanation]: _Aurum_ is Latin for _gold_.


## Example usage

```go
func init() {
  aurum.Init()
}

func Test(t *testing.T) {
  g := aurum.Golden{
    Dir: "./testdata",
  }
  g.Assert(t, "example", []string{"expected", "value"})
}
```

A more complete code example can be found in the [`example`
directory](./example/). To update the golden files:

```shell
go test -update_golden_files
```


## Alternatives

* [github.com/xorcare/golden](https://pkg.go.dev/github.com/xorcare/golden)
* [gotest.tools/v3/golden](https://pkg.go.dev/gotest.tools/v3/golden)

[releases]: https://github.com/hansmi/aurum/releases/latest

<!-- vim: set sw=2 sts=2 et : -->
