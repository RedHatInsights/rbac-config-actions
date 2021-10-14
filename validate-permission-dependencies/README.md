# RBAC Validate Permission Dependencies Action
This will validate permission verb dependencies are valid.

Usage:
```
on:
  pull_request:
    branches:
      - mani
name: PR Workflow
jobs:
  validate_configurations:
    name: Validate JSON Config
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Validate Permissions' Dependencies
        uses: RedHatInsights/rbac-config-actions/validate-permission-dependencies@main
        with:
          permissions_path_pattern: 'configs/**/*/permissions/*.json'
```