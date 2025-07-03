include .env

MIGRATIONS_DIR= services/room-svc/db/migrations


DB_URL = "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&TimeZone=Asia/Jakarta"

CQL_DB_URL = "host=localhost port=9042 keyspace=messaging_service"
migrate-up :
	goose -dir ${MIGRATIONS_DIR} postgres ${DB_URL} up

migrate-down :
	goose -dir ${MIGRATIONS_DIR} postgres ${DB_URL} down

migrate-down-to-zero :
	goose -dir ${MIGRATIONS_DIR} postgres ${DB_URL} down-to 0

migrate-cql-up :
	cd services/chat-command-cql-svc/cmd/migration && go run main.go

generate-proto-user:
	cd proto && protoc --go_out=. --go-grpc_out=. user.proto
	
generate-proto-room:
	cd proto && protoc --go_out=. --go-grpc_out=. room.proto

generate-proto-query:
	cd proto && protoc --go_out=. --go-grpc_out=. query.proto

start-user-svc:
	cd services/user-svc/cmd/web && go run main.go

start-room-svc:
	cd services/room-svc/cmd/web && go run main.go

start-chat-query-svc:
	cd services/chat-query-svc/cmd/web && go run main.go

start-chat-realtime-svc:
	cd services/chat-realtime-svc/cmd/web && go run main.go

start-chat-ingestion-svc:
	cd services/chat-ingestion-svc/cmd/worker && go run main.go

start-chat-command-cql-svc:
	cd services/chat-command-cql-svc/cmd/worker && go run main.go

start-chat-command-elastic-svc:
	cd services/chat-command-elastic-svc/cmd/worker && go run main.go