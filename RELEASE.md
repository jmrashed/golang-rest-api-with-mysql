# Release Process

This document outlines the release process for the Golang REST API project.

## Prerequisites

- Git repository with main branch
- GitHub repository with Actions enabled
- Proper permissions to create tags and releases

## Release Steps

### 1. Prepare Release

```bash
# Update version and run pre-release checks
make prepare-release VERSION=1.1.0
```

This will:
- Update version in Makefile
- Update CHANGELOG.md with release date
- Run tests and build checks
- Build Docker image

### 2. Commit Changes

```bash
git add .
git commit -m "Prepare release v1.1.0"
git push origin main
```

### 3. Create and Push Tag

```bash
# Create and push release tag
make create-release VERSION=1.1.0
```

This will:
- Create annotated Git tag
- Push tag to origin
- Trigger GitHub Actions release workflow

### 4. GitHub Release

The GitHub Actions workflow will automatically:
- Build binaries for multiple platforms
- Create checksums
- Build Docker image
- Extract changelog notes
- Create GitHub release with assets

## Manual Release (Alternative)

If you prefer manual control:

```bash
# 1. Prepare release
bash scripts/prepare-release.sh 1.1.0

# 2. Commit changes
git add .
git commit -m "Prepare release v1.1.0"

# 3. Create tag
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0

# 4. GitHub release will be created automatically
```

## Release Assets

Each release includes:
- Linux AMD64 binary
- Linux ARM64 binary  
- Windows AMD64 binary
- macOS AMD64 binary
- macOS ARM64 binary
- Docker image (tar.gz)
- Checksums file

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- MAJOR.MINOR.PATCH (e.g., 1.2.3)
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes (backward compatible)

## Rollback

If issues are found after release:

```bash
# Delete tag locally and remotely
git tag -d v1.1.0
git push origin :refs/tags/v1.1.0

# Delete GitHub release manually
# Fix issues and create new release
```