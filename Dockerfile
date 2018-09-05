FROM golang:1.10 as builder

COPY . /go/src/github.com/JulienBalestra/debug-app

RUN make -C /go/src/github.com/JulienBalestra/debug-app re

FROM busybox:latest

COPY --from=builder /go/src/github.com/JulienBalestra/debug-app/probe-failures/probe-failures /usr/local/bin/probe-failures
COPY --from=builder /go/src/github.com/JulienBalestra/debug-app/write-tail/write-tail /usr/local/bin/write-tail
COPY --from=builder /go/src/github.com/JulienBalestra/debug-app/fork/fork /usr/local/bin/fork
