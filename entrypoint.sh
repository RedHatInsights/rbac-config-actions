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
   export INPUT_BRANCH_NAME=main
fi

git config http.sslVerify false
git config user.name "[GitHub] - Automated Action"
git config user.email "actions@github.com"

indent() { sed '2,$s/^/  /'; }

for ENV in ci qa stage prod
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
    echo '  annotations:
    qontract.recycle: "true"' >> $f
  done
done

git fetch
git checkout ${INPUT_BRANCH_NAME}
git pull origin ${INPUT_BRANCH_NAME} --rebase

# push the changes
git add .
timestamp=$(date -u)
git commit -m "[GitHub] - Automated ConfigMap Generation: ${timestamp} - ${GITHUB_SHA}" || exit 0
git push origin ${INPUT_BRANCH_NAME}
