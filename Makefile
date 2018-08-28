TARGET=debug-app
CC=go
GOOS=linux
CGO_ENABLED=0

$(TARGET):
	GOOS=$(GOOS) CGO_ENABLED=$(CGO_ENABLED) $(CC) build -i .

$(TARGET)-docker:
	docker build -t quay.io/julienbalestra/$(TARGET):master .

clean:
	$(RM) $(TARGET)

re: clean $(TARGET)