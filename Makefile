default: test-all
 
test-all: test test-coverage test-integration


.PHONY: test
test:
	go test -v  ./...

# runs coverage tests and generates the coverage report
.PHONY: test-coverage
test-coverage:
	go test ./... -v -coverpkg=./...

# runs integration tests
.PHONY: test-coverpkg
test-integration:
	go test ./... -tags=integration ./...