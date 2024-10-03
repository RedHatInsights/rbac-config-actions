# Generate Schema Validation Action
This will validate generated schema from the underlying ksl files.

Usage:
```
on:
  pull_request:
    branches:
      - main
name: Validate generated schema
jobs:
  validate_schema:
    name: validate schema
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Generate & validate schema
        uses: RedHatInsights/rbac-config-actions/validate-schema@main
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
```