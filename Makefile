-include .env
export

VERSION := $(shell git rev-parse --short HEAD)
DATABASE_URL := postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable
MIGRATION_PATH := database/migrations
IMG := alextanhongpin/url-shortener

start: generate up
	@go run main.go

init:
	@go mod init

install:
	@go get -u github.com/jteeuwen/go-bindata/...
	@go mod tidy
	@brew install golang-migrate

up:
	@docker-compose up -d

down:
	@docker-compose down

clean: # Clears the local database volume.
	@rm -rf ./tmp

db: # Access the postgres cli.
	@psql -h ${DB_HOST} -U ${DB_USER} -d ${DB_NAME}

sql-%: # Creates a new migration file.
	@migrate create -ext sql -dir ${MIGRATION_PATH} -seq $* 

migrate: # Run the database migration manually.
	@migrate -database ${DATABASE_URL} -path ${MIGRATION_PATH} up

rollback: # Rollback the database migration manually.
	@migrate -database ${DATABASE_URL} -path ${MIGRATION_PATH} drop
	@migrate -database ${DATABASE_URL} -path ${MIGRATION_PATH} down 

generate: # Generate the binary data for the database migration.
	@go-bindata -prefix ${MIGRATION_PATH} -pkg migrations -o ${MIGRATION_PATH}/bindata.go ${MIGRATION_PATH}/...

docker: generate
	@docker build --build-arg VERSION=${VERSION} -t ${IMG} .

docker-start:
	@docker run -d -p 8080:8080 ${IMG}
