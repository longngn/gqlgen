package errors

import (
	"fmt"

	"github.com/vektah/gqlgen/graphql"
)

type QueryError struct {
	Message       string        `json:"message"`
	Locations     []Location    `json:"locations,omitempty"`
	Path          []interface{} `json:"path,omitempty"`
	Rule          string        `json:"-"`
	ResolverError error         `json:"-"`
}

type Location struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func (a Location) Before(b Location) bool {
	return a.Line < b.Line || (a.Line == b.Line && a.Column < b.Column)
}

func Errorf(format string, a ...interface{}) *QueryError {
	return &QueryError{
		Message: fmt.Sprintf(format, a...),
	}
}

func (err *QueryError) Error() string {
	if err == nil {
		return "<nil>"
	}
	str := fmt.Sprintf("graphql: %s", err.Message)
	for _, loc := range err.Locations {
		str += fmt.Sprintf(" (line %d, column %d)", loc.Line, loc.Column)
	}
	return str
}

var _ error = &QueryError{}

type Builder struct {
	Errors []*QueryError
}

func (c *Builder) Errorf(format string, args ...interface{}) {
	c.Errors = append(c.Errors, Errorf(format, args...))
}

func (c *Builder) Error(err error) {
	c.Errors = append(c.Errors, Errorf("%s", err.Error()))
}

func ErrorWriter(errs []*QueryError) graphql.Marshaler {
	res := graphql.Array{}

	for _, err := range errs {
		if err == nil {
			res = append(res, graphql.Null)
			continue
		}

		errObj := &graphql.OrderedMap{}

		errObj.Add("message", graphql.String(err.Message))

		if len(err.Locations) > 0 {
			locations := graphql.Array{}
			for _, location := range err.Locations {
				locationObj := &graphql.OrderedMap{}
				locationObj.Add("line", graphql.Int(location.Line))
				locationObj.Add("column", graphql.Int(location.Column))

				locations = append(locations, locationObj)
			}

			errObj.Add("locations", locations)
		}
		res = append(res, errObj)
	}

	return res
}
