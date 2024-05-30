docker-start:
	docker-compose -f ./docker-compose.yml --profile infra --profile app up -d