package invocationmetadataextraction

import (
	"github.com/buildbarn/bb-storage/pkg/jmespath"
)

// SourceControl contains source control data for a invocation
type SourceControl struct {
	Repo      *string
	RepoURL   *string
	Ref       *string
	RefURL    *string
	Commit    *string
	CommitURL *string
}

// InvocationMetadata contains metadata for a invocation
type InvocationMetadata struct {
	Username       *string
	Hostname       *string
	SourceControls []SourceControl
	InvocationTags map[string]string
	BuildTags      map[string]string
}

// NewInvocationMetadata creates a new empty InvocationMetadata
func NewInvocationMetadata() *InvocationMetadata {
	return &InvocationMetadata{
		Username:       nil,
		Hostname:       nil,
		SourceControls: []SourceControl{},
		InvocationTags: map[string]string{},
		BuildTags:      map[string]string{},
	}
}

func parseSourceControlMap(sc map[string]any) *SourceControl {
	if sc == nil {
		return nil
	}

	shouldSave := false

	getField := func(key string) *string {
		if value, ok := sc[key].(string); ok && value != "" {
			shouldSave = true
			return &value
		}
		return nil
	}

	scStruct := SourceControl{
		Repo:      getField("repo"),
		RepoURL:   getField("repoUrl"),
		Ref:       getField("ref"),
		RefURL:    getField("refUrl"),
		Commit:    getField("commit"),
		CommitURL: getField("commitUrl"),
	}

	if !shouldSave {
		return nil
	}

	return &scStruct
}

// ExtractInvocationMetadata extracts invocation metadata from the environment
// variables using the provided JMESPath expression.
func ExtractInvocationMetadata(extractor *jmespath.Expression, envVars map[string]string) *InvocationMetadata {
	if extractor == nil {
		return nil
	}

	// Convert map[string]string to map[string]any that the extractor needs
	searchVars := make(map[string]any, len(envVars))
	for k, v := range envVars {
		searchVars[k] = v
	}

	searchResults, err := extractor.Search(map[string]any{
		"env": searchVars,
	})
	if err != nil || searchResults == nil {
		return nil
	}

	resMap, ok := searchResults.(map[string]any)
	if !ok {
		return nil
	}

	metadata := NewInvocationMetadata()

	if username, ok := resMap["username"].(string); ok && username != "" {
		metadata.Username = &username
	}
	if hostname, ok := resMap["hostname"].(string); ok && hostname != "" {
		metadata.Hostname = &hostname
	}
	if sourceControls, ok := resMap["sourceControls"].([]any); ok {
		for _, sourceControl := range sourceControls {
			if sc, ok := sourceControl.(map[string]any); ok {
				if scres := parseSourceControlMap(sc); scres != nil {
					metadata.SourceControls = append(metadata.SourceControls, *scres)
				}
			}
		}
	}
	if invocationTags, ok := resMap["invocationTags"].(map[string]any); ok {
		for key, anyValue := range invocationTags {
			if value, ok := anyValue.(string); ok && value != "" {
				metadata.InvocationTags[key] = value
			}
		}
	}
	if buildTags, ok := resMap["buildTags"].(map[string]any); ok {
		for key, anyValue := range buildTags {
			if value, ok := anyValue.(string); ok && value != "" {
				metadata.BuildTags[key] = value
			}
		}
	}

	return metadata
}
