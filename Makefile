.PHONY: build


build:
	docker-compose build

run:
	docker-compose run --rm --entrypoint go app run cmd/main.go
