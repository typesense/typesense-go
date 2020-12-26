//go:generate oapi-codegen -package api -generate types -o types_gen.go openapi.yml
//go:generate oapi-codegen -package api -generate client -o client_gen.go openapi.yml

package api
