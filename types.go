package main

type spanOpening struct {
	ParentID    string
	OriginalRow string
}

type spanClosing struct {
	ParentID      string
	OperationName string
	Tags          string
	OriginalRow   string
}

type trace struct {
	Openings map[string]spanOpening
	Closings map[string]spanClosing
}
