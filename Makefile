CURRENT_DIR=$(shell pwd)

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto-gen:
	./scripts/genproto.sh	${CURRENT_DIR}





DB_URL := "postgres://postgres:4444@localhost:5432/examdb?sslmode=disable"



migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up
migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down
migrate_file:
	migrate create -ext sql -dir migrations/ -seq  users
migrate_dirty:
	migrate -path ./migrations/ -database "postgres://postgres:4444@localhost:5432/examdb?sslmode=disable" force <No>



