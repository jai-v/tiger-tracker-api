start_db:
	docker-compose -f ./docker-compose.yml up db --wait

make migrate_db:
	docker-compose -f ./docker-compose.yml up db-migrate