.PHONY: build


build:
	docker-compose build

run:
	docker-compose run --rm app
