# GOproxy
PoC for learning how work proxy and load balancer like nginx server


## Load Balancer command

```bash
GOOS=linux go build -o goproxy .
GOOS=windows go build -o goproxy.exe .
```

> Start Proxy server

`genkey` command in true generate random secretkey storaged on badgerdb on `badger.data`

```bash
./goproxy --genkey true --port <port>
```

> Start load balancer

```bash
./goproxy --backends http://goproxy.com:3000,http://goproxy.com:3001,http://goproxy.com:3002 --port <port>
```


> Start API PoC service

```bash
go run  services/micro-a/api.go --port <port>
```

Install badger db on window if you don`t use CGO
```bash
CGO_ENABLED=0 go get github.com/dgraph-io/badger/v3
```
