#!/bin/sh -l

GITHUB_TOKEN=$1
BRANCH_NAME=$2
WORKSPACE=/github/workspace
CONFIGMAPS_TARGET_PATH=$3
CONFIGMAPS_TARGET=${WORKSPACE}/${CONFIGMAPS_TARGET_PATH}

if [ -z "${CONFIGMAPS_TARGET_PATH}" ]; then
   export CONFIGMAPS_TARGET=${WORKSPACE}/_private/configmaps
fi

mkdir -p $CONFIGMAPS_TARGET

# permission setup
PERMISSION_DIR=${WORKSPACE}/configs/permissions
PERMISSION_CONFIGMAP_FILE=${CONFIGMAPS_TARGET}/model-access-permissions.configmap.yml

# role setup
ROLE_DIR=${WORKSPACE}/configs/roles
ROLE_CONFIGMAP_FILE=${CONFIGMAPS_TARGET}/rbac-config.yml

if [ -z "${GITHUB_TOKEN}" ]; then
    echo "error: no GITHUB_TOKEN supplied"
    exit 1
fi

if [ -z "${BRANCH_NAME}" ]; then
   export BRANCH_NAME=main
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
remote_repo="https://${GITHUB_ACTOR}:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"
git config http.sslVerify false
git config user.name "[GitHub] - Automated Action"
git config user.email "actions@github.com"

git checkout ${BRANCH_NAME}
git add .
timestamp=$(date -u)
git commit -m "[GitHub] - Automated ConfigMap Generation: ${timestamp} - ${GITHUB_SHA}" || exit 0
git pull --rebase origin ${BRANCH_NAME}
git push origin ${BRANCH_NAME}