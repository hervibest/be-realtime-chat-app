include .env

DB_URL = "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&TimeZone=Asia/Jakarta"
migrate-up :
	goose -dir ${MIGRATIONS_DIR} postgres ${DB_URL} up
migrate-down :
	goose -dir ${MIGRATIONS_DIR} postgres ${DB_URL} down
migrate-down-to-zero :
	goose -dir ${MIGRATIONS_DIR} postgres ${DB_URL} down-to 0