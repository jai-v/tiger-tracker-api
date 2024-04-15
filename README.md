# tiger-tracker-api
Microservice for tracking the population of tigers in the wild

## Prerequistes:
    go 1.22
    docker

## Set below environment variables:
    export PATH=$(go env GOPATH)/bin:$PATH
    export GO111MODULE=on
    export DB_HOST=localhost
    export DB_PORT=3306
    export DB_USER=tiger
    export DB_PASSWORD=kitten
    export DB_NAME=tigerhall
    export MIGRATION_CONTEXT=local

## Commands:

### To start db
    make start_db

### To run migrations
    make migrate_db

### To run service
    make run

### To generate swagger docs
    //Install Swagger
    go get -u github.com/swaggo/swag/cmd/swag
    
    //Generate Swagger Docs
    swag init -g router/router.go

### Swagger URLs

| Environment | URL                                                           |
|-------------|---------------------------------------------------------------|
| Local       | http://localhost:8080/api/tiger-tracker/v1/swagger/index.html |

### To generate or update mocks
    //Install mockgen
    go install go.uber.org/mock/mockgen@latest

