.PHONY: postgres adminer migrate

postgres:
	docker run --rm -ti -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=postgres -e POSTGRES_HOST_AUTH_METHOD=md5 postgres

adminer:
	docker run --rm -ti -p 5555:8080 adminer

migrate:
	migrate -source file://migrations \
		-database postgres://postgres:password@localhost:5432/postgres?sslmode=disable \
		up

migrate-down:
	migrate -source file://migrations \
		-database postgres://postgres:password@localhost:5432/postgres?sslmode=disable \
		down