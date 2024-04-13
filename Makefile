start_db:
	docker-compose -f ./docker-compose.yml up mysqldb --wait

make migrate_db:
	docker-compose -f ./docker-compose.yml up mysqldb-migrate

make run:
	docker-compose -f ./docker-compose.yml up tiger-tracker-api