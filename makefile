build:
	GOOS=linux go build -o goproxy ./cmd/

app:
	./goproxy