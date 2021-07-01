# egosystem
PoC for learning how work proxy and load balancer like nginx server


## Load Balancer command

```bash
go build
```

> Start load balancer

```bash
./micros --backends http://localhost:3000,http://localhost:3001,http://localhost:3002 --port <port>
```

> Start API mock

```bash
go run  services/micro-a/api.go --port <port>
```
