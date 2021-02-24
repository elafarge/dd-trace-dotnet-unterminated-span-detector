package main

import (
	"bufio"
	"io"
	"regexp"
)

var openingRegexp = regexp.MustCompile(`Span started: \[s_id: (.+?), p_id: (.+?), t_id: (.+?)\]`)

var closingRegexp = regexp.MustCompile(`Span closed: \[s_id: (.+?), p_id: (.+?), t_id: (.+?)\]`)
var closingOperationNameRegexp = regexp.MustCompile(`Operation: (.+),`)
var closingTagsRegexp = regexp.MustCompile(`Tags: [(.+)]`)

func extractTracesFromLogStream(reader io.Reader) map[string]trace {
	traces := map[string]trace{}
	// read it line by line
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if openingRegexp.MatchString(line) {
			var matches = openingRegexp.FindStringSubmatch(line)
			originalRow := matches[0]
			traceID := matches[3]
			spanID := matches[1]
			parentID := matches[2]

			if _, ok := traces[traceID]; !ok {
				traces[traceID] = trace{
					Openings: map[string]spanOpening{},
					Closings: map[string]spanClosing{},
				}
			}

			traces[traceID].Openings[spanID] = spanOpening{parentID, originalRow}
		}

		if closingRegexp.MatchString(line) {
			var matches = closingRegexp.FindStringSubmatch(line)
			originalRow := matches[0]
			traceID := matches[3]
			spanID := matches[1]
			parentID := matches[2]

			var operationName string
			var match = closingOperationNameRegexp.FindStringSubmatch(line)
			if match == nil {
				for scanner.Scan() {
					line = scanner.Text()
					match = closingOperationNameRegexp.FindStringSubmatch(line)
					if match != nil {
						operationName = match[1]
						break
					}
				}
			} else {
				operationName = match[1]
			}

			var tags string
			match = closingTagsRegexp.FindStringSubmatch(line)
			if match != nil {
				tags = match[1]
			}

			traces[traceID].Closings[spanID] = spanClosing{parentID, operationName, tags, originalRow}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return traces
}
