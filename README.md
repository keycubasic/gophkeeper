
GophKeeper Password Manager
===============
### GophKeeper is a client-server system that allows the user to reliably and securely store logins, passwords, binary data and other private information.
---

![Gopher, the Go mascot](/assets/images/gophkeeper.png)

Service use [Gateway](https://github.com/grpc-ecosystem/grpc-gateway).

Run go mod tidy to resolve the versions. Install by running

```shell
$ go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
$ go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
$ go install google.golang.org/protobuf/cmd/protoc-gen-go
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
$ go install github.com/googleapis/googleapis@latest
$ go install github.com/weblfe/go-swagger-merger
```

Generate stubs
```shell
$ scripts/proto.sh
```

Middleware problem. You may have to force the resource to be registered in go.mod.
```
require github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.1.0
```

Swagger
```
http://localhost:8080/swagger
```

Client build
```shell
$ go build -o cmd/gkclient/gkclient cmd/gkclient/main.go
```
