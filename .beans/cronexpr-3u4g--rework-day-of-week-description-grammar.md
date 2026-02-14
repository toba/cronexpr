---
# cronexpr-3u4g
title: Rework day-of-week description grammar
status: completed
type: task
priority: normal
created_at: 2026-02-14T17:03:10Z
updated_at: 2026-02-14T17:04:06Z
---

Move 'only' suffix to end of DOW descriptions and use full day names for single days in short mode. Changes: 'only on Sunday' → 'Sunday only', 'only on Tue and Thu' → 'Tue and Thu only', short single day uses full name.

## Summary of Changes\n\nReworked day-of-week description grammar in `describeDayOfWeek`:\n- Single day: "only on Sunday" → "Sunday only" (always uses full day name, even in short mode)\n- Day list: "only on Tue and Thu" → "Tue and Thu only" (preserves short/full names from caller)\n- Ranges and special patterns (L, #) unchanged\n- Updated 7 test expectations across TestDescribe, TestDescribeShort, and TestDescribeTimezone
