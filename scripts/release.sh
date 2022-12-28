#! /bin/bash

set -e

help() {
    cat <<- EOF
Usage: TAG=tag $0

Updates version in go.mod files and pushes a new branch to GitHub.

VARIABLES:
  TAG        git tag, for example, v1.0.0
EOF
    exit 0
}

if [ -z "$TAG" ]
then
    printf "TAG is required\n\n"
    help
fi

TAG_REGEX="^v(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)(\\-[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?(\\+[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?$"
if ! [[ "${TAG}" =~ ${TAG_REGEX} ]]; then
    printf "TAG is not valid: ${TAG}\n\n"
    exit 1
fi

TAG_FOUND=`git tag --list ${TAG}`
if [[ ${TAG_FOUND} = ${TAG} ]] ; then
    printf "tag ${TAG} already exists\n\n"
#    exit 1
fi

if ! git diff --quiet
then
    printf "working tree is not clean\n\n"
    git status
    exit 1
fi

PACKAGE_DIRS=$(find . -mindepth 2 -type f -name 'go.mod' -exec dirname {} \; \
  | sed 's/^\.\///' \
  | sort)

for dir in $PACKAGE_DIRS
do
    sed --in-place \
      "s/bxcodec\/gotcha\([^ ]*\) v.*/bxcodec\/gotcha\1 ${TAG}/" "${dir}/go.mod"
done

pwd

sed --in-place "s/\(return \)\"[^\"]*\"/\1\"${TAG#v}\"/" version.go

git checkout -b release/${TAG}

git add -u
git commit -m "chore: release $TAG"
git tag ${TAG}
git push origin ${TAG}

for dir in $PACKAGE_DIRS
do
    printf "${dir}: go get -u && go mod tidy -compat=1.19\n"
    go get github.com/bxcodec/gotcha@${TAG}
    (cd ./${dir} && go get -u && go mod tidy) # -compat=1.19
done

git push -u origin release/${TAG}
