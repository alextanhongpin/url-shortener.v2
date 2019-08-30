-include .env
export

DATABASE_URL := postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

start: generate
	@go run main.go

init:
	@go mod init

tidy:
	@brew install golang-migrate
	@go mod tidy

up:
	@docker-compose up -d

down:
	@docker-compose down

clean:
	@rm -rf ./tmp

db:
	@psql -h ${DB_HOST} -U ${DB_USER} -d ${DB_NAME}

sql-%:
	@migrate create -ext sql -dir infra/migrations -seq $* 

migrate:
	@migrate -database ${DATABASE_URL} -path infra/migrations up

rollback:
	@migrate -database ${DATABASE_URL} -path infra/migrations drop
	@migrate -database ${DATABASE_URL} -path infra/migrations down 

generate:
	@go-bindata -prefix 'infra/migrations/' -pkg migrations -o infra/migrations/bindata.go infra/migrations/...
