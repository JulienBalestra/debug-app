FROM golang:1.10 as builder

COPY . /go/src/github.com/JulienBalestra/debug-app

RUN make -C /go/src/github.com/JulienBalestra/debug-app re

FROM busybox:latest

COPY --from=builder /go/src/github.com/JulienBalestra/debug-app/debug-app /usr/local/bin/debug-app
