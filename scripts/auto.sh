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
git push origin release/$TAG

echo "check the latest version, sleeping for 1 min for the go registry to sync up"

sleep 1

go list -m -json github.com/bxcodec/gotcha@latest
go list -m -json github.com/bxcodec/gotcha/examples/basic@latest

go run github.com/bxcodec/gotcha/examples/basic@$tag
