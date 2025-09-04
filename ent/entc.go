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
		log.Fatalf("creating entgql extension: %v", err)
	}
	extensions := []entc.Extension{ex}
	extensions = append(extensions, entviz.Extension{})
	opts := []entc.Option{
		entc.Extensions(extensions...),
		entc.TemplateDir("./ent/template"),
		entc.FeatureNames("intercept"),
	}
	if err := os.RemoveAll("./ent/gen"); err != nil {
		log.Fatalf("failed to remove ./ent/gen: %v", err)
	}
	if err := entc.Generate("./ent/schema", &gen.Config{
		Target:  "./ent/gen/ent",
		Package: "github.com/buildbarn/bb-portal/ent/gen/ent",
	}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
