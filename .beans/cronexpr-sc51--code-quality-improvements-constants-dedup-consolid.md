---
# cronexpr-sc51
title: 'Code quality improvements: constants, dedup, consolidation'
status: completed
type: task
priority: normal
created_at: 2026-02-07T16:52:50Z
updated_at: 2026-02-07T16:54:57Z
---

Behavior-preserving refactors: fix duplicate month regex, extract year range constants, field count constants, daysPerWeek constant, replace manual weekday index with loop, consolidate field handlers, extract interval validation helper.

## Summary of Changes

All 7 planned refactors applied across 3 files:

1. **Fixed duplicate month regex** — `march|april` duplicated where `may` should be (cronexpr_parse.go:121)
2. **Extracted year range constants** — `minYear`/`maxYear` constants, `makeIntRange` helper replaces 50+ hardcoded lines, `numberTokens` generated via IIFE (cronexpr_parse.go:13-44)
3. **Extracted field count constants** — `minCronFields`/`maxCronFields` replace magic 5/7 (cronexpr.go:46-49)
4. **Extracted `daysPerWeek` constant** — replaces 5 magic `7`s in weekday arithmetic (cronexpr_next.go:8)
5. **Replaced manual weekday index access with loop** — fragile `w[0]..w[4]` replaced with `for _, day := range w` (cronexpr_next.go:185-189)
6. **Consolidated 5 field handlers** — single `parseField()` replaces `secondFieldHandler`, `minuteFieldHandler`, `hourFieldHandler`, `monthFieldHandler`, `yearFieldHandler` (cronexpr_parse.go:174-178)
7. **Extracted `validateStep` helper** — replaces 3 identical interval validation blocks (cronexpr_parse.go:336-341)
