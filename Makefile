CC = $(shell which go 2>/dev/null)

ifeq ($(CC),)
$(error "go is not in your system PATH")
else
$(info "go found")
endif

TARGET_FILE = messages
SWAGGER_FILE = ./dist/swagger.json

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: clean test build swagger docker docker-push publish

clean:
	$(RM) messages ./dist/swagger.json

test: build
	$(CC) test ./... -v --cover

$(TARGET_FILE): $(SRC)
	$(CC) get github.com/golang/dep/cmd/dep
	dep ensure
	$(CC) build
build: $(TARGET_FILE)

docker: build
	docker build -t shauer/messages .
docker-push: docker
	docker push shauer/messages

$(SWAGGER_FILE): build
	$(CC) get github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate spec -m -o $(SWAGGER_FILE)
swagger: $(SWAGGER_FILE)
	
publish: build swagger
	tar cvf $(TARGET_FILE).tar $(TARGET_FILE) config.yml ./dist
