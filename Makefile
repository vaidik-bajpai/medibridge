SCHEMA_PATH = ./internal/prisma/schema.prisma

prisma-gen:
	go run github.com/steebchen/prisma-client-go generate --schema $(SCHEMA_PATH)
 
prisma-push:
	go run github.com/steebchen/prisma-client-go db push --schema $(SCHEMA_PATH)
 
postgres-shell: 
	docker exec -it medibridgeDB psql -U root -d postgres

mock-gen:
	mockery --name=UserStorer --output=./mocks --outpkg=mocks /
	mockery --name=PatientStorer --output=./mocks --outpkg=mocks /
	mockery --name=SessionStorer --output=./mocks --outpkg=mocks

db:
	docker start medibridgeDB

swag-gen:
	swagger generate spec -o ./swagger.yaml –scan-models