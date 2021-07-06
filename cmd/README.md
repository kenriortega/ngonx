# GOproxy
PoC for learning how work proxy and load balancer like nginx server


## GOproxy command

> Build

```bash
make build
```

> Start Proxy server first time

`genkey` command in true generate random secretkey storaged on badgerdb on `badger.data`

```bash
./goproxy -portProxy 5000 -genkey true
```

> Start Proxy server

```bash
./goproxy -portProxy 5001
```

> Start load balancer

```bash
./goproxy -type lb --backends "http://localhost:5000,http://localhost:5001,http://localhost:5002"
```

> Start API PoC service

```bash
go run  services/micro-a/api.go --port <port>
```

Install badger db on window if you don`t use CGO
```bash
CGO_ENABLED=0 go get github.com/dgraph-io/badger/v3
```
