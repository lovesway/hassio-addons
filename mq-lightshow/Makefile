VERSION_TAG = $(shell git describe --abbrev=0 --tags)
VERSION_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
TEST_DIRS = .
COVERAGE_DATA_FILE = coverage.dat

ifeq ($(VERSION_BRANCH), master)
        VERSION = $(VERSION_TAG)
else
        VERSION = $(VERSION_TAG)_$(VERSION_BRANCH)
endif

.PHONY: default
default: native;
all: clean lint test native docker
lint:
	@echo "-[ Linting with golangci-lint ]-"
	golangci-lint run --enable-all --exclude-use-default=false
test:
	@echo "-[ Running Tests with Race Condition Check ]-"
	go test -race -covermode=atomic $(TEST_DIRS)
test-cover:
	@echo "-[ Running Tests with Coverage in HTML Format ]-"
	go test -covermode=atomic -coverprofile=$(COVERAGE_DATA_FILE) $(TEST_DIRS)
	go tool cover -html=$(COVERAGE_DATA_FILE)
docker:
	@echo "-[ Building for Docker ]-"
	 docker build --build-arg ARCH=amd64 --build-arg GOOS=linux --build-arg GOARCH=amd64 \
	 	--build-arg BUILD_VERSION=$(VERSION) -t amd64-mq-lightshow:$(VERSION) .
docker-all-arch:
	@echo "-[ Building for Docker - All Architectures]-"
	 docker buildx build --platform linux/amd64 --build-arg ARCH=amd64 --build-arg GOOS=linux \
	 	--build-arg GOARCH=amd64 --build-arg BUILD_VERSION=$(VERSION) -t lovesway/amd64-mq-lightshow:$(VERSION) .
	 docker buildx build --platform linux/arm/v7 --build-arg ARCH=arm32v7 --build-arg GOOS=linux \
	 	--build-arg GOARCH=arm --build-arg BUILD_VERSION=$(VERSION) -t lovesway/armhf-mq-lightshow:$(VERSION) .
	 docker buildx build --platform linux/arm/v7 --build-arg ARCH=arm32v7 --build-arg GOOS=linux \
	 	--build-arg GOARCH=arm --build-arg BUILD_VERSION=$(VERSION) -t lovesway/armv7-mq-lightshow:$(VERSION) .
	 docker buildx build --platform linux/arm64 --build-arg ARCH=arm64v8 --build-arg GOOS=linux \
	 	--build-arg GOARCH=arm64 --build-arg BUILD_VERSION=$(VERSION) -t lovesway/aarch64-mq-lightshow:$(VERSION) .
	 docker buildx build --platform linux/386 --build-arg ARCH=i386 --build-arg GOOS=linux \
	 	--build-arg GOARCH=386 --build-arg BUILD_VERSION=$(VERSION) -t lovesway/i386-mq-lightshow:$(VERSION) .
native:
	@echo "-[ Building Natively ]-"
	go build -ldflags "-X main.version=$(VERSION)"
clean:
	@echo "-[ Cleaning ]-"
	if [ -f mq-lightshow ]; then rm mq-lightshow; fi;
	if [ -f $(COVERAGE_DATA_FILE) ]; then rm $(COVERAGE_DATA_FILE); fi;
help:
	@echo "-[ Showing Help ]-"
	@echo "usage: make [target]\n"
	@echo "  native		Build using native golang (default if no target supplied)."
	@echo "  docker		Build using Docker."
	@echo "  docker-all-arch	Build all architectures using Docker.\n"
	@echo "  lint			Lint with golangci-lint."
	@echo "  clean			Deletes unneeded files produced by testing and building.\n"
	@echo "  test			Run tests with race condition check."
	@echo "  test-cover		Run tests with code coverage report in HTML format.\n"
	@echo "  all			Runs: clean lint test native docker.\n"
	@echo "  help			Show this help."
