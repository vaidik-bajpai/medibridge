SCHEMA_PATH = ./internal/prisma/schema.prisma

prisma-gen:
	go run github.com/steebchen/prisma-client-go generate --schema $(SCHEMA_PATH)
 
prisma-push:
	go run github.com/steebchen/prisma-client-go db push --schema $(SCHEMA_PATH)
 
postgres-shell: 
	docker exec -it my-postgres psql -U postgres -d postgres

mock-gen:
	mockery --name=UserStorer --output=internal/mocks --outpkg=mocks --dir=internal/store && \
	mockery --name=PatientStorer --output=internal/mocks --outpkg=mocks --dir=internal/store && \
	mockery --name=SessionStorer --output=internal/mocks --outpkg=mocks --dir=internal/store && \
	mockery --name=DiagnosesStorer --output=internal/mocks --outpkg=mocks --dir=internal/store && \
	mockery --name=VitalsStorer --output=internal/mocks --outpkg=mocks --dir=internal/store && \
	mockery --name=ConditionStorer --output=internal/mocks --outpkg=mocks --dir=internal/store && \
	mockery --name=AllergyStorer --output=internal/mocks --outpkg=mocks --dir=internal/store

db:
	docker start my-postgres

swag-gen:
	swag init -g cmd/main.go