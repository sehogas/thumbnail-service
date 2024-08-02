########################################################
override TARGET=thumbnail_service
VERSION=1.0.0
OS=linux
ARCH=amd64
FLAGS_VERSION="-X 'main.Version=$(VERSION)'"
FLAGS="-s -w -X 'main.Version=$(VERSION)'"
CGO=0
########################################################

run:
	go run -ldflags=$(FLAGS_VERSION) ./cmd/api/main.go 

bin:
	CGO_ENABLED=$(CGO) GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags=$(FLAGS) -o ./bin/$(TARGET) ./cmd/api/main.go 

install: 
	CGO_ENABLED=$(CGO) GOOS=$(OS) GOARCH=$(ARCH) go install -ldflags=$(FLAGS) 

clean:
	rm -rf ./bin/$(TARGET)

build:
	docker build --build-arg VERSION=$(VERSION) -t $(TARGET):$(VERSION) .
	docker tag $(TARGET):$(VERSION) $(TARGET):latest

start:
	docker run --rm -d --name $(TARGET) -p 3010:3010 $(TARGET):latest

start_with_network:
	docker run --rm -d --name $(TARGET) --network $(NETWORK) -p 3000:3000 $(TARGET):latest

stop:
	docker stop $(TARGET)

.PHONY: bin clean run install 



