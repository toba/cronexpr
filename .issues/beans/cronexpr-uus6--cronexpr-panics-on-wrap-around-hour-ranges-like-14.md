---
# cronexpr-uus6
title: cronexpr panics on wrap-around hour ranges like 14-3
status: completed
type: bug
priority: normal
created_at: 2026-02-07T16:28:17Z
updated_at: 2026-02-07T16:32:41Z
---

## Problem

gorhill/cronexpr panics on wrap-around hour ranges like `14-3` (meaning 14 through 3, wrapping past midnight). `NextRunTime` recovers the panic but returns an error, causing 'Never' display on job detail page.

## Expected Behavior

Wrap-around ranges like `14-3` in the hour field should be interpreted as `14-23,0-3` (14 through midnight, then midnight through 3).

## Fix

Normalize wrap-around ranges before passing to cronexpr. When the start of a range is greater than the end (e.g. `14-3`), expand it into the equivalent union of two ranges (e.g. `14-23,0-3`).

## Summary of Changes

Fixed in `cronexpr_parse.go`:

- Modified `populateMany()` to accept optional field bounds (`fieldMin`, `fieldMax`) and handle wrap-around ranges where `lo > hi` (e.g. hour `14-3` → populate 14..23 then 0..3)
- Updated all call sites (`genericFieldHandler`, `dowFieldHandler`, `domFieldHandler`) to pass field descriptor bounds

Tests added in `cronexpr_test.go`:

- Wrap-around hour range: `0 14-3 * * *` (6 cases)
- Wrap-around hour range with step: `0 22-4/2 * * *` (5 cases)
- Wrap-around minute range: `45-15 * * * *` (4 cases)
- Wrap-around day-of-week range: `0 0 * * 5-1` (5 cases)
- Wrap-around month range: `0 0 1 10-2 *` (4 cases)
