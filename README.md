# nGOnx

It`s a simple proxy server written with [go](https://github.com/golang/go).
The core services are based on [nginx](https://github.com/nginx) server, and [traefik](https://github.com/traefik/traefik)
Our roadmap you can find here [project board](https://piqba.notion.site/f98799ba3d384526ac6591247b12481c?v=4ca0c832682749e99dcc72a66fdd71a1)

## Features

* Run ngonx as a reverse proxy
* Run ngonx as a grpc proxy
* Run ngonx as a load balancer (round robin)
* Run ngonx as a static web server
* Project collaborative and open source

## Download ngonxctl

Go to last release and download [ngonxctl](https://github.com/kenriortega/ngonx/releases) binary for you OS

## **nGOnx** commands


```bash
ngonxctl -h
```

```bash
This is Ngonx ctl a proxy reverse inspired on nginx & traefik

Usage:
  ngonxctl [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  grpc        Run ngonx as a grpc proxy
  help        Help about any command
  lb          Run ngonx as a load balancer (round robin)
  proxy       Run ngonx as a reverse proxy
  setup       Create configuration file it`s doesn`t exist
  static      Run ngonx as a static web server
  version     Print the version number of ngonxctl

Flags:
  -f, --cfgfile string   File setting.yml (default "ngonx.yaml")
  -p, --cfgpath string   Config path only  (default "path/binary")
  -h, --help             help for ngonxctl

Use "ngonxctl [command] --help" for more information about a command.

```

> Setup cmd

Create the yaml file that contain a generic data

```bash
Create configuration file it`s doesn`t exist

Usage:
  ngonxctl setup [flags]

Flags:
  -h, --help   help for setup

Global Flags:
  -f, --cfgfile string   File setting.yml (default "ngonx.yaml")
  -p, --cfgpath string   Config path only  (default "path/binary")

```

> Yaml file with generic data


```yaml
# Static web server like nginx
static_server:
  host_server: 0.0.0.0
  port_server: 8080
  static_files: ./examples/dist
  ssl_server:
    enable: true
    ssl_port: 8443
    crt_file: ./ssl/cert.pem
    key_file: ./ssl/key.pem
# Grpc reverse proxy transparent
grpc:
  listener_grpc: "0.0.0.0:50000"
  client_crt: ./scripts/ca.crt
  ssl_grpc:
    enable: false
    ssl_port: 50443
    crt_file: ./scripts/server.crt
    key_file: ./scripts/server.pem
  endpoints_grpc:
    - name: /calculator.CalculatorService
      host_uri: 0.0.0.0:50050
# Reverse Proxy
proxy:
  host_proxy: 0.0.0.0
  port_proxy: 30000
  port_exporter_proxy: 10000
  ssl_proxy:
    enable: false
    ssl_port: 443
    crt_file: ./ssl/cert.pem
    key_file: ./ssl/key.pem
  cache_proxy:
    engine: badger # badgerDB|redis
    key: secretKey
  security:
    type: apikey # apikey|jwt|none
  # maps of microservices with routes
  services_proxy:
      - name: microA
        host_uri: http://localhost:3000
        endpoints:
          - path_endpoints: /api/v1/health/
            path_proxy: /health/
            path_protected: false

          - path_endpoints: /api/v1/version/
            path_proxy: /version/
            path_protected: true
```


> Version cmd show Build Time, version hash and version for the current binary

```bash
Print the version number of ngonxctl

Usage:
  ngonxctl version [flags]

Flags:
  -h, --help   help for version

Global Flags:
  -f, --cfgfile string   File setting.yml (default "ngonx.yaml")
  -p, --cfgpath string   Config path only  (default "path/binary")

```


> Start Proxy server first time

```bash
Run ngonx as a reverse proxy

Usage:
  ngonxctl proxy [flags]

Flags:
      --genkey           Action for generate hash for protected routes
  -h, --help             help for proxy
      --port int         Port to serve to run proxy (default 5000)
      --prevkey string   Action for save a previous hash for protected routes to validate JWT

Global Flags:
  -f, --cfgfile string   File setting.yml (default "ngonx.yaml")
  -p, --cfgpath string   Config path only  (default "path/binary")

```

`genkey` command in true generate random secretkey and save on badgerdb on `badger.data`

```bash
./ngonxctl proxy -port 5000 -genkey true
```

`prevkey` command receive a custom secretkey and save this on badgerdb on `badger.data`

```bash
./ngonxctl proxy -port 5000 -prevkey <secretKey>
```

> Start Proxy server

```bash
./ngonxctl proxy -port 5000
```

> Start load balancer

```bash
Run ngonx as a load balancer (round robin)

Usage:
  ngonxctl lb [flags]

Flags:
      --backends string   Load balanced backends, use commas to separate (default "ngonx.yaml")
  -h, --help              help for lb
      --port int          Port to serve to run load balancing  (default 4000)

Global Flags:
  -f, --cfgfile string   File setting.yml (default "ngonx.yaml")
  -p, --cfgpath string   Config path only  (default "path/binary")

```


```bash
./ngonxctl lb --backends "http://localhost:5000,http://localhost:5001,http://localhost:5002"
```
> Start static files server

```bash
Run ngonx as a static web server

Usage:
  ngonxctl static [flags]

Flags:
  -h, --help   help for static

Global Flags:
  -f, --cfgfile string   File setting.yml (default "ngonx.yaml")
  -p, --cfgpath string   Config path only  (default "path/binary")

```

```bash
./ngonxctl static
```

> Start grpc proxy server

```bash
Run ngonx as a grpc proxy

Usage:
  ngonxctl grpc [flags]

Flags:
  -h, --help   help for grpc

Global Flags:
  -f, --cfgfile string   File setting.yml (default "ngonx.yaml")
  -p, --cfgpath string   Config path only  (default "path/binary")

```

```bash
./ngonxctl grpc
```

Metrics
-----------

Currently ngonx use prometheus as a metric collector. The main service `proxy` expose for port 10000 on route `/metrics`

```bash
curl http://localhost:10000/metrics
```


Management API & Web(coming...)
-----------

Currently ngonx use port 10001 for export a simple api to check all services 

```bash
curl http://localhost:10001/api/v1/mngt/
```

  Method   | Path | 
  ---------|----------------
  GET | /      
  GET | /health      
  GET | /readiness      
  GET | /wss      

And you can view on the ngonx-UI SPA

![Service Discovery](/docs/service1.jpeg)
![Service Discovery](/docs/service2.jpeg)


How to interact with app for testing reasson
-----

> Start API PoC service

```bash
go run  services/micro-a/api.go --port <port>
```

> Start Gprc services from example folder

Excute from  makefile the following cmds

```bash
make grpcsvr # for start server

make make grpccli #for start client
```

Install badger db on window if you don`t use CGO
```bash
CGO_ENABLED=0 go get github.com/dgraph-io/badger/v3
```

## Generate customs SSL (key,cert,pem,crs,cnf)

```bash

cd scripts

chmod +x ./generate.sh

./generate.sh
```

## Generate SSL (throug let`s encrypt)

Check this post ["Contratar un certificado SSL Gratis con Let's Encrypt y configurar Nginx"](https://www.albertcoronado.com/2020/05/05/contratar-un-certificado-ssl-gratis-con-lets-encrypt-y-configurar-nginx/ "Contratar un certificado SSL Gratis con Let's Encrypt y configurar Nginx") by [@acoronado](https://www.albertcoronado.com/) for more information.

```bash

apt-get install -y certbot

certbot certonly \
    -d midominio.com \
    --noninteractive \
    --standalone \
    --agree-tos \
    --register-unsafely-without-email
```

Then configure our `ngonx.yaml` and add the ssl files generated by `certbot`

```yaml
# For SSL static server
  ssl_server:
    enable: true
    ssl_port: 8443
    crt_file: /etc/letsencrypt/live/mydomain.com/fullchain.pem
    key_file: /etc/letsencrypt/live/mydomain.com/privkey.pem
# For SSL proxy server
  ssl_proxy:
    enable: true
    ssl_port: 443
    crt_file: /etc/letsencrypt/live/mydomain.com/fullchain.pem
    key_file: /etc/letsencrypt/live/mydomain.com/privkey.pem
```

When the ssl certification was expired you need to renew all certificates
with the following commnad

```bash
certbot renew
```

---


> Yaml file fields

BenchMarking
------------

```bash
ab -c 1000 -n 10000 http://localhost:<proxyPort>/health
```