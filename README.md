# RBAC Post-commit Action
A GitHub action to run post-commit actions for RBAC config. This will generate
ConfigMaps from the JSON config for permissions and roles.

Usage:
```
on:
  push:
    branches:
      - main
name: RBAC Config to ConfigMap
jobs:
  convert_config:
    name: JSON to ConfigMap
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Converting JSON config to ConfigMaps
        uses: coderbydesign/rbac-config-post-commit-action@main
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          branch_name: <BRANCH_NAME> # optional - defaults to `main`
```