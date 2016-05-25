package builtin

import (
	"fmt"
	"github.com/zalando/skipper/filters"
	"net/http"
)

type static struct {
	handler http.Handler
}

// Returns a filter Spec to serve static content from a file system
// location. Behaves similarly to net/http.FileServer. It shunts the route.
//
// Filter instances of this specification expect two parameters: a
// request path prefix and a local directory path. When processing a
// request, it clips the prefix from the request path, and appends the
// rest of the path to the directory path. Then, it uses the resulting
// path to serve static content from the file system.
//
// Name: "static".
func NewStatic() filters.Spec { return &static{} }

// "static"
func (spec *static) Name() string { return StaticName }

// Creates instances of the static filter. Expects two parameters: request path
// prefix and file system root.
func (spec *static) CreateFilter(config []interface{}) (filters.Filter, error) {
	if len(config) != 2 {
		return nil, fmt.Errorf("invalid number of args: %d, expected 1", len(config))
	}

	webRoot, ok := config[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid parameter type, expected string for web root prefix")
	}

	root, ok := config[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid parameter type, expected string for path to root dir")
	}

	return &static{http.StripPrefix(webRoot, http.FileServer(http.Dir(root)))}, nil
}

// Serves content from the file system and marks the request served.
func (f *static) Request(ctx filters.FilterContext) {
	ctx.Serve(ServeResponse(ctx.Request(), f.handler))
}

// Noop.
func (f *static) Response(filters.FilterContext) {}
