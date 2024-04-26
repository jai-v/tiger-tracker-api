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
    export DB_ROOT_PASSWORD=pasword
    export WEBAPP_HOST=localhost
    export WEBAPP_PORT=5001
    export MIGRATION_CONTEXT=local
    export HYDRA_DB_USER=hydrauser
    export HYDRA_DB_PASSWORD=secret
    export HYDRA_DB_NAME=hydradb
    export HYDRA_DB_PORT=5432
    export HYDRA_URLS_LOGIN=http://$WEBAPP_HOST:$WEBAPP_PORT/authentication/login
    export HYDRA_URLS_CONSENT=http://$WEBAPP_HOST:$WEBAPP_PORT/authentication/consent
    export HYDRA_SERVE_PUBLIC_HOST=0.0.0.0
    export HYDRA_HOST=localhost
    export HYDRA_PUBLIC_PORT=4444
    export HYDRA_ADMIN_PORT=4445
    export ADMINER_PORT=9000
    export HYDRA_CLIENT_ID=ui-web
    export HYDRA_CLIENT_SECRET=topsecret
    

## Commands:

### To run the WebApp with API Service & Hydra
    make run

    #run the below docker command (one time activity)
    docker-compose exec hydra hydra clients create \
    --endpoint http://127.0.0.1:4445/ \
    --name Tiger Tracker \
    --id ui-web \
    --secret topsecret \
    --grant-types authorization_code,refresh_token \
    --response-types code,id_token \
    --callbacks http://localhost:5001/callback \
    --token-endpoint-auth-method client_secret_post \
    --scope offline

    #open below url in a browser:
    http://localhost:5001

    #login with
    #username: testuser
    #password: testuser

### To stop all services
    make stop

### To start only database
    make start_db

### To run the migrations
    make migrate_db


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

