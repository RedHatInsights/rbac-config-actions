# Generate V1-Only Permissions Data Action
This will generate a KSIL file representing permissions declared in V1 but not yet migrated to the V2 model.

Usage:
```
on:
  pull_request:
    branches:
      - master
name: Generate V1-Only Permissions Data
jobs:
  generate_v1_permissions:
    name: Generate V1-Only Permissions Data
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4.1.7
      - name: Run Generate V1-Only Permissions Data for stage
        uses: RedHatInsights/rbac-config-actions/generate-v1-only-permissions@main
        with:
          ksl: /ksl/project/files/path/ # required
          rbac_permissions: /json/source/files/path/ # required
```