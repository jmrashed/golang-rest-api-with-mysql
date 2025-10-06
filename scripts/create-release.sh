#!/bin/bash

set -e

VERSION=${1:-""}

if [ -z "$VERSION" ]; then
    echo "Error: Version number is required"
    echo "Usage: $0 <version>"
    exit 1
fi

echo "Creating and pushing tag v$VERSION"

git tag -a "v$VERSION" -m "Release v$VERSION"
git push origin "v$VERSION"

echo "Tag v$VERSION created and pushed successfully!"
echo "GitHub release will be created automatically via workflow."