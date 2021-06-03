start:
	docker-compose up upload-api

stop:
	docker-compose down

tests:
	docker-compose up tests