PROJECTNAME=$(shell basename "$(PWD)")

.PHONY:
.SILENT:

CFLAGS = -local

build:
	go build -o ./.bin/sim cmd/*.go
compose:
	docker-compose up --build
run: build
	./.bin/sim ${CFLAGS}

run-docker: compose


