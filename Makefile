PROJECT := "tiger-tracker-api"
start_db:
	docker compose --project-name $(PROJECT) -f ./docker-compose.yml up mysqldb --wait

make migrate_db:
	docker compose --project-name $(PROJECT) -f ./docker-compose.yml up mysqldb-migrate

make run:
	docker compose --project-name $(PROJECT) -f ./docker-compose.yml up --build --force-recreate

make stop:
	docker compose --project-name $(PROJECT) -f ./docker-compose.yml down
