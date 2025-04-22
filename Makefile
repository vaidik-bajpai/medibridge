SCHEMA_PATH = ./internal/prisma/schema.prisma

prisma-gen:
	go run github.com/steebchen/prisma-client-go generate --schema $(SCHEMA_PATH)
 
prisma-push:
	go run github.com/steebchen/prisma-client-go db push --schema $(SCHEMA_PATH)
 
postgres-shell: 
	docker exec -it medibridgeDB psql -U root -d postgres

