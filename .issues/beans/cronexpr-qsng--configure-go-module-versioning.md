---
# cronexpr-qsng
title: Configure Go module versioning
status: completed
type: task
priority: normal
created_at: 2026-02-07T16:46:22Z
updated_at: 2026-02-07T16:46:38Z
---

Replace GoReleaser workflow with library-appropriate release workflow and update CLAUDE.md Go version

## Summary of Changes\n\n- Replaced GoReleaser workflow with a library-appropriate release workflow that runs tests then creates a GitHub Release with auto-generated release notes via softprops/action-gh-release@v2\n- Updated CLAUDE.md to reflect Go 1.25+ (matching go.mod)
