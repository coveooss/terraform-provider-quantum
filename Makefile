SOURCES = $(wildcard **/*.go)

.PHONY: test
test:
	go test ./...

terraform-provider-quantum: $(SOURCES)
	go build ./...

.PHONY: build
build: terraform-provider-quantum

.PHONY: install
install: terraform-provider-quantum
	mv terraform-provider-quantum $(shell dirname $(shell which terraform))