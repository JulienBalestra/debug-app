TARGET=debug-app
CC=go
GOOS=linux
CGO_ENABLED=0

$(TARGET):
	GOOS=$(GOOS) CGO_ENABLED=$(CGO_ENABLED) $(CC) build -i -o probe-failures/probe-failures probe-failures/main.go
	GOOS=$(GOOS) CGO_ENABLED=$(CGO_ENABLED) $(CC) build -i -o write-tail/write-tail write-tail/main.go

$(TARGET)-docker:
	docker build -t quay.io/julienbalestra/$(TARGET):master .

clean:
	$(RM) probe-failures/probe-failures write-tail/write-tail

re: clean $(TARGET)