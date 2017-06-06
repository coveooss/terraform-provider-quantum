build:
	go build -o /usr/local/bin/terraform-provider-quantum

deploy:
	GOOS=linux go build -o .pkg/terraform-provider-quantum_linux
	GOOS=darwin go build -o .pkg/terraform-provider-quantum_darwin
	GOOS=windows go build -o .pkg/terraform-provider-quantum.exe
