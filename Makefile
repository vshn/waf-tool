.PHONY: all build clean dist fmt tests

all: clean lint tests build

fmt:
	@echo 'Reformat Go code ...'
	go fmt ./...

vet:
	@echo 'Examine Go code ...'
	go vet ./...

lint: fmt vet
	@echo 'Check for uncommitted changes ...'
	git diff --exit-code

test:
	@echo 'Run all tests ...'
	go test --cover ./...

build:
	@echo 'Build waf-tool binary ...'
	go build

dist:
	@echo 'Build all distributions ...'
	goreleaser release --snapshot --rm-dist --skip-sign

clean:
	@echo 'Clean up test cache and build artifacts ...'
	go clean -testcache
	rm -rf waf-tool dist/