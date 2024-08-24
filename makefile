PROTOC = protoc
PROTO_PATH = api/proto
PROTO_FILES_MASK = ${PROTO_PATH}/*.proto

proto:
	@${PROTOC} --proto_path=${PROTO_PATH} --go_out=. ${PROTO_FILES_MASK}
	@${PROTOC} --proto_path=${PROTO_PATH} --go-grpc_out=. ${PROTO_FILES_MASK}
	@echo "the protobuf files have been rebuild"


MIGRATE = migrate
MIGR_DIR = migrations
DB_SOURCE = postgres://amicie:admin@localhost:5432/$(1)?sslmode=disable
MIGRATE_BODY = ${MIGRATE} -path ${MIGR_DIR} -database $(call DB_SOURCE,$(db))

migrate: check_db
	@${MIGRATE_BODY} up

migrate_force: check_db check_v
	@${MIGRATE_BODY} force ${v}

migrate_down:
	@${MIGRATE_BODY} down

check_db:
ifndef db
	@$(error parameter db is required [database name])
endif

check_v:
ifndef v
	@$(error parameter v is required [version])
endif

ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif
