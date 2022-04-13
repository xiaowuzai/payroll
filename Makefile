GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GORUN=$(GOCMD) run
BINARY_NAME=payroll
DOCKER=docker
DOCKERBUILD=$(DOCKER) build

all: docker
build: clean
	cd ./$(BINARY_NAME) && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) github.com/xiaowuzai/$(BINARY_NAME)/cmd/server
docker: build
	tar -cvf - $(BINARY_NAME) $(BINARY_NAME)/configs | $(DOCKERBUILD) -t $(BINARY_NAME) -f $(BINARY_NAME)/Dockerfile -
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)/$(BINARY_NAME)
