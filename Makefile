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

# to run migration file
.PHONY: migrateUpProd
migrateUpProd:
	migrate -database "postgres://postgres:ohN6Nei0ugiRena5@project-sprint-postgres.cavsdeuj9ixh.ap-southeast-1.rds.amazonaws.com:5432/postgres?sslmode=verify-full&sslrootcert=ap-southeast-1-bundle.pem"  -path db/migrations up

# to run rollback migration
.PHONY: migrateDownProd
migrateDown:
	migrate -database "postgres://postgres:ohN6Nei0ugiRena5@project-sprint-postgres.cavsdeuj9ixh.ap-southeast-1.rds.amazonaws.com:5432/postgres?sslmode=disable" -path db/migrations down