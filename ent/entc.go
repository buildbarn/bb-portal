//go:build entc

package main

import (
	"log"
	"os"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/hedwigz/entviz"
)

func main() {
	ex, err := entgql.NewExtension(
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("./internal/graphql/schema/ent.graphql"),
		entgql.WithWhereInputs(true),
		entgql.WithConfigPath("./gqlgen.yml"),
		entgql.WithRelaySpec(true),
		entgql.WithNodeDescriptor(false),
	)
	if err != nil {
		log.Fatalf("Error creating entgql extension: %v", err)
	}
	config := &gen.Config{
		Target:  "./ent/gen/ent",
		Package: "github.com/buildbarn/bb-portal/ent/gen/ent",
	}
	opts := []entc.Option{
		entc.Extensions(
			ex,
			entviz.Extension{},
		),
		entc.FeatureNames("sql/execquery"),
		entc.FeatureNames("sql/upsert"),
		entc.FeatureNames("sql/modifier"),
		entc.FeatureNames("privacy"),
		entc.FeatureNames("entql"),
	}
	err = os.RemoveAll("./ent/gen/ent")
	if err != nil {
		log.Fatalf("Error removing old generated code: %v", err)
	}
	if err := entc.Generate("./ent/schema", config, opts...); err != nil {
		log.Fatalf("Error running first ent codegen: %v", err)
	}
	if err := entc.Generate("./ent/authschema", config, opts...); err != nil {
		log.Fatalf("Error running second ent codegen: %v\nMake sure that all schemas have been added to ent/authschema/schema.go correctly. Read `ent/README.md` to understand how the generation system works.", err)
	}
}
