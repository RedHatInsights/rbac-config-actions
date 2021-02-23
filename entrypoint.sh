#!/bin/sh -l

WORKSPACE=/github/workspace
CONFIGMAPS_TARGET=${WORKSPACE}/${INPUT_TARGET_PATH}

if [ -z "${INPUT_TARGET_PATH}" ]; then
   export CONFIGMAPS_TARGET=${WORKSPACE}/_private/configmaps
fi

mkdir -p $CONFIGMAPS_TARGET

# permission setup
PERMISSION_DIR=${WORKSPACE}/configs/permissions
PERMISSION_CONFIGMAP_FILE=${CONFIGMAPS_TARGET}/model-access-permissions.configmap.yml

# role setup
ROLE_DIR=${WORKSPACE}/configs/roles
ROLE_CONFIGMAP_FILE=${CONFIGMAPS_TARGET}/rbac-config.yml

if [ -z "${INPUT_TOKEN}" ]; then
    echo "error: no INPUT_TOKEN supplied"
    exit 1
fi

if [ -z "${INPUT_BRANCH_NAME}" ]; then
   export INPUT_BRANCH_NAME=main
fi

# create configmaps
kubectl create configmap model-access-permissions --from-file $PERMISSION_DIR --dry-run=client --validate=false -o yaml > $PERMISSION_CONFIGMAP_FILE
kubectl create configmap rbac-config --from-file $ROLE_DIR --dry-run=client --validate=false -o yaml > $ROLE_CONFIGMAP_FILE

# add annotations
for f in $PERMISSION_CONFIGMAP_FILE $ROLE_CONFIGMAP_FILE
do
  echo '  annotations:
      qontract.recycle: "true"' >> $f
done

# push the changes
git config http.sslVerify false
git config user.name "[GitHub] - Automated Action"
git config user.email "actions@github.com"

git add .
timestamp=$(date -u)
git commit -m "[GitHub] - Automated ConfigMap Generation: ${timestamp} - ${GITHUB_SHA}" || exit 0
git pull --rebase
git push origin ${INPUT_BRANCH_NAME}
