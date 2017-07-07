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
	cp terraform-provider-quantum $(shell dirname $(shell which terraform))

.PHONY: deploy
deploy:
	GOARCH=amd64 GOOS=linux go build -o terraform-provider-quantum_linux_x64 .
	GOARCH=amd64 GOOS=darwin go build -o terraform-provider-quantum_darwin_x64 .
	GOARCH=amd64 GOOS=windows go build -o terraform-provider-quantum_x64.exe .
