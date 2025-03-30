WORK_DIR = $(shell pwd)

PROJECT := scan-service
REVISION := latest

BUILD_VENDOR := go mod vendor && chmod -R +w vendor

install_deps:
	docker compose -f infrastructure/build.yaml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "apk update && apk add git && $(BUILD_VENDOR)"

build: install_deps
	docker compose -f infrastructure/build.yaml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go build -mod=vendor -o ./bin/$(PROJECT)"

vet: install_deps
	docker-compose -f infrastructure/build.yaml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go vet -mod=vendor ./..."

create_user:
	docker exec -it scanner-db bash -c "/sql/create_user.sh"

start: build
	docker-compose -f docker-compose.local-app.yml up -d

stop:
	docker-compose -f docker-compose.local-app.yml down -v

clean:
	chmod -R +w ./.gopath vendor || true
