package summary

// IsBuildEnvKey checks if an environment variable key is a build-level
// environment variable - one that is expected to be the same for multiple
// invocations that occur as part of the same build.
//
// In this context, "Build" refers to a CI build, like a Jenkins
// build, not a 'bazel build'.
func IsBuildEnvKey(k string) bool {
	switch k {
	case // Jenkins variables, from
		// https://www.jenkins.io/doc/book/pipeline/jenkinsfile/#using-environment-variables.
		"BUILD_ID", "BUILD_NUMBER", "BUILD_TAG", "BUILD_URL", "EXECUTOR_NUMBER", "JOB_NAME":
		return true

	case //	Gerrit plugin variables. NOTE: GERRIT_CHANGE_COMMIT_MESSAGE is excluded because it may be too large.
		"GERRIT_BRANCH", "GERRIT_CHANGE_COMMIT_MESSAGE", "GERRIT_CHANGE_ID", "GERRIT_CHANGE_NUMBER",
		"GERRIT_CHANGE_SUBJECT", "GERRIT_PATCHSET_NUMBER":
		return true

	case // Jenkins git / github plugin
		"GIT_SHA":
		return true
	case // Git variables
		// NOTE: Left out ones already in Jenkins, like BUILD_URL.
		// NOTE: Left off GIT_COMMIT_MESSAGE because it may be too long for the JSONB index.
		"GIT_AUTHOR_EMAILS", "GIT_AUTHOR_EMAIL", "GIT_AUTHOR_NAMES", "GIT_AUTHOR_NAME", "GIT_BRANCH",
		"GIT_COMMITTER_EMAILS", "GIT_COMMITTER_EMAIL", "GIT_COMMITTER_NAMES", "GIT_COMMITTER_NAME", "GIT_COMMIT_SHORT", "GIT_COMMIT", "GIT_HASH", "GIT_PREVIOUS_COMMIT", "GIT_PREVIOUS_SUCCESSFUL_COMMIT", "GIT_URL", "GIT_TAG", "PIPELINE_SPEC_ID", "PROJECT_ID", "PROJECT":
		return true

	case //	Git PR variables
		// Excludes ones already in Git variables above.
		// Excludes GIT_PR_COMMIT_MESSAGE becaues it may be too long.
		"GIT_PR_COMMIT_AUTHOR", "GIT_PR_COMMIT_COMMITTER", "GIT_PR_COMMIT", "GIT_PR_ID", "GIT_PR_SOURCE_BRANCH",
		"GIT_PR_TARGET_BRANCH", "GIT_PR_TARGET_COMMIT", "GIT_PR_URL":
		return true
	default:
		return false
	}
}
