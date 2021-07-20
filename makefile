build:
	GOOS=linux go build -o goproxy ./cmd/
	# CGO_ENABLED=0 GOOS=windows go build -o goproxy.exe ./cmd/

app:
	./goproxy

gencert:
	go run ./examples/tools/generate_cert.go