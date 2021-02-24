package main

import (
	"io"
)

// describes a generic parser that - given a stream of logs - will extract the generated traces from
// these logs "Span Opened" and "Span closed" events
type parser interface {
	extractTracesFromLogStream(reader io.Reader) map[string]trace
}
