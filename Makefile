# to run server
.PHONY: runServerLinux
runServerLinux:
	cd cmd && GOOS=linux GOARCH=amd64 go build -o main

.PHONY: runServerMac
runServerMac:
	cd cmd && go build -o main && ./main

.PHONY: migrateUp
migrateUp:
	migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)"  -path db/migrations up

# to run rollback migration
.PHONY: migrateDown
migrateDown:
	migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)"  -path db/migrations down
	