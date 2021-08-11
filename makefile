build:
	GOOS=linux go build -o proxyctl ./cmd/
	# CGO_ENABLED=0 GOOS=windows go build -o proxyctl.exe ./cmd/


gencert:
	go run ./examples/tools/generate_cert.go