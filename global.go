package aurum

import (
	"errors"
	"flag"
	"sync"
)

const DefaultUpdateFlagName = "update_golden_files"

type globalOptions struct {
	mu             sync.Mutex
	initialized    bool
	flagSet        *flag.FlagSet
	flagName       string
	updatesEnabled bool
}

func (g *globalOptions) init(opts []InitOption) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.initialized {
		return errors.New("package can only be initialized once")
	}

	if g.flagSet == nil {
		g.flagSet = flag.CommandLine
	}

	for _, opt := range opts {
		opt.apply(g)
	}

	if g.flagName != "" {
		g.flagSet.BoolVar(&g.updatesEnabled, g.flagName, g.updatesEnabled,
			"Update golden test files in-place.")
	}

	g.initialized = true

	return nil
}

func (g *globalOptions) checkUpdatesEnabled() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.updatesEnabled
}

var global = &globalOptions{
	flagName: DefaultUpdateFlagName,
}

// Interface implemented by initialization options.
type InitOption interface {
	apply(*globalOptions)
}

type withFlagName string

func (n withFlagName) apply(opt *globalOptions) {
	opt.flagName = string(n)
}

// Override the flag name.
func WithFlagName(name string) InitOption {
	return withFlagName(name)
}

// Initialize the package and register a command line flag. Must be called
// before parsing flags. Example usage in a test file:
//
//	func init() {
//	  aurum.Init()
//	}
//
// If not called the default values are used.
func Init(opts ...InitOption) {
	if err := global.init(opts); err != nil {
		panic(err)
	}
}
