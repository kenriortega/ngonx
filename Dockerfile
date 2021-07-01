FROM golang:alpine as build
WORKDIR /go/src/github.com/egosystem.org/micros
COPY . .
RUN go build -o /proxy

FROM alpine:3.13.2
COPY --from=build /proxy /proxy
ENTRYPOINT ["/proxy"]