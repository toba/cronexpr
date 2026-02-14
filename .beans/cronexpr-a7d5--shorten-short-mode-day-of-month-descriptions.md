---
# cronexpr-a7d5
title: Shorten short-mode day-of-month descriptions
status: completed
type: task
priority: normal
created_at: 2026-02-14T16:47:08Z
updated_at: 2026-02-14T16:50:57Z
---

When `DescribeOptions.Short` is true, the day-of-month descriptions should be more compact:

| Long (current) | Short (current) | Short (desired) |
|---|---|---|
| on days 1 and 15 of the month | on days 1 and 15 of the month | 1 and 15 of month |
| on day 5 of the month | on day 5 of the month | day 5 of month |
| on days 1–15 of the month | on days 1–15 of the month | days 1–15 of month |
| on the last day of the month | on the last day of the month | last day of month |
| on the weekday nearest day 5 of the month | on the weekday nearest day 5 of the month | weekday nearest day 5 of month |

## Changes needed

- [x] Add `short bool` parameter to `describeDayOfMonth`
- [x] Pass `opts.Short` through `describeDate` → `describeDayOfMonth`
- [x] In short mode, drop "on " / "on the " prefix and use "month" instead of "the month"
- [x] Update tests in `cronexpr_describe_test.go`

Triggered by Pacer core's job dashboard showing "At 4AM, on days 1 and 15 of the month" for the snapshot-reviews job — too long for the short rendering.

## Summary of Changes

Added `short bool` parameter to `describeDayOfMonth` and threaded it through `describeDate`. In short mode, day-of-month descriptions drop the "on "/"on the " prefix and use "month" instead of "the month" (e.g. "on days 1 and 15 of the month" → "1 and 15 of month"). Updated existing short-mode tests and added new cases for day list, day range, last day, and weekday-nearest patterns.
