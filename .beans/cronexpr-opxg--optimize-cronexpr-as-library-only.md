---
# cronexpr-opxg
title: Optimize cronexpr as library-only
status: completed
type: task
priority: normal
created_at: 2026-02-07T16:26:58Z
updated_at: 2026-02-07T16:27:34Z
---

Remove CLI artifacts: delete .goreleaser.yaml, remove dist/ from .gitignore, remove dead expression field, update README.md and CLAUDE.md to remove CLI references.

## Summary of Changes

- Deleted `.goreleaser.yaml` (no binary to release for a library)
- Removed `dist/` and goreleaser comment from `.gitignore`
- Removed dead `expression string` field from `Expression` struct in `cronexpr.go`
- Removed CLI reference sentence from `README.md`
- Removed `go build ./cronexpr` and CLI subdirectory reference from `CLAUDE.md`
- Verified: `go vet`, `go test`, and `go test -race` all pass
