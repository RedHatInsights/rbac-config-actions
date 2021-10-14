# RBAC Convert Config Action
This will generate configMaps from the JSON config for permissions and roles.

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
        uses: RedHatInsights/rbac-config-actions/convert-config@main
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          branch_name: <BRANCH_NAME> # optional - defaults to `main`
          target_path: /_your/destination/path/ # optional - defaults to `/_private/configmaps/`
```