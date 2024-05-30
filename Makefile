docker-start:
	docker-compose -f ./docker-compose.yml --profile infra up -d

config-local:
	./config.sh local