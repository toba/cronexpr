#!/bin/bash
set -e

# Pre-commit checks
echo "==> Running pre-commit checks..."
golangci-lint run
go test ./...

# Stage and show changes
echo "==> Staging changes..."
git add -A
git status --short
echo ""
echo "==> Staged diff:"
git diff --staged

# Get commit message from arguments
if [ -z "$1" ]; then
    echo ""
    echo "ERROR: Commit subject required as first argument"
    exit 1
fi

SUBJECT="$1"
DESCRIPTION="${2:-}"

# Build commit message
if [ -n "$DESCRIPTION" ]; then
    COMMIT_MSG="$SUBJECT

$DESCRIPTION"
else
    COMMIT_MSG="$SUBJECT"
fi

# Create commit
echo ""
echo "==> Creating commit..."
git commit -m "$COMMIT_MSG"
git status

# Push and release
CURRENT_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [ -n "$CURRENT_TAG" ]; then
    echo ""
    echo "==> Current version: $CURRENT_TAG"
fi

if [ "$PUSH" = "true" ]; then
    echo "==> Pushing commits..."
    git push

    # If NEW_VERSION is set, create and push tag
    # GitHub Actions release workflow will create the release automatically
    if [ -n "$NEW_VERSION" ]; then
        echo "==> Creating tag $NEW_VERSION..."
        git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"

        echo "==> Pushing tag (GitHub Actions will create release)..."
        git push origin "$NEW_VERSION"
        echo "==> Tag $NEW_VERSION pushed, release workflow will create GitHub release"
    fi
else
    echo "==> Commit is local only (use PUSH=true to push and release)"
    if [ -n "$NEW_VERSION" ]; then
        echo "==> NEW_VERSION=$NEW_VERSION will be used when pushed"
    fi
fi
echo ""
echo "==> Done!"
