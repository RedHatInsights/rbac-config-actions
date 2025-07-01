# Generate V1-Only Permissions Data

This utility generates a KSIL file representing permissions declared in V1 but not yet migrated to the V2 model.

## Installation

### As a Go Tool

You can install this utility directly using `go install`:

```bash
# Install from the repository
go install github.com/RedHatInsights/rbac-config-actions/generate-v1-only-permissions/cmd/generate-v1-only-permissions@latest

# Or install locally after cloning the repository
cd generate-v1-only-permissions
go install ./cmd/generate-v1-only-permissions
```

### Usage as Command Line Tool

After installation, you can use the tool directly:

```bash
generate-v1-only-permissions -ksl /path/to/ksl/project -rbac-permissions-json /path/to/rbac/permissions
```

#### Command Line Options

- `-ksl`: The path to the ksl project directory (where the migrated_apps.lst file is) - **required**
- `-rbac-permissions-json`: The path to the directory containing RBAC permissions .json files for the current environment - **required**

#### Example

```bash
generate-v1-only-permissions \
  -ksl /home/user/ksl-project \
  -rbac-permissions-json /home/user/rbac-permissions
```

## Usage as GitHub Action

You can also use this tool as a GitHub Action in your workflows:

```yaml
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