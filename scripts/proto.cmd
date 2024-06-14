protoc -I . ^
    -I %GOPATH%/pkg/mod/github.com/googleapis/googleapis@v0.0.0-20240605231841-21c839628513 ^
    --go_out . ^
    --go_opt paths=import ^
    --go-grpc_out . ^
    --go-grpc_opt paths=import ^
    --grpc-gateway_out . ^
    --grpc-gateway_opt paths=import ^
    --grpc-gateway_opt logtostderr=true ^
    --grpc-gateway_opt generate_unbound_methods=true ^
    --openapiv2_out . ^
    --openapiv2_opt logtostderr=true ^
    cmd/proto/auth.proto

protoc -I . ^
    -I %GOPATH%/pkg/mod/github.com/googleapis/googleapis@v0.0.0-20240605231841-21c839628513 ^
    -I %GOPATH%/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.20.0 ^
    --go_out . ^
    --go_opt paths=import ^
    --go-grpc_out . ^
    --go-grpc_opt paths=import ^
    --grpc-gateway_out . ^
    --grpc-gateway_opt paths=import ^
    --grpc-gateway_opt logtostderr=true ^
    --grpc-gateway_opt generate_unbound_methods=true ^
    --openapiv2_out . ^
    --openapiv2_opt logtostderr=true ^
    cmd/proto/gophkeeper.proto


go-swagger-merger -o internal/server/swagger/swagger.json  cmd/proto/auth.swagger.json cmd/proto/gophkeeper.swagger.json 