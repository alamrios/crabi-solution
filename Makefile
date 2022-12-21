.PHONY: test

serve:
	docker-compose -f infra/deploy/local/docker-compose.yml up -d --build

stop:
	docker-compose -f infra/deploy/local/docker-compose.yml down

test:
	go test -cover \
	./internal/app/user \
	./internal/infra/http/pld

proto-compile:
	protoc --proto_path=proto/ pld/v1/pld.proto --go_out=proto/pld --go-grpc_out=proto/pld