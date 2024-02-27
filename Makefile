generate:
	protoc -I proto proto/investments.proto --go_out=./gen --go_opt=paths=source_relative --go-grpc_out=./gen --go-grpc_opt=paths=source_relative
migrate-up:
	migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5430/postgres?sslmode=disable" up
migrate-down:
	migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5430/postgres?sslmode=disable" down
