BINARY_NAME=RateQuerySvr
GOBIN=./bin/
build:
	go mod tidy
	go build -o ${GOBIN}${BINARY_NAME} main.go

clean:
	@echo "  >  Cleaning build cache"
	rm $(GOBIN)$(BINARY_NAME)