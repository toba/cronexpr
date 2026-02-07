# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go library for parsing cron expressions and computing the next matching time(s). Fork of `github.com/gorhill/cronexpr`. Module path is `github.com/toba/cronexpr`, requires Go 1.25+. Supports 5-field (standard), 6-field (with seconds), and 7-field (with seconds and year) cron expressions, plus predefined aliases (`@yearly`, `@monthly`, `@weekly`, `@daily`, `@hourly`).

## Build & Test Commands

```bash
go test ./...              # run all tests
go test -run TestZero      # run a single test by name
go test -bench .           # run benchmarks
go test -v                 # verbose test output
```

Tests are in-package (package `cronexpr`), so they have access to unexported types.

## Architecture

Three source files, one exported type (`Expression`), two entry points (`Parse`/`MustParse`):

- **cronexpr.go** — Public API. Defines the `Expression` struct and the `Parse`/`MustParse` constructors, plus `Next`/`NextN` methods. `Next` walks fields from year down to second using `slices.BinarySearch` against pre-computed sorted value lists; on mismatch it delegates to the appropriate `next*` function in cronexpr_next.go.
- **cronexpr_parse.go** — Parsing engine. Tokenizes each cron field using regex patterns, handles special characters (`L`, `W`, `#`), and populates the `Expression` fields. Uses `fieldDescriptor` structs to parameterize parsing per field type. Regex patterns are lazily compiled and cached in a `sync.Map`.
- **cronexpr_next.go** — Time advancement. Contains `nextYear` through `nextSecond` functions that cascade when a field overflows. `calculateActualDaysOfMonth` merges day-of-month and day-of-week constraints per the crontab spec (if both are restricted, either match triggers).

## Cron Expression Semantics

- Day-of-week 7 is normalized to 0 (both mean Sunday)
- 5 fields → seconds default to 0, year defaults to wildcard
- 6 fields → seconds default to 0
- When both day-of-month and day-of-week are restricted (not `*`), a day matches if **either** field matches (union, per crontab spec)
- `W` (nearest weekday) does not cross month boundaries
- Year range is 1970–2099
