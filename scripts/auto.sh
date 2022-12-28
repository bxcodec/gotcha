#!/bin/bash

echo "updating repo"

git pull

go work sync

git describe --tags

echo "what's the tag"

read tag

tag=$(echo "$tag" | sed 's/ //g')

export TAG=$tag

chmod +x scripts/release.sh
chmod +x scripts/tag.sh

./scripts/release.sh
./scripts/tag.sh

git add .
git commit -m "chore: update deps-$tag"

git push origin --tags
git push

echo "verify the latest versionx"

go list -m -json github.com/bxcodec/gotcha@latest
go list -m -json github.com/bxcodec/gotcha/examples/basic@latest
