delete-docker-build-cache:
	docker buildx prune -f

delete-unused-docker-images:
	docker rmi $(docker images --filter "dangling=true" -q --no-trunc)

run-docker-all:
	docker compose -f docker-compose-nats.yml up -d
	docker compose -f docker-compose-user.yml up -d
	docker compose -f docker-compose-bank.yml up -d

rebuild-user:
	docker compose -f docker-compose-user.yml up -d --build user_service

rebuild-bank:
	docker compose -f docker-compose-bank.yml up -d --build bank_service

stop-docker-all:
	docker compose -f docker-compose-nats.yml down
	docker compose -f docker-compose-user.yml down
	docker compose -f docker-compose-bank.yml down