// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocationproblem"
	"github.com/buildbarn/bb-portal/ent/gen/ent/blob"
	"github.com/buildbarn/bb-portal/ent/gen/ent/build"
	"github.com/buildbarn/bb-portal/ent/gen/ent/eventfile"
)

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (bi *BazelInvocationQuery) CollectFields(ctx context.Context, satisfies ...string) (*BazelInvocationQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return bi, nil
	}
	if err := bi.collectField(ctx, false, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return bi, nil
}

func (bi *BazelInvocationQuery) collectField(ctx context.Context, oneNode bool, opCtx *graphql.OperationContext, collected graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	var (
		unknownSeen    bool
		fieldSeen      = make(map[string]struct{}, len(bazelinvocation.Columns))
		selectedFields = []string{bazelinvocation.FieldID}
	)
	for _, field := range graphql.CollectFields(opCtx, collected.Selections, satisfies) {
		switch field.Name {

		case "eventFile":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = (&EventFileClient{config: bi.config}).Query()
			)
			if err := query.collectField(ctx, oneNode, opCtx, field, path, mayAddCondition(satisfies, eventfileImplementors)...); err != nil {
				return err
			}
			bi.withEventFile = query

		case "build":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = (&BuildClient{config: bi.config}).Query()
			)
			if err := query.collectField(ctx, oneNode, opCtx, field, path, mayAddCondition(satisfies, buildImplementors)...); err != nil {
				return err
			}
			bi.withBuild = query
		case "invocationID":
			if _, ok := fieldSeen[bazelinvocation.FieldInvocationID]; !ok {
				selectedFields = append(selectedFields, bazelinvocation.FieldInvocationID)
				fieldSeen[bazelinvocation.FieldInvocationID] = struct{}{}
			}
		case "startedAt":
			if _, ok := fieldSeen[bazelinvocation.FieldStartedAt]; !ok {
				selectedFields = append(selectedFields, bazelinvocation.FieldStartedAt)
				fieldSeen[bazelinvocation.FieldStartedAt] = struct{}{}
			}
		case "endedAt":
			if _, ok := fieldSeen[bazelinvocation.FieldEndedAt]; !ok {
				selectedFields = append(selectedFields, bazelinvocation.FieldEndedAt)
				fieldSeen[bazelinvocation.FieldEndedAt] = struct{}{}
			}
		case "changeNumber":
			if _, ok := fieldSeen[bazelinvocation.FieldChangeNumber]; !ok {
				selectedFields = append(selectedFields, bazelinvocation.FieldChangeNumber)
				fieldSeen[bazelinvocation.FieldChangeNumber] = struct{}{}
			}
		case "patchsetNumber":
			if _, ok := fieldSeen[bazelinvocation.FieldPatchsetNumber]; !ok {
				selectedFields = append(selectedFields, bazelinvocation.FieldPatchsetNumber)
				fieldSeen[bazelinvocation.FieldPatchsetNumber] = struct{}{}
			}
		case "bepCompleted":
			if _, ok := fieldSeen[bazelinvocation.FieldBepCompleted]; !ok {
				selectedFields = append(selectedFields, bazelinvocation.FieldBepCompleted)
				fieldSeen[bazelinvocation.FieldBepCompleted] = struct{}{}
			}
		case "stepLabel":
			if _, ok := fieldSeen[bazelinvocation.FieldStepLabel]; !ok {
				selectedFields = append(selectedFields, bazelinvocation.FieldStepLabel)
				fieldSeen[bazelinvocation.FieldStepLabel] = struct{}{}
			}
		case "id":
		case "__typename":
		default:
			unknownSeen = true
		}
	}
	if !unknownSeen {
		bi.Select(selectedFields...)
	}
	return nil
}

type bazelinvocationPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []BazelInvocationPaginateOption
}

func newBazelInvocationPaginateArgs(rv map[string]any) *bazelinvocationPaginateArgs {
	args := &bazelinvocationPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[whereField].(*BazelInvocationWhereInput); ok {
		args.opts = append(args.opts, WithBazelInvocationFilter(v.Filter))
	}
	return args
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (bip *BazelInvocationProblemQuery) CollectFields(ctx context.Context, satisfies ...string) (*BazelInvocationProblemQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return bip, nil
	}
	if err := bip.collectField(ctx, false, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return bip, nil
}

func (bip *BazelInvocationProblemQuery) collectField(ctx context.Context, oneNode bool, opCtx *graphql.OperationContext, collected graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	var (
		unknownSeen    bool
		fieldSeen      = make(map[string]struct{}, len(bazelinvocationproblem.Columns))
		selectedFields = []string{bazelinvocationproblem.FieldID}
	)
	for _, field := range graphql.CollectFields(opCtx, collected.Selections, satisfies) {
		switch field.Name {

		case "bazelInvocation":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = (&BazelInvocationClient{config: bip.config}).Query()
			)
			if err := query.collectField(ctx, oneNode, opCtx, field, path, mayAddCondition(satisfies, bazelinvocationImplementors)...); err != nil {
				return err
			}
			bip.withBazelInvocation = query
		case "problemType":
			if _, ok := fieldSeen[bazelinvocationproblem.FieldProblemType]; !ok {
				selectedFields = append(selectedFields, bazelinvocationproblem.FieldProblemType)
				fieldSeen[bazelinvocationproblem.FieldProblemType] = struct{}{}
			}
		case "label":
			if _, ok := fieldSeen[bazelinvocationproblem.FieldLabel]; !ok {
				selectedFields = append(selectedFields, bazelinvocationproblem.FieldLabel)
				fieldSeen[bazelinvocationproblem.FieldLabel] = struct{}{}
			}
		case "id":
		case "__typename":
		default:
			unknownSeen = true
		}
	}
	if !unknownSeen {
		bip.Select(selectedFields...)
	}
	return nil
}

type bazelinvocationproblemPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []BazelInvocationProblemPaginateOption
}

func newBazelInvocationProblemPaginateArgs(rv map[string]any) *bazelinvocationproblemPaginateArgs {
	args := &bazelinvocationproblemPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[whereField].(*BazelInvocationProblemWhereInput); ok {
		args.opts = append(args.opts, WithBazelInvocationProblemFilter(v.Filter))
	}
	return args
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (b *BlobQuery) CollectFields(ctx context.Context, satisfies ...string) (*BlobQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return b, nil
	}
	if err := b.collectField(ctx, false, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *BlobQuery) collectField(ctx context.Context, oneNode bool, opCtx *graphql.OperationContext, collected graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	var (
		unknownSeen    bool
		fieldSeen      = make(map[string]struct{}, len(blob.Columns))
		selectedFields = []string{blob.FieldID}
	)
	for _, field := range graphql.CollectFields(opCtx, collected.Selections, satisfies) {
		switch field.Name {
		case "uri":
			if _, ok := fieldSeen[blob.FieldURI]; !ok {
				selectedFields = append(selectedFields, blob.FieldURI)
				fieldSeen[blob.FieldURI] = struct{}{}
			}
		case "sizeBytes":
			if _, ok := fieldSeen[blob.FieldSizeBytes]; !ok {
				selectedFields = append(selectedFields, blob.FieldSizeBytes)
				fieldSeen[blob.FieldSizeBytes] = struct{}{}
			}
		case "archivingStatus":
			if _, ok := fieldSeen[blob.FieldArchivingStatus]; !ok {
				selectedFields = append(selectedFields, blob.FieldArchivingStatus)
				fieldSeen[blob.FieldArchivingStatus] = struct{}{}
			}
		case "reason":
			if _, ok := fieldSeen[blob.FieldReason]; !ok {
				selectedFields = append(selectedFields, blob.FieldReason)
				fieldSeen[blob.FieldReason] = struct{}{}
			}
		case "archiveURL":
			if _, ok := fieldSeen[blob.FieldArchiveURL]; !ok {
				selectedFields = append(selectedFields, blob.FieldArchiveURL)
				fieldSeen[blob.FieldArchiveURL] = struct{}{}
			}
		case "id":
		case "__typename":
		default:
			unknownSeen = true
		}
	}
	if !unknownSeen {
		b.Select(selectedFields...)
	}
	return nil
}

type blobPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []BlobPaginateOption
}

func newBlobPaginateArgs(rv map[string]any) *blobPaginateArgs {
	args := &blobPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[whereField].(*BlobWhereInput); ok {
		args.opts = append(args.opts, WithBlobFilter(v.Filter))
	}
	return args
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (b *BuildQuery) CollectFields(ctx context.Context, satisfies ...string) (*BuildQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return b, nil
	}
	if err := b.collectField(ctx, false, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *BuildQuery) collectField(ctx context.Context, oneNode bool, opCtx *graphql.OperationContext, collected graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	var (
		unknownSeen    bool
		fieldSeen      = make(map[string]struct{}, len(build.Columns))
		selectedFields = []string{build.FieldID}
	)
	for _, field := range graphql.CollectFields(opCtx, collected.Selections, satisfies) {
		switch field.Name {

		case "invocations":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = (&BazelInvocationClient{config: b.config}).Query()
			)
			if err := query.collectField(ctx, false, opCtx, field, path, mayAddCondition(satisfies, bazelinvocationImplementors)...); err != nil {
				return err
			}
			b.WithNamedInvocations(alias, func(wq *BazelInvocationQuery) {
				*wq = *query
			})
		case "buildURL":
			if _, ok := fieldSeen[build.FieldBuildURL]; !ok {
				selectedFields = append(selectedFields, build.FieldBuildURL)
				fieldSeen[build.FieldBuildURL] = struct{}{}
			}
		case "buildUUID":
			if _, ok := fieldSeen[build.FieldBuildUUID]; !ok {
				selectedFields = append(selectedFields, build.FieldBuildUUID)
				fieldSeen[build.FieldBuildUUID] = struct{}{}
			}
		case "id":
		case "__typename":
		default:
			unknownSeen = true
		}
	}
	if !unknownSeen {
		b.Select(selectedFields...)
	}
	return nil
}

type buildPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []BuildPaginateOption
}

func newBuildPaginateArgs(rv map[string]any) *buildPaginateArgs {
	args := &buildPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[whereField].(*BuildWhereInput); ok {
		args.opts = append(args.opts, WithBuildFilter(v.Filter))
	}
	return args
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (ef *EventFileQuery) CollectFields(ctx context.Context, satisfies ...string) (*EventFileQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return ef, nil
	}
	if err := ef.collectField(ctx, false, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return ef, nil
}

func (ef *EventFileQuery) collectField(ctx context.Context, oneNode bool, opCtx *graphql.OperationContext, collected graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	var (
		unknownSeen    bool
		fieldSeen      = make(map[string]struct{}, len(eventfile.Columns))
		selectedFields = []string{eventfile.FieldID}
	)
	for _, field := range graphql.CollectFields(opCtx, collected.Selections, satisfies) {
		switch field.Name {

		case "bazelInvocation":
			var (
				alias = field.Alias
				path  = append(path, alias)
				query = (&BazelInvocationClient{config: ef.config}).Query()
			)
			if err := query.collectField(ctx, oneNode, opCtx, field, path, mayAddCondition(satisfies, bazelinvocationImplementors)...); err != nil {
				return err
			}
			ef.withBazelInvocation = query
		case "url":
			if _, ok := fieldSeen[eventfile.FieldURL]; !ok {
				selectedFields = append(selectedFields, eventfile.FieldURL)
				fieldSeen[eventfile.FieldURL] = struct{}{}
			}
		case "modTime":
			if _, ok := fieldSeen[eventfile.FieldModTime]; !ok {
				selectedFields = append(selectedFields, eventfile.FieldModTime)
				fieldSeen[eventfile.FieldModTime] = struct{}{}
			}
		case "protocol":
			if _, ok := fieldSeen[eventfile.FieldProtocol]; !ok {
				selectedFields = append(selectedFields, eventfile.FieldProtocol)
				fieldSeen[eventfile.FieldProtocol] = struct{}{}
			}
		case "mimeType":
			if _, ok := fieldSeen[eventfile.FieldMimeType]; !ok {
				selectedFields = append(selectedFields, eventfile.FieldMimeType)
				fieldSeen[eventfile.FieldMimeType] = struct{}{}
			}
		case "status":
			if _, ok := fieldSeen[eventfile.FieldStatus]; !ok {
				selectedFields = append(selectedFields, eventfile.FieldStatus)
				fieldSeen[eventfile.FieldStatus] = struct{}{}
			}
		case "reason":
			if _, ok := fieldSeen[eventfile.FieldReason]; !ok {
				selectedFields = append(selectedFields, eventfile.FieldReason)
				fieldSeen[eventfile.FieldReason] = struct{}{}
			}
		case "id":
		case "__typename":
		default:
			unknownSeen = true
		}
	}
	if !unknownSeen {
		ef.Select(selectedFields...)
	}
	return nil
}

type eventfilePaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []EventFilePaginateOption
}

func newEventFilePaginateArgs(rv map[string]any) *eventfilePaginateArgs {
	args := &eventfilePaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[whereField].(*EventFileWhereInput); ok {
		args.opts = append(args.opts, WithEventFileFilter(v.Filter))
	}
	return args
}

const (
	afterField     = "after"
	firstField     = "first"
	beforeField    = "before"
	lastField      = "last"
	orderByField   = "orderBy"
	directionField = "direction"
	fieldField     = "field"
	whereField     = "where"
)

func fieldArgs(ctx context.Context, whereInput any, path ...string) map[string]any {
	field := collectedField(ctx, path...)
	if field == nil || field.Arguments == nil {
		return nil
	}
	oc := graphql.GetOperationContext(ctx)
	args := field.ArgumentMap(oc.Variables)
	return unmarshalArgs(ctx, whereInput, args)
}

// unmarshalArgs allows extracting the field arguments from their raw representation.
func unmarshalArgs(ctx context.Context, whereInput any, args map[string]any) map[string]any {
	for _, k := range []string{firstField, lastField} {
		v, ok := args[k]
		if !ok {
			continue
		}
		i, err := graphql.UnmarshalInt(v)
		if err == nil {
			args[k] = &i
		}
	}
	for _, k := range []string{beforeField, afterField} {
		v, ok := args[k]
		if !ok {
			continue
		}
		c := &Cursor{}
		if c.UnmarshalGQL(v) == nil {
			args[k] = c
		}
	}
	if v, ok := args[whereField]; ok && whereInput != nil {
		if err := graphql.UnmarshalInputFromContext(ctx, v, whereInput); err == nil {
			args[whereField] = whereInput
		}
	}

	return args
}

// mayAddCondition appends another type condition to the satisfies list
// if it does not exist in the list.
func mayAddCondition(satisfies []string, typeCond []string) []string {
Cond:
	for _, c := range typeCond {
		for _, s := range satisfies {
			if c == s {
				continue Cond
			}
		}
		satisfies = append(satisfies, c)
	}
	return satisfies
}
