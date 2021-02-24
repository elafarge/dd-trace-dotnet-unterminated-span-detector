FROM golang:1-alpine AS builder

RUN apk update && apk add --no-cache git
RUN echo $GOPATH

WORKDIR $GOPATH/src/app
COPY . .
# disable CGO, we don't need it :)
RUN pwd
RUN go install && go get -v -d
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $GOPATH/bin/dd-trace-dotnet-unterminated-span-detector

FROM scratch
COPY --from=builder /go/bin/dd-trace-dotnet-unterminated-span-detector .
ENTRYPOINT ["dd-trace-dotnet-unterminated-span-detector"]
