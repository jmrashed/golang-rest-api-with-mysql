#!/bin/bash

set -e

VERSION=${1:-""}

if [ -z "$VERSION" ]; then
    echo "Error: Version number is required"
    echo "Usage: $0 <version>"
    exit 1
fi

echo "Preparing release v$VERSION"

if ! [[ $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Invalid version format. Use semantic versioning (e.g., 1.1.0)"
    exit 1
fi

if ! git diff-index --quiet HEAD --; then
    echo "Error: You have uncommitted changes. Please commit or stash them first."
    exit 1
fi

sed -i "s/VERSION=.*/VERSION=$VERSION/" Makefile

TODAY=$(date +%Y-%m-%d)
sed -i "s/## \[Unreleased\]/## [Unreleased]\n\n## [$VERSION] - $TODAY/" CHANGELOG.md

make test
make build

echo "Release v$VERSION prepared successfully!"
echo "Next steps:"
echo "1. git add . && git commit -m 'Prepare release v$VERSION'"
echo "2. git tag v$VERSION && git push origin v$VERSION"