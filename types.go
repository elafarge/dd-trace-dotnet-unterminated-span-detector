package main

type spanOpening struct {
	ParentID    string
	OriginalRow string
}

type spanClosing struct {
	ParentID    string
	OriginalRow string

	ServiceName   string
	OperationName string
	Tags          string
}

// map of unclosed span parent chain string ==> unclosed occurences
type unterminatedSpans map[string]unterminatedSpan

type unterminatedSpan struct {
	Occurences int
	Traces     []*spanNode
}

type trace struct {
	Openings map[string]spanOpening
	Closings map[string]spanClosing
}

type spanNode struct {
	SpanID     string
	ParentID   string
	Name       string
	Tags       string
	OpeningLog string
	ClosingLog string

	Children []*spanNode
}
