.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build ./...

.PHONY: install
install:
	go install ./...

.PHONY: deploy
deploy:
	GOARCH=amd64 GOOS=linux go build -o .pkg/terraform-provider-quantum_linux_x64
	GOARCH=amd64 GOOS=darwin go build -o .pkg/terraform-provider-quantum_darwin_x64
	GOARCH=amd64 GOOS=windows go build -o .pkg/terraform-provider-quantum_x64.exe
