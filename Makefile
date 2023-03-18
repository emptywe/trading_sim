PROJECTNAME=$(shell basename "$(PWD)")

.PHONY:
.SILENT:

CFLAGS = -local

build:
	go build -o ./.bin/sim cmd/*.go
compose:
	docker-compose up -d --build --remove-orphans
run: build
	./.bin/sim ${CFLAGS}
down:
	docker-compose down
show_logs:
	docker-compose logs

run-docker: compose

