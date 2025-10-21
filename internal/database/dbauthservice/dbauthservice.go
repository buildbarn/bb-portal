package dbauthservice

import (
	"context"
	"log/slog"
	"time"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/buildbarn/bb-storage/pkg/digest"
)

// DbAuthService is a service that provides a list of instance names that the
// user is authorized to access, based on the instance names stored in the
// database. It caches the list of instance names for a short period of time to
// avoid excessive database queries.
type DbAuthService struct {
	db                     *ent.Client
	clock                  clock.Clock
	instanceNames          []digest.InstanceName
	lastUpdated            time.Time
	instanceNameAuthorizer auth.Authorizer
	updateThreshold        time.Duration
}

// NewDbAuthService creates a new DbAuthService
func NewDbAuthService(db *ent.Client, clock clock.Clock, instanceNameAuthroizer auth.Authorizer, updateThreshold time.Duration) *DbAuthService {
	return &DbAuthService{
		db:                     db,
		clock:                  clock,
		instanceNames:          []digest.InstanceName{},
		lastUpdated:            time.Time{},
		instanceNameAuthorizer: instanceNameAuthroizer,
		updateThreshold:        updateThreshold,
	}
}

// GetInstanceNames returns the list of all instance names stored in the
// database. The list is cached for a short period of time to avoid excessive
// database queries.
func (s *DbAuthService) GetInstanceNames(ctx context.Context) []digest.InstanceName {
	now := s.clock.Now()
	if now.Sub(s.lastUpdated) < s.updateThreshold {
		return s.instanceNames
	}

	dbInstanceNames, err := s.db.InstanceName.Query().All(ctx)
	if err != nil {
		slog.Warn("Failed to update instance name cache", "error", err)
		return s.instanceNames
	}

	s.instanceNames = make([]digest.InstanceName, 0, len(dbInstanceNames))
	for _, dbInstanceName := range dbInstanceNames {
		instanceName, err := digest.NewInstanceName(dbInstanceName.Name)
		if err != nil {
			slog.Warn("Invalid instance name found in database", "instanceName", dbInstanceName.Name, "error", err)
			continue
		}
		s.instanceNames = append(s.instanceNames, instanceName)
	}

	s.lastUpdated = now
	return s.instanceNames
}

// GetAuthorizedInstanceNames returns the list of instance names that the
// user is authorized to access.
func (s *DbAuthService) GetAuthorizedInstanceNames(ctx context.Context) []any {
	instanceNames := s.GetInstanceNames(ctx)
	errors := s.instanceNameAuthorizer.Authorize(ctx, instanceNames)

	authorizedInstanceNames := make([]any, 0, len(instanceNames))
	for i, instanceName := range instanceNames {
		if errors[i] == nil {
			authorizedInstanceNames = append(authorizedInstanceNames, instanceName.String())
		}
	}
	return authorizedInstanceNames
}

type dbAuthServiceKey struct{}

// NewContextWithDbAuthService creates a new Context object
// that has a DbAuthService attached to it.
func NewContextWithDbAuthService(ctx context.Context, dbAuthService *DbAuthService) context.Context {
	return context.WithValue(ctx, dbAuthServiceKey{}, dbAuthService)
}

// FromContext reobtains the DbAuthService that was attached to the
// Context object.
func FromContext(ctx context.Context) *DbAuthService {
	if value := ctx.Value(dbAuthServiceKey{}); value != nil {
		return value.(*DbAuthService)
	}
	return nil
}
