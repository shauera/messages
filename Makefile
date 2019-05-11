CC = $(shell which go 2>/dev/null)

ifeq ($(CC),)
$(error "go is not in your system PATH")
else
$(info "go found")
endif

.PHONY: clean generate test all

all: clean generate test build

clean:
	$(RM) messages
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
publish: build
	tar cvf ds-appliance-controller.tar ds-appliance-controller ds-appliance-controller.version config.yaml
