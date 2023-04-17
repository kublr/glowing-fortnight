#!/bin/bash

# make sure that the script fails on any failure
set -e
set -o pipefail

# verify that git is available
git version >/dev/null 2>&1 || { echo "ERROR: workspace_status.sh: git is not available" >&2 ; false ; }

# check current tree state
echo BUILD_GIT_TREE_STATE "$(
    if git_status="$(git status --porcelain 2>/dev/null)" ; then
        if [ -z "${git_status}" ] ; then
            echo clean
        else
            echo dirty
        fi
    else
        echo na
    fi
)"

echo STABLE_BUILD_GIT_COMMIT "$(git rev-parse HEAD 2>/dev/null || true)"

# this is to make sure that stable and volatile files are always exactly the same when stamping data is the same
echo BUILD_TIMESTAMP 1
echo BUILD_EMBED_LABEL
echo BUILD_HOST na
echo BUILD_USER na
