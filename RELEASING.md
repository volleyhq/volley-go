# Releasing Volley Go SDK

This guide explains how to release a new version of the Volley Go SDK.

## Prerequisites

1. Ensure all tests pass:
   ```bash
   go test ./...
   ```

2. Ensure the code compiles:
   ```bash
   go build ./...
   ```

3. Update the version in `version.go`:
   ```go
   const Version = "1.0.0"  // Update to new version
   ```

4. Update `go.mod` if needed (Go version requirements)

## Release Steps

### 1. Create a Git Tag

Tag the release with a semantic version (e.g., `v1.0.0`):

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

**Important**: The tag name must start with `v` and match the version in `version.go` (without the `v` prefix).

### 2. Publish to GitHub

If this is a GitHub repository, the tag will make the release available via:

```bash
go get github.com/volleyhooks/volley-go@v1.0.0
```

### 3. Verify Release

Test that the release can be installed:

```bash
# In a clean directory
go mod init test-release
go get github.com/volleyhooks/volley-go@v1.0.0
go build
```

### 4. Update Documentation

- Update the README if needed
- Create a GitHub release with release notes
- Update any changelog or release notes

## Versioning

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for backwards-compatible bug fixes

## Example Release Workflow

```bash
# 1. Update version.go
# const Version = "1.0.0"

# 2. Commit changes
git add version.go
git commit -m "Bump version to 1.0.0"

# 3. Create and push tag
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
git push origin main  # or master

# 4. Verify
go get github.com/volleyhooks/volley-go@v1.0.0
```

## Notes

- Go modules automatically use the latest tagged version
- Users can pin to a specific version: `go get github.com/volleyhooks/volley-go@v1.0.0`
- Users can use the latest version: `go get github.com/volleyhooks/volley-go@latest`
- The module path in `go.mod` must match the GitHub repository path

