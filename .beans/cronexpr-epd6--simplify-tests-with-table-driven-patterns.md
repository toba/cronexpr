---
# cronexpr-epd6
title: Simplify tests with table-driven patterns
status: completed
type: task
priority: normal
created_at: 2026-02-07T16:57:48Z
updated_at: 2026-02-07T16:59:48Z
---

Add name fields to crontest, refactor TestZero/TestNextN/TestInterval into table-driven tests, fix bugs

## Summary of Changes

All changes in `cronexpr_test.go`:

1. **Added `name` field to `crontest` struct** and populated descriptive names from existing comments. Updated `TestExpressions` to use two-level `t.Run` (outer by name, inner by from-time). Changed `t.Errorf` to `t.Fatalf` on parse error to avoid nil-pointer panic.

2. **Refactored `TestZero`** into table-driven test with 3 cases (`PastYear`, `FutureYear`, `ZeroTime`). Fixed pre-existing copy-paste bug where error message said `2014` but expression was `2099`.

3. **Merged `TestNextN` and `TestNextN_every5min`** into a single table-driven `TestNextN` with two cases (`FifthSaturday`, `Every5Min`). Uses `t.Fatalf` on length mismatch.

4. **Refactored `TestInterval_Interval60Issue`** into table-driven test with 4 named cases.

5. **No changes** to benchmarks or `example_test.go`.
