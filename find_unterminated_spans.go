package main

import "fmt"

// map of unclosed span parent chain string ==> unclosed occurences
type unterminatedSpans map[string]int

func findUnterminatedSpans(traces map[string]trace) unterminatedSpans {
	unclosedSpans := unterminatedSpans{}
	for _, trace := range traces {
		for spanID, spanOpening := range trace.Openings {
			if _, ok := trace.Closings[spanID]; ok {
				continue
			}

			var parentChain = fmt.Sprintf("%s UNCLOSED SPAN", computeParentChain(trace, spanOpening.ParentID))
			if _, ok := unclosedSpans[parentChain]; !ok {
				unclosedSpans[parentChain] = 0
			}
			unclosedSpans[parentChain] = unclosedSpans[parentChain] + 1
		}
	}
	return unclosedSpans
}

func computeParentChain(trace trace, parentID string) string {
	if parentID == "null" {
		return ""
	}

	var parentSpan, ok = trace.Closings[parentID]
	var nextParentID = parentSpan.ParentID
	var operationName = parentSpan.OperationName
	var tags = parentSpan.Tags
	if !ok {
		// missing parent closing, let's try with the opening
		var parentSpanOpening, okdoki = trace.Openings[parentID]
		if !okdoki {
			return "MISSING PARENT SPAN - BROKEN CHAIN"
		}
		nextParentID = parentSpanOpening.ParentID
		operationName = "UNCLOSED SPAN"
	}

	return fmt.Sprintf("%s %s[%s] >", computeParentChain(trace, nextParentID), operationName, tags)
}
