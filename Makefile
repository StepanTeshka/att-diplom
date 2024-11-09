include .env

LOCAL_BIN := $(CURDIR)/LDE/bin

# Dependencies ===============

install:
	mkdir -p $(LOCAL_BIN)
	make install-go-air-livereload


install-go-air-livereload:
	GOBIN=$(LOCAL_BIN) go install github.com/cosmtrek/air@v1.51.0

# ============================


# Go Dependencies ============

t: 
	go mod tidy

# ============================




# Main =======================

run:
	go run ./cmd/server/main.go

run-watch:
	$(LOCAL_BIN)/air -c .air.toml

# ============================




DOCKER_APP_NAME:=skymine

docker-run:
	docker compose -f ./docker-compose.yml --env-file ./.env -p ${DOCKER_APP_NAME} up -d

docker-stop:
	docker compose -f ./docker-compose.yml --env-file ./.env -p ${DOCKER_APP_NAME} down



# Make sys ===================

_required_param-%:
	@if [ "${${*}}" == "" ]; then echo "\n\033[0;91mPlease provide arg: \"$*\"\033[0m\n"; exit 1; fi

# ============================