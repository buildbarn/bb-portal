package authschema

import (
	"context"
	"log/slog"

	"entgo.io/ent/entql"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/invocationtarget"
	"github.com/buildbarn/bb-portal/ent/gen/ent/privacy"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testsummary"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
)

func addDefaultFilter(f privacy.Filter, authorizedInstanceNames []any) entql.P {
	return entql.HasEdgeWith(
		"instance_name",
		entql.FieldIn("name", authorizedInstanceNames...),
	)
}

func addAuthenticatedUserFilter(f privacy.Filter, authorizedInstanceNames []any) entql.P {
	return entql.HasEdgeWith(
		"bazel_invocations",
		addDefaultFilter(f, authorizedInstanceNames),
	)
}

func addTestSummaryFilter(f privacy.Filter, authorizedInstanceNames []any) entql.P {
	return entql.HasEdgeWith(
		testsummary.InvocationTargetColumn,
		entql.HasEdgeWith(
			invocationtarget.TargetColumn,
			addDefaultFilter(f, authorizedInstanceNames),
		),
	)
}

func privacyFilterFunc(ctx context.Context, f privacy.Filter, filterFunc func(privacy.Filter, []any) entql.P) error {
	if dbauthservice.BypassDbAuthServiceFromContext(ctx) {
		return privacy.Skip
	}

	dbAuthService := dbauthservice.FromContext(ctx)
	if dbAuthService == nil {
		slog.WarnContext(ctx, "No DbAuthService or DbAuthServiceBypass present in context; denying all access. This is most likely a bug.")
		return privacy.Deny
	}
	authorizedInstanceNames := dbAuthService.GetAuthorizedInstanceNames(ctx)
	// If no instance names are authorized, add a nil value to the list to
	// avoid a syntax error in the generated SQL, e.g., "IN ()".
	authorizedInstanceNames = append(authorizedInstanceNames, nil)

	f.Where(
		filterFunc(f, authorizedInstanceNames),
	)
	return privacy.Skip
}

// Policy for AuthenticatedUser.
func (AuthenticatedUser) Policy() ent.Policy {
	return privacy.FilterFunc(func(ctx context.Context, f privacy.Filter) error {
		return privacyFilterFunc(ctx, f, addAuthenticatedUserFilter)
	})
}

// Policy for BazelInvocation.
func (BazelInvocation) Policy() ent.Policy {
	return privacy.FilterFunc(func(ctx context.Context, f privacy.Filter) error {
		return privacyFilterFunc(ctx, f, addDefaultFilter)
	})
}

// Policy for Build.
func (Build) Policy() ent.Policy {
	return privacy.FilterFunc(func(ctx context.Context, f privacy.Filter) error {
		return privacyFilterFunc(ctx, f, addDefaultFilter)
	})
}

// Policy for Target.
func (Target) Policy() ent.Policy {
	return privacy.FilterFunc(func(ctx context.Context, f privacy.Filter) error {
		return privacyFilterFunc(ctx, f, addDefaultFilter)
	})
}

// Policy for TestSummary.
func (TestSummary) Policy() ent.Policy {
	return privacy.FilterFunc(func(ctx context.Context, f privacy.Filter) error {
		return privacyFilterFunc(ctx, f, addTestSummaryFilter)
	})
}
