package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"context"
	"fmt"

	"entgo.io/contrib/entgql"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocationproblem"
	"github.com/buildbarn/bb-portal/ent/gen/ent/blob"
	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/ent/gen/ent/eventfile"
	"github.com/buildbarn/bb-portal/internal/graphql/helpers"
)

// ID is the resolver for the id field.
func (r *bazelInvocationResolver) ID(ctx context.Context, obj *ent.BazelInvocation) (string, error) {
	return helpers.GraphQLIDFromTypeAndID("BazelInvocation", obj.ID), nil
}

// ID is the resolver for the id field.
func (r *bazelInvocationProblemResolver) ID(ctx context.Context, obj *ent.BazelInvocationProblem) (string, error) {
	return helpers.GraphQLIDFromTypeAndID("BazelInvocationProblem", obj.ID), nil
}

// ID is the resolver for the id field.
func (r *blobResolver) ID(ctx context.Context, obj *ent.Blob) (string, error) {
	return helpers.GraphQLIDFromTypeAndID("Blob", obj.ID), nil
}

// ID is the resolver for the id field.
func (r *buildResolver) ID(ctx context.Context, obj *ent.Build) (string, error) {
	return helpers.GraphQLIDFromTypeAndID("Build", obj.ID), nil
}

// ID is the resolver for the id field.
func (r *eventFileResolver) ID(ctx context.Context, obj *ent.EventFile) (string, error) {
	return helpers.GraphQLIDFromTypeAndID("EventFile", obj.ID), nil
}

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id string) (ent.Noder, error) {
	typ, intID, err := helpers.GraphQLTypeAndIntIDFromID(id)
	if err != nil {
		return nil, err
	}
	table := map[string]string{
		"BazelInvocation":        bazelinvocation.Table,
		"BazelInvocationProblem": bazelinvocationproblem.Table,
		"ActionProblem":          bazelinvocationproblem.Table,
		"ProgressProblem":        bazelinvocationproblem.Table,
		"TargetProblem":          bazelinvocationproblem.Table,
		"TestProblem":            bazelinvocationproblem.Table,
		"Blob":                   blob.Table,
		"Build":                  build.Table,
		"EventFile":              eventfile.Table,
	}[typ]

	var n ent.Noder
	if table == bazelinvocationproblem.Table {
		var dbProblem *ent.BazelInvocationProblem
		dbProblem, err = r.client.BazelInvocationProblem.Get(ctx, intID)
		if err != nil {
			return nil, err
		}
		n, err = r.helper.DBProblemToAPIProblem(ctx, dbProblem)
		if err != nil {
			return nil, err
		}
	} else {
		n, err = r.client.Noder(ctx, intID, ent.WithFixedNodeType(table))
	}

	if err != nil {
		return nil, err
	}
	return n, nil
}

// Nodes is the resolver for the nodes field.
func (r *queryResolver) Nodes(ctx context.Context, ids []string) ([]ent.Noder, error) {
	noders := make([]ent.Noder, 0, len(ids))
	for _, id := range ids {
		noder, err := r.Node(ctx, id)
		if err != nil {
			return nil, err
		}
		noders = append(noders, noder)
	}
	return noders, nil
}

// FindBazelInvocations is the resolver for the findBazelInvocations field.
func (r *queryResolver) FindBazelInvocations(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int, where *ent.BazelInvocationWhereInput) (*ent.BazelInvocationConnection, error) {
	return r.client.BazelInvocation.Query().Paginate(ctx, after, first, before, last, ent.WithBazelInvocationFilter(where.Filter))
}

// FindBuilds is the resolver for the findBuilds field.
func (r *queryResolver) FindBuilds(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int, where *ent.BuildWhereInput) (*ent.BuildConnection, error) {
	return r.client.Build.Query().Paginate(ctx, after, first, before, last, ent.WithBuildFilter(where.Filter))
}

// ID is the resolver for the id field.
func (r *bazelInvocationProblemWhereInputResolver) ID(ctx context.Context, obj *ent.BazelInvocationProblemWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: ID - id"))
}

// IDNeq is the resolver for the idNEQ field.
func (r *bazelInvocationProblemWhereInputResolver) IDNeq(ctx context.Context, obj *ent.BazelInvocationProblemWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDNeq - idNEQ"))
}

// IDIn is the resolver for the idIn field.
func (r *bazelInvocationProblemWhereInputResolver) IDIn(ctx context.Context, obj *ent.BazelInvocationProblemWhereInput, data []string) error {
	panic(fmt.Errorf("not implemented: IDIn - idIn"))
}

// IDNotIn is the resolver for the idNotIn field.
func (r *bazelInvocationProblemWhereInputResolver) IDNotIn(ctx context.Context, obj *ent.BazelInvocationProblemWhereInput, data []string) error {
	panic(fmt.Errorf("not implemented: IDNotIn - idNotIn"))
}

// IDGt is the resolver for the idGT field.
func (r *bazelInvocationProblemWhereInputResolver) IDGt(ctx context.Context, obj *ent.BazelInvocationProblemWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDGt - idGT"))
}

// IDGte is the resolver for the idGTE field.
func (r *bazelInvocationProblemWhereInputResolver) IDGte(ctx context.Context, obj *ent.BazelInvocationProblemWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDGte - idGTE"))
}

// IDLt is the resolver for the idLT field.
func (r *bazelInvocationProblemWhereInputResolver) IDLt(ctx context.Context, obj *ent.BazelInvocationProblemWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDLt - idLT"))
}

// IDLte is the resolver for the idLTE field.
func (r *bazelInvocationProblemWhereInputResolver) IDLte(ctx context.Context, obj *ent.BazelInvocationProblemWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDLte - idLTE"))
}

// ID is the resolver for the id field.
func (r *bazelInvocationWhereInputResolver) ID(ctx context.Context, obj *ent.BazelInvocationWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: ID - id"))
}

// IDNeq is the resolver for the idNEQ field.
func (r *bazelInvocationWhereInputResolver) IDNeq(ctx context.Context, obj *ent.BazelInvocationWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDNeq - idNEQ"))
}

// IDIn is the resolver for the idIn field.
func (r *bazelInvocationWhereInputResolver) IDIn(ctx context.Context, obj *ent.BazelInvocationWhereInput, data []string) error {
	panic(fmt.Errorf("not implemented: IDIn - idIn"))
}

// IDNotIn is the resolver for the idNotIn field.
func (r *bazelInvocationWhereInputResolver) IDNotIn(ctx context.Context, obj *ent.BazelInvocationWhereInput, data []string) error {
	panic(fmt.Errorf("not implemented: IDNotIn - idNotIn"))
}

// IDGt is the resolver for the idGT field.
func (r *bazelInvocationWhereInputResolver) IDGt(ctx context.Context, obj *ent.BazelInvocationWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDGt - idGT"))
}

// IDGte is the resolver for the idGTE field.
func (r *bazelInvocationWhereInputResolver) IDGte(ctx context.Context, obj *ent.BazelInvocationWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDGte - idGTE"))
}

// IDLt is the resolver for the idLT field.
func (r *bazelInvocationWhereInputResolver) IDLt(ctx context.Context, obj *ent.BazelInvocationWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDLt - idLT"))
}

// IDLte is the resolver for the idLTE field.
func (r *bazelInvocationWhereInputResolver) IDLte(ctx context.Context, obj *ent.BazelInvocationWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDLte - idLTE"))
}

// ID is the resolver for the id field.
func (r *blobWhereInputResolver) ID(ctx context.Context, obj *ent.BlobWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: ID - id"))
}

// IDNeq is the resolver for the idNEQ field.
func (r *blobWhereInputResolver) IDNeq(ctx context.Context, obj *ent.BlobWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDNeq - idNEQ"))
}

// IDIn is the resolver for the idIn field.
func (r *blobWhereInputResolver) IDIn(ctx context.Context, obj *ent.BlobWhereInput, data []string) error {
	panic(fmt.Errorf("not implemented: IDIn - idIn"))
}

// IDNotIn is the resolver for the idNotIn field.
func (r *blobWhereInputResolver) IDNotIn(ctx context.Context, obj *ent.BlobWhereInput, data []string) error {
	panic(fmt.Errorf("not implemented: IDNotIn - idNotIn"))
}

// IDGt is the resolver for the idGT field.
func (r *blobWhereInputResolver) IDGt(ctx context.Context, obj *ent.BlobWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDGt - idGT"))
}

// IDGte is the resolver for the idGTE field.
func (r *blobWhereInputResolver) IDGte(ctx context.Context, obj *ent.BlobWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDGte - idGTE"))
}

// IDLt is the resolver for the idLT field.
func (r *blobWhereInputResolver) IDLt(ctx context.Context, obj *ent.BlobWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDLt - idLT"))
}

// IDLte is the resolver for the idLTE field.
func (r *blobWhereInputResolver) IDLte(ctx context.Context, obj *ent.BlobWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDLte - idLTE"))
}

// ID is the resolver for the id field.
func (r *buildWhereInputResolver) ID(ctx context.Context, obj *ent.BuildWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: ID - id"))
}

// IDNeq is the resolver for the idNEQ field.
func (r *buildWhereInputResolver) IDNeq(ctx context.Context, obj *ent.BuildWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDNeq - idNEQ"))
}

// IDIn is the resolver for the idIn field.
func (r *buildWhereInputResolver) IDIn(ctx context.Context, obj *ent.BuildWhereInput, data []string) error {
	panic(fmt.Errorf("not implemented: IDIn - idIn"))
}

// IDNotIn is the resolver for the idNotIn field.
func (r *buildWhereInputResolver) IDNotIn(ctx context.Context, obj *ent.BuildWhereInput, data []string) error {
	panic(fmt.Errorf("not implemented: IDNotIn - idNotIn"))
}

// IDGt is the resolver for the idGT field.
func (r *buildWhereInputResolver) IDGt(ctx context.Context, obj *ent.BuildWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDGt - idGT"))
}

// IDGte is the resolver for the idGTE field.
func (r *buildWhereInputResolver) IDGte(ctx context.Context, obj *ent.BuildWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDGte - idGTE"))
}

// IDLt is the resolver for the idLT field.
func (r *buildWhereInputResolver) IDLt(ctx context.Context, obj *ent.BuildWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDLt - idLT"))
}

// IDLte is the resolver for the idLTE field.
func (r *buildWhereInputResolver) IDLte(ctx context.Context, obj *ent.BuildWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDLte - idLTE"))
}

// ID is the resolver for the id field.
func (r *eventFileWhereInputResolver) ID(ctx context.Context, obj *ent.EventFileWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: ID - id"))
}

// IDNeq is the resolver for the idNEQ field.
func (r *eventFileWhereInputResolver) IDNeq(ctx context.Context, obj *ent.EventFileWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDNeq - idNEQ"))
}

// IDIn is the resolver for the idIn field.
func (r *eventFileWhereInputResolver) IDIn(ctx context.Context, obj *ent.EventFileWhereInput, data []string) error {
	panic(fmt.Errorf("not implemented: IDIn - idIn"))
}

// IDNotIn is the resolver for the idNotIn field.
func (r *eventFileWhereInputResolver) IDNotIn(ctx context.Context, obj *ent.EventFileWhereInput, data []string) error {
	panic(fmt.Errorf("not implemented: IDNotIn - idNotIn"))
}

// IDGt is the resolver for the idGT field.
func (r *eventFileWhereInputResolver) IDGt(ctx context.Context, obj *ent.EventFileWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDGt - idGT"))
}

// IDGte is the resolver for the idGTE field.
func (r *eventFileWhereInputResolver) IDGte(ctx context.Context, obj *ent.EventFileWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDGte - idGTE"))
}

// IDLt is the resolver for the idLT field.
func (r *eventFileWhereInputResolver) IDLt(ctx context.Context, obj *ent.EventFileWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDLt - idLT"))
}

// IDLte is the resolver for the idLTE field.
func (r *eventFileWhereInputResolver) IDLte(ctx context.Context, obj *ent.EventFileWhereInput, data *string) error {
	panic(fmt.Errorf("not implemented: IDLte - idLTE"))
}

// BazelInvocation returns BazelInvocationResolver implementation.
func (r *Resolver) BazelInvocation() BazelInvocationResolver { return &bazelInvocationResolver{r} }

// BazelInvocationProblem returns BazelInvocationProblemResolver implementation.
func (r *Resolver) BazelInvocationProblem() BazelInvocationProblemResolver {
	return &bazelInvocationProblemResolver{r}
}

// Blob returns BlobResolver implementation.
func (r *Resolver) Blob() BlobResolver { return &blobResolver{r} }

// Build returns BuildResolver implementation.
func (r *Resolver) Build() BuildResolver { return &buildResolver{r} }

// EventFile returns EventFileResolver implementation.
func (r *Resolver) EventFile() EventFileResolver { return &eventFileResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// BazelInvocationProblemWhereInput returns BazelInvocationProblemWhereInputResolver implementation.
func (r *Resolver) BazelInvocationProblemWhereInput() BazelInvocationProblemWhereInputResolver {
	return &bazelInvocationProblemWhereInputResolver{r}
}

// BazelInvocationWhereInput returns BazelInvocationWhereInputResolver implementation.
func (r *Resolver) BazelInvocationWhereInput() BazelInvocationWhereInputResolver {
	return &bazelInvocationWhereInputResolver{r}
}

// BlobWhereInput returns BlobWhereInputResolver implementation.
func (r *Resolver) BlobWhereInput() BlobWhereInputResolver { return &blobWhereInputResolver{r} }

// BuildWhereInput returns BuildWhereInputResolver implementation.
func (r *Resolver) BuildWhereInput() BuildWhereInputResolver { return &buildWhereInputResolver{r} }

// EventFileWhereInput returns EventFileWhereInputResolver implementation.
func (r *Resolver) EventFileWhereInput() EventFileWhereInputResolver {
	return &eventFileWhereInputResolver{r}
}

type bazelInvocationResolver struct{ *Resolver }
type bazelInvocationProblemResolver struct{ *Resolver }
type blobResolver struct{ *Resolver }
type buildResolver struct{ *Resolver }
type eventFileResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type bazelInvocationProblemWhereInputResolver struct{ *Resolver }
type bazelInvocationWhereInputResolver struct{ *Resolver }
type blobWhereInputResolver struct{ *Resolver }
type buildWhereInputResolver struct{ *Resolver }
type eventFileWhereInputResolver struct{ *Resolver }
