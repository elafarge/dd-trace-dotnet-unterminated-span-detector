package main

import (
	"bufio"
	"io"
	"regexp"
)

// TODO: "Unit" test these regexps and methods
var dotnetOpeningRegexp = regexp.MustCompile(`Span started: \[s_id: (.+?), p_id: (.+?), t_id: (.+?)\]`)

var dotnetClosingRegexp = regexp.MustCompile(`Span closed: \[s_id: (.+?), p_id: (.+?), t_id: (.+?)\]`)
var dotnetClosingOperationNameRegexp = regexp.MustCompile(`Operation: (.+?),`)
var dotnetClosingTagsRegexp = regexp.MustCompile(`Tags: \[(.+)\]`)
var dotnetClosingServiceNameRegexp = regexp.MustCompile(`Service: (.+?),`)

type dotnetParser struct{}

func (p *dotnetParser) extractTracesFromLogStream(reader io.Reader) map[string]trace {
	traces := map[string]trace{}
	// read it line by line
	lineReader := bufio.NewReader(reader)
	for {
		lineBytes, err := lineReader.ReadString('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		line := string(lineBytes)

		if dotnetOpeningRegexp.MatchString(line) {
			var matches = dotnetOpeningRegexp.FindStringSubmatch(line)
			traceID := matches[3]
			spanID := matches[1]
			parentID := matches[2]

			if _, ok := traces[traceID]; !ok {
				traces[traceID] = trace{
					Openings: map[string]spanOpening{},
					Closings: map[string]spanClosing{},
				}
			}

			traces[traceID].Openings[spanID] = spanOpening{
				ParentID:    parentID,
				OriginalRow: line,
			}
		}

		if dotnetClosingRegexp.MatchString(line) {
			var matches = dotnetClosingRegexp.FindStringSubmatch(line)
			traceID := matches[3]
			spanID := matches[1]
			parentID := matches[2]

			var operationName string
			var match = dotnetClosingOperationNameRegexp.FindStringSubmatch(line)
			if match == nil {
				for {
					lineBytes, err := lineReader.ReadString('\n')

					if err == io.EOF {
						break
					} else if err != nil {
						panic(err)
					}

					line := string(lineBytes)

					match = dotnetClosingOperationNameRegexp.FindStringSubmatch(line)
					if match != nil {
						operationName = match[1]
						break
					}
				}
			} else {
				operationName = match[1]
			}

			var tags string
			match = dotnetClosingTagsRegexp.FindStringSubmatch(line)
			if match != nil {
				tags = match[1]
			}

			var serviceName string
			match = dotnetClosingServiceNameRegexp.FindStringSubmatch(line)
			if match != nil {
				serviceName = match[1]
			}

			if _, ok := traces[traceID]; !ok {
				traces[traceID] = trace{
					Openings: map[string]spanOpening{},
					Closings: map[string]spanClosing{},
				}
			}

			traces[traceID].Closings[spanID] = spanClosing{
				ParentID:      parentID,
				OriginalRow:   line,
				ServiceName:   serviceName,
				OperationName: operationName,
				Tags:          tags,
			}
		}
	}

	return traces
}
