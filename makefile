build:
	GOOS=linux go build -o goproxy ./cmd/goproxy.go

app:
	./goproxy