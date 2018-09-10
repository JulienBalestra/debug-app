TARGET=debug-app
CC=go
GOOS=linux
CGO_ENABLED=0

$(TARGET):
	GOOS=$(GOOS) CGO_ENABLED=$(CGO_ENABLED) $(CC) build -i -o probe-failures/probe-failures probe-failures/main.go
	GOOS=$(GOOS) CGO_ENABLED=$(CGO_ENABLED) $(CC) build -i -o write-tail/write-tail write-tail/main.go
	GOOS=$(GOOS) CGO_ENABLED=$(CGO_ENABLED) $(CC) build -i -o fork/fork fork/main.go
	GOOS=$(GOOS) CGO_ENABLED=$(CGO_ENABLED) $(CC) build -i -o probe-sleep/probe-sleep probe-sleep/main.go
	GOOS=$(GOOS) CGO_ENABLED=$(CGO_ENABLED) $(CC) build -i -o probe-sleep/probe-write probe-write/main.go

$(TARGET)-docker:
	docker build -t quay.io/julienbalestra/$(TARGET):master .

clean:
	$(RM) probe-failures/probe-failures write-tail/write-tail

re: clean $(TARGET)