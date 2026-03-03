# Changelog

## Week of Feb 9 – Feb 15, 2026

### 🗜️ Tweaks

- Rework day-of-week description grammar; "only on Sunday" → "Sunday only"
- Shorten short-mode day-of-month descriptions; drop "on" prefix, use "month" instead of "the month"
- Use ordinal suffixes in day-of-month descriptions; "day 5" → "the 5th"

## Week of Feb 2 – Feb 8, 2026

### 🐞 Fixes

- Fix panic on wrap-around hour ranges like `14-3`; normalize to equivalent union of two ranges

### 🗜️ Tweaks

- Replace all regex usage with string operations; `BenchmarkParse` 3x faster
- Rewrite README.md as library-focused documentation
- Configure Go module versioning; replace GoReleaser with library-appropriate release workflow
- Optimize cronexpr as library-only; remove CLI artifacts
- Simplify tests with table-driven patterns
- Code quality improvements; constants, dedup, consolidation
