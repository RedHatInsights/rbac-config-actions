#!/bin/sh -l

WORKSPACE=/github/workspace
CONFIGMAPS_TARGET=${WORKSPACE}/${INPUT_TARGET_PATH}

if [ -z "${INPUT_TARGET_PATH}" ]; then
   export CONFIGMAPS_TARGET=${WORKSPACE}/_private/configmaps
fi

if [ -z "${INPUT_TOKEN}" ]; then
    echo "error: no INPUT_TOKEN supplied"
    exit 1
fi

if [ -z "${INPUT_BRANCH_NAME}" ]; then
   export INPUT_BRANCH_NAME=configmaps-schema
fi

git config --global --add safe.directory ${WORKSPACE}
git config http.sslVerify false
git config user.name "[GitHub] - Automated Action"
git config user.email "actions@github.com"

# Remove configmaps branch from remote if it already exists.
if [ "${INPUT_BRANCH_NAME}" != "main" ] || [  "${INPUT_BRANCH_NAME}" != "master" ]; then
  if git ls-remote --exit-code --heads origin "${INPUT_BRANCH_NAME}" >/dev/null 2>&1; then
    git push origin --delete "${INPUT_BRANCH_NAME}"
  fi
fi

# sync the input branch
git fetch
git checkout -b "${INPUT_BRANCH_NAME}" --no-track origin/master

indent() { sed '2,$s/^/  /'; }

for ENV in stage prod
do
  mkdir -p $CONFIGMAPS_TARGET/${ENV}

  PERMISSION_DIR=${WORKSPACE}/configs/${ENV}/permissions
  PERMISSION_CONFIGMAP_FILE=${CONFIGMAPS_TARGET}/${ENV}/model-access-permissions.configmap.yml

  ROLE_DIR=${WORKSPACE}/configs/${ENV}/roles
  ROLE_CONFIGMAP_FILE=${CONFIGMAPS_TARGET}/${ENV}/rbac-config.yml

  # create templates
  for f in $PERMISSION_CONFIGMAP_FILE $ROLE_CONFIGMAP_FILE
  do
    echo -n 'kind: Template
apiVersion: v1
objects:
- ' > $f
  done

  # create configmaps
  kubectl create configmap model-access-permissions --from-file $PERMISSION_DIR --dry-run=client --validate=false -o yaml | indent >> $PERMISSION_CONFIGMAP_FILE
  kubectl create configmap rbac-config --from-file $ROLE_DIR --dry-run=client --validate=false -o yaml | indent >> $ROLE_CONFIGMAP_FILE

  # add annotations
  for f in $PERMISSION_CONFIGMAP_FILE $ROLE_CONFIGMAP_FILE
  do
    echo '    annotations:
      qontract.recycle: "true"' >> $f
  done
done

# generate schemas
# Clone ksl-schema-language repo
git clone --depth=1 https://github.com/project-kessel/ksl-schema-language.git

# Build KSL compiler binary
cd ksl-schema-language
mkdir -p bin/
go build -gcflags "all=-N -l" -o ./bin/ ./cmd/ksl/...
cd ..

# Run KSL compiler
for ENV in stage prod
do # ex. configs/stage/schemas
  ./ksl-schema-language/bin/ksl -o configs/${ENV}/schemas/schema.zed configs/${ENV}/schemas/src/*.ksl configs/${ENV}/schemas/src/*.json || exit 1
  rm configs/${ENV}/schemas/src/rbac_v1_permissions.json
done

# Remove KSL repo
rm -rf ksl-schema-language/

# push the changes
git add .
git add --force configs/${ENV}/schemas/schema.zed #Override .gitignore
timestamp=$(date -u)
git commit -m "[GitHub] - Automated ConfigMap & Schema Generation: ${timestamp} - ${GITHUB_SHA}" || exit 0
git push origin ${INPUT_BRANCH_NAME}
