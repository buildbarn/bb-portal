# CI integration

These are example integrations against different CI environments.

## Github Actions

Backend extractor: Use `config/gh-actions.jmespath`, 

Frontend config:
```
additionalBuildColumns: [
  { title: 'Repo', value_key: 'repo', url_key: 'repo_url' },
  { title: 'PR', value_key: 'pull_request', url_key: 'pull_request_url' },
  { title: 'Workflow', value_key: 'workflow', url_key: 'workflow_url' },
],
additionalBuildInvocationColumns: [
  { title: 'Job', value_key: 'job' },
  { title: 'Action', value_key: 'action' },
],
```

## Gitlab CI

Backend extractor: Use `config/gitlab.jmespath`, 

Frontend config:
```
additionalBuildColumns: [
  { title: 'Repo', value_key: 'repo', url_key: 'repo_url' },
  { title: 'Pipeline', value_key: 'pipeline', url_key: 'pipeline_url' },
],
additionalBuildInvocationColumns: [
  { title: 'Job', value_key: 'job', url_key: 'job_url' },
  { title: 'Job stage', value_key: 'job_stage' },
],
```

## SemaphoreCI

Backend extractor: Use `config/semaphore.jmespath`

Frontend config:
```
additionalBuildColumns: [
  { title: 'Repo', value_key: 'repo', url_key: 'repo_url' },
  { title: 'Workflow', value_key: 'semaphore_workflow', url_key: 'semaphore_workflow' },
],
additionalBuildInvocationColumns: [
  { title: 'Pipeline', value_key: 'semaphore_pipeline', url_key: 'semaphore_pipeline_url' },
  { title: 'Job', value_key: 'semaphore_job', url_key: 'semaphore_job_url' },
],
```
