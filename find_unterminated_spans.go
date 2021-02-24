package main

import (
	"fmt"
	"strings"
)

// TODO: unit test these methods

func findUnterminatedSpans(traces map[string]trace) unterminatedSpans {
	unclosedSpans := unterminatedSpans{}
	for _, trace := range traces {
		for spanID := range trace.Openings {
			if _, ok := trace.Closings[spanID]; ok {
				continue
			}

			var parentSpanChain = computeParentChain(trace, spanID)
			var operationChain = computeOperationChainString(trace, parentSpanChain)

			if _, ok := unclosedSpans[operationChain]; !ok {
				unclosedSpans[operationChain] = unterminatedSpan{}
			}
			unclosedSpans[operationChain] = unterminatedSpan{
				unclosedSpans[operationChain].Occurences + 1,
				append(unclosedSpans[operationChain].Traces, buildTraceTree(trace, parentSpanChain[0])),
			}

		}
	}
	return unclosedSpans
}

func buildTraceTree(trace trace, spanID string) *spanNode {
	opening, ok := trace.Openings[spanID]
	if !ok {
		return nil
	}

	children := []*spanNode{}
	for childSpanID, opening := range trace.Openings {
		if opening.ParentID == spanID {
			children = append(children, buildTraceTree(trace, childSpanID))
		}
	}

	closing, ok := trace.Closings[spanID]
	if !ok {
		return &spanNode{
			spanID,
			opening.ParentID,
			"unclosed",
			"unclosed",
			opening.OriginalRow,
			"unclosed",
			children,
		}
	}

	return &spanNode{
		spanID,
		opening.ParentID,
		fmt.Sprintf("%s:%s", closing.ServiceName, closing.OperationName),
		closing.Tags,
		opening.OriginalRow,
		closing.OriginalRow,
		children,
	}
}

func computeParentChain(trace trace, parentID string) []string {
	if parentID == "null" {
		return nil
	}

	var parentSpan, ok = trace.Openings[parentID]
	if !ok {
		return append([]string{"NO PARENT FOUND"}, parentID)
	}
	return append(computeParentChain(trace, parentSpan.ParentID), parentID)
}

func computeOperationChainString(trace trace, spanChain []string) string {
	result := []string{}
	for _, spanID := range spanChain {
		closingSpan, ok := trace.Closings[spanID]
		if !ok {
			result = append(result, "UNCLOSED SPAN")
		} else {
			result = append(result, fmt.Sprintf("%s:%s", closingSpan.ServiceName, closingSpan.OperationName))
		}
	}
	return strings.Join(result, " > ")
}
