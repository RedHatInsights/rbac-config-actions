#!/bin/sh -l

set -e

# generate schemas
# Clone ksl-schema-language repo
git clone --depth=1 https://github.com/project-kessel/ksl-schema-language.git

# Build KSL compiler binary
cd ksl-schema-language
mkdir -p bin/
go build -gcflags "all=-N -l" -o ./bin/ ./...
cd ..

# Run KSL compiler
for ENV in stage prod
do # ex. configs/stage/schemas
  ./ksl-schema-language/bin/ksl -o configs/${ENV}/schemas/schema.zed configs/${ENV}/schemas/src/*.ksl configs/${ENV}/schemas/src/*.json
  rm configs/${ENV}/schemas/src/rbac_v1_permissions.json
done

# Remove KSL repo
rm -rf ksl-schema-language/
