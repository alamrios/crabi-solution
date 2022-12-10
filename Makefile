
serve:
	docker-compose -f infra/deploy/local/docker-compose.yml up -d --build

stop:
	docker-compose -f infra/deploy/local/docker-compose.yml down