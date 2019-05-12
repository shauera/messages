CC = $(shell which go 2>/dev/null)

ifeq ($(CC),)
$(error "go is not in your system PATH")
else
$(info "go found")
endif

.PHONY: clean generate test build swagger all

all: clean generate test build swagger

clean:
	$(RM) messages ./dist/swagger.json
generate:
	$(CC) generate
test: build
	$(CC) test ./... -v --cover
build: generate
	$(CC) get github.com/golang/dep/cmd/dep
	dep ensure
	$(CC) build
docker: build
	docker build -t shauer/messages .
docker-push: docker
	docker push shauer/messages
publish: build swagger
	tar cvf messages.tar messages ./dist
swagger:
	$(CC) get github.com/go-swagger/go-swagger/cmd/swagger
	swagger generate spec -o ./dist/swagger.json
