redis:
	docker run -d -p 6379:6379 redis:alpine3.17

proto:
	rm -f grpc/gen/*.go
	protoc --proto_path=grpc/porto --go_out=grpc/pb --go_opt=paths=source_relative \
	--go-grpc_out=grpc/pb --go-grpc_opt=paths=source_relative \
	grpc/porto/*.proto

evans:
	evans --host 127.0.0.1 --port 9090 -r repl