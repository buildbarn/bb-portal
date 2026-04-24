package invocationmetadataextraction_test

import (
	"testing"

	"github.com/buildbarn/bb-portal/pkg/invocationmetadataextraction"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/buildbarn/bb-storage/pkg/jmespath"
	jmespath_config "github.com/buildbarn/bb-storage/pkg/proto/configuration/jmespath"
	"github.com/stretchr/testify/require"
)

func ptr(s string) *string {
	return &s
}

func extractorFromExpression(t *testing.T, expression string) *jmespath.Expression {
	jmespathExpression := &jmespath_config.Expression{
		Expression: expression,
	}
	extractor, err := jmespath.NewExpressionFromConfiguration(jmespathExpression, nil, clock.SystemClock)
	require.NoError(t, err)
	return extractor
}

func TestInvocationMetadataExtraction(t *testing.T) {
	fullExtractor := extractorFromExpression(t, `{
		"username": env.USER
		"hostname": env.HOSTNAME
		"sourceControls": [
			{
			"repo": env.REPO_1
			"repoUrl": env.REPO_URL_1
			"ref": env.REF_1
			"refUrl": env.REF_URL_1
			"commit": env.COMMIT_1
			"commitUrl": env.COMMIT_URL_1
			"dummyKey": env.DUMMY_VALUE
			},
			{
			"repo": env.REPO_2
			"repoUrl": env.REPO_URL_2
			"ref": env.REF_2
			"refUrl": env.REF_URL_2
			"commit": env.COMMIT_2
			"commitUrl": env.COMMIT_URL_2
			"dummyKey": env.DUMMY_VALUE
			}
		]
		"invocationTags": {
			"invTag1": env.INV_TAG_1
			"invTag2": env.INV_TAG_2
		}
		"buildTags": {
			"buildTag1": env.BUILD_TAG_1
			"buildTag2": env.BUILD_TAG_2
		}
		"dummyKey": env.DUMMY_VALUE
	}`)

	t.Run("NoExpression", func(t *testing.T) {
		extractor := extractorFromExpression(t, `{
			"foo": "bar"
		}`)
		metadata := invocationmetadataextraction.ExtractInvocationMetadata(extractor, map[string]string{
			"USER":     "user",
			"HOSTNAME": "hostname",
		})
		require.Equal(t, invocationmetadataextraction.NewInvocationMetadata(), metadata)
	})

	t.Run("NoData", func(t *testing.T) {
		extractor := extractorFromExpression(t, `{
			"username": env.USER
			"hostname": env.HOSTNAME
		}`)
		metadata := invocationmetadataextraction.ExtractInvocationMetadata(extractor, map[string]string{})
		require.Equal(t, invocationmetadataextraction.NewInvocationMetadata(), metadata)
	})

	t.Run("MissingUsername", func(t *testing.T) {
		metadata := invocationmetadataextraction.ExtractInvocationMetadata(fullExtractor, map[string]string{
			"HOSTNAME":     "hostname",
			"REPO_1":       "repo1",
			"REPO_URL_1":   "repoUrl1",
			"REF_1":        "ref1",
			"REF_URL_1":    "refUrl1",
			"COMMIT_1":     "commit1",
			"COMMIT_URL_1": "commitUrl1",
			"REPO_2":       "repo2",
			"REPO_URL_2":   "repoUrl2",
			"REF_2":        "ref2",
			"REF_URL_2":    "refUrl2",
			"COMMIT_2":     "commit2",
			"COMMIT_URL_2": "commitUrl2",
			"INV_TAG_1":    "invTag1",
			"INV_TAG_2":    "invTag2",
			"BUILD_TAG_1":  "buildTag1",
			"BUILD_TAG_2":  "buildTag2",
		})

		require.Equal(t, &invocationmetadataextraction.InvocationMetadata{
			Username: nil,
			Hostname: ptr("hostname"),
			SourceControls: []invocationmetadataextraction.SourceControl{
				{
					Repo:      ptr("repo1"),
					RepoURL:   ptr("repoUrl1"),
					Ref:       ptr("ref1"),
					RefURL:    ptr("refUrl1"),
					Commit:    ptr("commit1"),
					CommitURL: ptr("commitUrl1"),
				},
				{
					Repo:      ptr("repo2"),
					RepoURL:   ptr("repoUrl2"),
					Ref:       ptr("ref2"),
					RefURL:    ptr("refUrl2"),
					Commit:    ptr("commit2"),
					CommitURL: ptr("commitUrl2"),
				},
			},
			InvocationTags: map[string]string{
				"invTag1": "invTag1",
				"invTag2": "invTag2",
			},
			BuildTags: map[string]string{
				"buildTag1": "buildTag1",
				"buildTag2": "buildTag2",
			},
		}, metadata)
	})

	t.Run("MissingHostname", func(t *testing.T) {
		metadata := invocationmetadataextraction.ExtractInvocationMetadata(fullExtractor, map[string]string{
			"USER":         "user",
			"REPO_1":       "repo1",
			"REPO_URL_1":   "repoUrl1",
			"REF_1":        "ref1",
			"REF_URL_1":    "refUrl1",
			"COMMIT_1":     "commit1",
			"COMMIT_URL_1": "commitUrl1",
			"REPO_2":       "repo2",
			"REPO_URL_2":   "repoUrl2",
			"REF_2":        "ref2",
			"REF_URL_2":    "refUrl2",
			"COMMIT_2":     "commit2",
			"COMMIT_URL_2": "commitUrl2",
			"INV_TAG_1":    "invTag1",
			"INV_TAG_2":    "invTag2",
			"BUILD_TAG_1":  "buildTag1",
			"BUILD_TAG_2":  "buildTag2",
		})

		require.Equal(t, &invocationmetadataextraction.InvocationMetadata{
			Username: ptr("user"),
			Hostname: nil,
			SourceControls: []invocationmetadataextraction.SourceControl{
				{
					Repo:      ptr("repo1"),
					RepoURL:   ptr("repoUrl1"),
					Ref:       ptr("ref1"),
					RefURL:    ptr("refUrl1"),
					Commit:    ptr("commit1"),
					CommitURL: ptr("commitUrl1"),
				},
				{
					Repo:      ptr("repo2"),
					RepoURL:   ptr("repoUrl2"),
					Ref:       ptr("ref2"),
					RefURL:    ptr("refUrl2"),
					Commit:    ptr("commit2"),
					CommitURL: ptr("commitUrl2"),
				},
			},
			InvocationTags: map[string]string{
				"invTag1": "invTag1",
				"invTag2": "invTag2",
			},
			BuildTags: map[string]string{
				"buildTag1": "buildTag1",
				"buildTag2": "buildTag2",
			},
		}, metadata)
	})

	t.Run("MissingSourceControl", func(t *testing.T) {
		metadata := invocationmetadataextraction.ExtractInvocationMetadata(fullExtractor, map[string]string{
			"USER":        "user",
			"HOSTNAME":    "hostname",
			"INV_TAG_1":   "invTag1",
			"INV_TAG_2":   "invTag2",
			"BUILD_TAG_1": "buildTag1",
			"BUILD_TAG_2": "buildTag2",
			"DUMMY_VALUE": "dummyValue",
		})

		require.Equal(t, &invocationmetadataextraction.InvocationMetadata{
			Username:       ptr("user"),
			Hostname:       ptr("hostname"),
			SourceControls: []invocationmetadataextraction.SourceControl{},
			InvocationTags: map[string]string{
				"invTag1": "invTag1",
				"invTag2": "invTag2",
			},
			BuildTags: map[string]string{
				"buildTag1": "buildTag1",
				"buildTag2": "buildTag2",
			},
		}, metadata)
	})

	t.Run("SingleSourceControl", func(t *testing.T) {
		metadata := invocationmetadataextraction.ExtractInvocationMetadata(fullExtractor, map[string]string{
			"USER":         "user",
			"HOSTNAME":     "hostname",
			"REPO_1":       "repo1",
			"REPO_URL_1":   "repoUrl1",
			"REF_1":        "ref1",
			"REF_URL_1":    "refUrl1",
			"COMMIT_1":     "commit1",
			"COMMIT_URL_1": "commitUrl1",
			"INV_TAG_1":    "invTag1",
			"INV_TAG_2":    "invTag2",
			"BUILD_TAG_1":  "buildTag1",
			"BUILD_TAG_2":  "buildTag2",
			"DUMMY_VALUE":  "dummyValue",
		})

		require.Equal(t, &invocationmetadataextraction.InvocationMetadata{
			Username: ptr("user"),
			Hostname: ptr("hostname"),
			SourceControls: []invocationmetadataextraction.SourceControl{
				{
					Repo:      ptr("repo1"),
					RepoURL:   ptr("repoUrl1"),
					Ref:       ptr("ref1"),
					RefURL:    ptr("refUrl1"),
					Commit:    ptr("commit1"),
					CommitURL: ptr("commitUrl1"),
				},
			},
			InvocationTags: map[string]string{
				"invTag1": "invTag1",
				"invTag2": "invTag2",
			},
			BuildTags: map[string]string{
				"buildTag1": "buildTag1",
				"buildTag2": "buildTag2",
			},
		}, metadata)
	})

	t.Run("MissingInvocationTags", func(t *testing.T) {
		metadata := invocationmetadataextraction.ExtractInvocationMetadata(fullExtractor, map[string]string{
			"USER":         "user",
			"HOSTNAME":     "hostname",
			"REPO_1":       "repo1",
			"REPO_URL_1":   "repoUrl1",
			"REF_1":        "ref1",
			"REF_URL_1":    "refUrl1",
			"COMMIT_1":     "commit1",
			"COMMIT_URL_1": "commitUrl1",
			"REPO_2":       "repo2",
			"REPO_URL_2":   "repoUrl2",
			"REF_2":        "ref2",
			"REF_URL_2":    "refUrl2",
			"COMMIT_2":     "commit2",
			"COMMIT_URL_2": "commitUrl2",
			"BUILD_TAG_1":  "buildTag1",
			"BUILD_TAG_2":  "buildTag2",
		})

		require.Equal(t, &invocationmetadataextraction.InvocationMetadata{
			Username: ptr("user"),
			Hostname: ptr("hostname"),
			SourceControls: []invocationmetadataextraction.SourceControl{
				{
					Repo:      ptr("repo1"),
					RepoURL:   ptr("repoUrl1"),
					Ref:       ptr("ref1"),
					RefURL:    ptr("refUrl1"),
					Commit:    ptr("commit1"),
					CommitURL: ptr("commitUrl1"),
				},
				{
					Repo:      ptr("repo2"),
					RepoURL:   ptr("repoUrl2"),
					Ref:       ptr("ref2"),
					RefURL:    ptr("refUrl2"),
					Commit:    ptr("commit2"),
					CommitURL: ptr("commitUrl2"),
				},
			},
			InvocationTags: map[string]string{},
			BuildTags: map[string]string{
				"buildTag1": "buildTag1",
				"buildTag2": "buildTag2",
			},
		}, metadata)
	})

	t.Run("MissingBuildTags", func(t *testing.T) {
		metadata := invocationmetadataextraction.ExtractInvocationMetadata(fullExtractor, map[string]string{
			"USER":         "user",
			"HOSTNAME":     "hostname",
			"REPO_1":       "repo1",
			"REPO_URL_1":   "repoUrl1",
			"REF_1":        "ref1",
			"REF_URL_1":    "refUrl1",
			"COMMIT_1":     "commit1",
			"COMMIT_URL_1": "commitUrl1",
			"REPO_2":       "repo2",
			"REPO_URL_2":   "repoUrl2",
			"REF_2":        "ref2",
			"REF_URL_2":    "refUrl2",
			"COMMIT_2":     "commit2",
			"COMMIT_URL_2": "commitUrl2",
			"INV_TAG_1":    "invTag1",
			"INV_TAG_2":    "invTag2",
		})

		require.Equal(t, &invocationmetadataextraction.InvocationMetadata{
			Username: ptr("user"),
			Hostname: ptr("hostname"),
			SourceControls: []invocationmetadataextraction.SourceControl{
				{
					Repo:      ptr("repo1"),
					RepoURL:   ptr("repoUrl1"),
					Ref:       ptr("ref1"),
					RefURL:    ptr("refUrl1"),
					Commit:    ptr("commit1"),
					CommitURL: ptr("commitUrl1"),
				},
				{
					Repo:      ptr("repo2"),
					RepoURL:   ptr("repoUrl2"),
					Ref:       ptr("ref2"),
					RefURL:    ptr("refUrl2"),
					Commit:    ptr("commit2"),
					CommitURL: ptr("commitUrl2"),
				},
			},
			InvocationTags: map[string]string{
				"invTag1": "invTag1",
				"invTag2": "invTag2",
			},
			BuildTags: map[string]string{},
		}, metadata)
	})

	t.Run("FullDataWithDummyValues", func(t *testing.T) {
		metadata := invocationmetadataextraction.ExtractInvocationMetadata(fullExtractor, map[string]string{
			"USER":         "user",
			"HOSTNAME":     "hostname",
			"REPO_1":       "repo1",
			"REPO_URL_1":   "repoUrl1",
			"REF_1":        "ref1",
			"REF_URL_1":    "refUrl1",
			"COMMIT_1":     "commit1",
			"COMMIT_URL_1": "commitUrl1",
			"REPO_2":       "repo2",
			"REPO_URL_2":   "repoUrl2",
			"REF_2":        "ref2",
			"REF_URL_2":    "refUrl2",
			"COMMIT_2":     "commit2",
			"COMMIT_URL_2": "commitUrl2",
			"INV_TAG_1":    "invTag1",
			"INV_TAG_2":    "invTag2",
			"BUILD_TAG_1":  "buildTag1",
			"BUILD_TAG_2":  "buildTag2",
			"DUMMY_VALUE":  "dummyValue",
		})

		require.Equal(t, &invocationmetadataextraction.InvocationMetadata{
			Username: ptr("user"),
			Hostname: ptr("hostname"),
			SourceControls: []invocationmetadataextraction.SourceControl{
				{
					Repo:      ptr("repo1"),
					RepoURL:   ptr("repoUrl1"),
					Ref:       ptr("ref1"),
					RefURL:    ptr("refUrl1"),
					Commit:    ptr("commit1"),
					CommitURL: ptr("commitUrl1"),
				},
				{
					Repo:      ptr("repo2"),
					RepoURL:   ptr("repoUrl2"),
					Ref:       ptr("ref2"),
					RefURL:    ptr("refUrl2"),
					Commit:    ptr("commit2"),
					CommitURL: ptr("commitUrl2"),
				},
			},
			InvocationTags: map[string]string{
				"invTag1": "invTag1",
				"invTag2": "invTag2",
			},
			BuildTags: map[string]string{
				"buildTag1": "buildTag1",
				"buildTag2": "buildTag2",
			},
		}, metadata)
	})
}
