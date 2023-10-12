package main

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 -config ./oapi-codegen.yml openapi.yaml
//go:generate go run github.com/kyleconroy/sqlc/cmd/sqlc@f1eef01 generate
//go:generate go run github.com/jmattheis/goverter/cmd/goverter@v0.17.4 -output ./conv/convgen/convgen.go -packageName convgen ./conv
