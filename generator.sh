go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -package api -generate client -o $(pwd)/typesense/api/client_gen.go $(pwd)/typesense/api/generator/generator.yml
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -package api -generate types -o $(pwd)/typesense/api/types_gen.go $(pwd)/typesense/api/generator/generator.yml
