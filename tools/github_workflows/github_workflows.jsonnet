local workflows_template = import 'tools/github_workflows/workflows_template.libsonnet';

workflows_template.getWorkflows(
  [
    'bb_portal',
  ],
  [
    'bb_portal:bb_portal',
  ],
)
