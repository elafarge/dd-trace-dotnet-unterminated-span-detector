Datadog Tracing - Unterminated Span Detector
============================================

When manually instrumenting your code with datadog traces, it's possible to
unintentionnaly forget to close a span. Although the manual instrumentation
libraries provide interfaces that help prevent that (leveraging `defer` in Go,
or `using` in C# for instance).

However, when the execution workflow is complex, or involves stream
manipulation in .NET, forgetting to close a span is rather common.

If a trace contains an unclosed span, it seems (with .NET manual instrumentation
at least) that the entire trace isn't forwarded to datadog, leading to
unpleasant surprises (missing traces) at troubleshooting time.

This tiny tool analyses the manual instrumentation library's DEBUG logs to
detect unclosed spans and report them, along with all the information you need
to pinpoint the part of the code where closing the span has been forgotten.

TODO: screenshot of a report

Usage
-----

#### Where should I run this check ?

We strongly disrecommend colocating this program with your production services:
- enabling dd-trace DEBUG logging in production isn't recommended
- this program can be quite CPU heavy (it uses Go's standard `regexp` package,
  which isn't known for being CPU-savvy)

If you absolutely need to target your production environment, we recommend that
you first retrieve the dd-trace logs from your production instances, and run
this program on dedicated instances.

In our experience, running this check on our CI pipelines, right after
end-to-end tests worked well: the end-to-end tests should have invoked most of
our codepaths (and therefore, spans should have been emitted for these code
paths).
Also, it makes it possible to check for unterminated spans before even
releasing you code to a remote (staging/production/...) environment :)

#### Running the check, using our Docker image

We strongly recommend using our docker image (which weight only 2MiB). Feel free
to built it yourself (see the Contributing section below) if using public Docker
images is against your security policy.

1. Make sure you enable DEBUG logs by setting the `DD_TRACE_DEBUG=true` env. var
   before launching the traced process.

2. Fetch the tracer logs in `/var/log/datadog` (for .NET applications, these are
   the logs matching the following pattern
   `/var/log/datadog/dotnet/dotnet-tracer-managed-*.log`)

3. Run the check, it should output a JSON payload containing information about
   unclosed spans.

```shell
cat dd-log-file.log | docker run -it --rm elafarge/dd-trace-dotnet-unterminated-span-detector
```

Contributing
------------

Contributions are always welcome, in particular we warmly welcome "parser"
implementation for logs emitted by instrumentation for new language (currently,
only .NET is supported).

There's no particular exotic coding guideline, just make sure you:
* `gofmt` your code before opening PRs
* add unit tests :)
* write a comprehensive PR description to ease the reviewer's job

### Running the code locally

In this folder:
```shell
cat log-file.log | go run ./*.go
```

### Running tests
In this folder:
```shell
go test
```

### Building and pushing a custom docker image
In this folder:
```shell
docker build -t YOUR_REPO/dd-trace-dotnet-unterminated-span-detector .
# Then
docker push YOUR_REPO/dd-trace-dotnet-unterminated-span-detector
```

Authors
-------
* Étienne Lafarge <etienne.lafarge _at_ gmail.com>
