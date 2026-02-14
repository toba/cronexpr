---
# cronexpr-0s65
title: Replace all regex usage with string operations
status: completed
type: task
priority: normal
created_at: 2026-02-07T17:06:30Z
updated_at: 2026-02-07T17:09:08Z
---

Replace all regex-based parsing in cronexpr_parse.go with strings.Cut/strings.HasSuffix + map lookups. This eliminates the regexp and sync imports entirely and should measurably speed up BenchmarkParse.

## Tasks
- [x] Change fieldDescriptor.atoi signature to func(string) (int, bool)
- [x] Update all 7 descriptors with new atoi signature
- [x] Replace fieldFinder in cronexpr.go with strings.Fields
- [x] Replace entryFinder with splitEntries
- [x] Rewrite genericFieldParse with string operations
- [x] Rewrite domFieldHandler case none branch
- [x] Rewrite dowFieldHandler case none branch
- [x] Remove dead code (makeLayoutRegexp, layout vars, etc.)
- [x] Remove regexp and sync imports
- [x] Run tests to verify correctness
- [x] Run benchmarks to measure improvement

## Summary of Changes

Replaced all regex-based parsing with string operations (strings.Cut, strings.HasSuffix, map lookups, strconv.Atoi). Removed regexp, sync imports and all layout pattern variables. BenchmarkParse improved from ~7700 ns/op to ~2500 ns/op (3x faster), allocations from 74 to 46, memory from 5880 to 3266 B/op.
