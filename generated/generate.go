package generated

//go:generate oapi-codegen -generate types -package openapi -o ./openapi/types.go ../api/openapi.yaml
//go:generate oapi-codegen -generate echo -package openapi -o ./openapi/server.go ../api/openapi.yaml
