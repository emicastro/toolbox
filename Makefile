BINARY := tb
DIST := dist

.PHONY: all linux darwin clean

all: linux darwin

linux:
	cd cmd/tb && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -trimpath -o ../../$(DIST)/$(BINARY)-linux-amd64 .

darwin:
	cd cmd/tb && CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
		go build -trimpath -o ../../$(DIST)/$(BINARY)-darwin-arm64 .

clean:
	rm -rf $(DIST)
