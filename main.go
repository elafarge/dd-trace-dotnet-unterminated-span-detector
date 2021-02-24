package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// read from stdin
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: cat dd-trace-dotnet-log-file | dd-trace-dotnet-unterminated-span-count")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	traces := extractTracesFromLogStream(reader)

	unclosedSpans := findUnterminatedSpans(traces)

	result, err2 := json.MarshalIndent(map[string]interface{}{
		"analyzedTraceCount":    len(traces),
		"unterminatedSpanCount": len(unclosedSpans),
		"unterminatedSpans":     unclosedSpans,
	}, "", "  ")
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(string(result))
}