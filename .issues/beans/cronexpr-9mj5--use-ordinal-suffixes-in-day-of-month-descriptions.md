---
# cronexpr-9mj5
title: Use ordinal suffixes in day-of-month descriptions
status: completed
type: task
priority: normal
created_at: 2026-02-14T16:57:31Z
updated_at: 2026-02-14T16:58:20Z
---

Update describeDayOfMonth to use ordinal suffixes (1st, 2nd, 3rd, etc.) instead of plain numbers for both short and long modes.

## Todo
- [x] Add descOrdinal helper function
- [x] Rewrite describeDayOfMonth to use ordinals
- [x] Update long-mode test expectations
- [x] Update short-mode test expectations
- [x] Verify all tests pass

## Summary of Changes

Added `descOrdinal` helper that returns integers with English ordinal suffixes (1st, 2nd, 3rd, 11th, 12th, 13th, 21st, etc.). Rewrote `describeDayOfMonth` to use ordinals throughout: single days ("on the 5th"), lists ("on the 1st and 15th"), ranges (long: "on the 1st–15th of the month", short: "days 1–15th"), and W patterns ("weekday nearest the 5th"). L and interval patterns unchanged. Updated all affected test expectations in both long and short modes.
