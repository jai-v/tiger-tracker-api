PROJECT := "tiger-tracker-api"
start_db:
	docker-compose --project-name $(PROJECT) -f ./docker-compose.yml up mysqldb --wait

make migrate_db:
	docker-compose --project-name $(PROJECT) -f ./docker-compose.yml up mysqldb-migrate

make run:
	docker-compose --project-name $(PROJECT) -f ./docker-compose.yml up tiger-tracker-api --build --force-recreate

make stop:
	docker-compose --project-name $(PROJECT) -f ./docker-compose.yml down
