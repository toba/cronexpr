# cronexpr

A Go library for parsing cron expressions and computing the next matching time(s). Fork of [gorhill/cronexpr](https://github.com/gorhill/cronexpr) with the following changes:

- Modernize
- Fix panic on wrap-around ranges (e.g. `14-3` for hours) by correctly handling ranges where start > end
- Replace all regex-based parsing with `strings.Cut`/map lookups, eliminating the `regexp` and `sync` dependencies (~3x faster parsing, 44% less memory)

## Install

Requires Go 1.25+.

```
go get github.com/toba/cronexpr
```

## Usage

```go
import (
    "fmt"
    "time"

    "github.com/toba/cronexpr"
)
```

Parse an expression and get the next matching time:

```go
expr := cronexpr.MustParse("0 0 29 2 *")
next := expr.Next(time.Now())
fmt.Println(next)
```

`Parse` returns an error instead of panicking:

```go
expr, err := cronexpr.Parse("0 0 29 2 *")
if err != nil {
    log.Fatal(err)
}
```

Get the next _n_ matching times with `NextN`:

```go
nextTimes := cronexpr.MustParse("0 0 29 2 *").NextN(time.Now(), 5)
for _, t := range nextTimes {
    fmt.Println(t)
}
```

A zero time is returned when no future match exists. Use `IsZero` to check:

```go
next := cronexpr.MustParse("* * * * * 1980").Next(time.Now())
if next.IsZero() {
    fmt.Println("no matching time")
}
```

The time zone of returned times always matches the time zone of the input.

## Supported formats

| Format   | Fields                                                     |
| -------- | ---------------------------------------------------------- |
| 5 fields | minute, hour, day-of-month, month, day-of-week             |
| 6 fields | minute, hour, day-of-month, month, day-of-week, year       |
| 7 fields | second, minute, hour, day-of-month, month, day-of-week, year |

When 5 fields are given, seconds default to `0` and year defaults to `*`. When 6 fields are given, seconds default to `0`.

## Extensions

Beyond the standard `*`, `,`, `-`, and `/` operators, the following are supported:

- **`L`** in day-of-month -- last day of the month
- **`L`** in day-of-week -- last occurrence of that weekday in the month (e.g. `5L` = last Friday)
- **`W`** in day-of-month -- nearest weekday to the given day (e.g. `15W`); single days only
- **`LW`** in day-of-month -- last weekday (Mon-Fri) of the month
- **`#`** in day-of-week -- nth occurrence of a weekday (e.g. `5#3` = third Friday)
- **Wrap-around ranges** -- ranges where start > end wrap through the field boundary (e.g. `22-3` for hours means 22, 23, 0, 1, 2, 3)
- **Month names** -- `JAN`-`DEC` (case-insensitive)
- **Day-of-week names** -- `SUN`-`SAT` (case-insensitive); `7` is accepted as Sunday

When both day-of-month and day-of-week are restricted (not `*`), a day matches if **either** field matches (union semantics, per the crontab spec).

## Aliases

| Alias                  | Equivalent              |
| ---------------------- | ----------------------- |
| `@yearly`, `@annually` | `0 0 0 1 1 * *`        |
| `@monthly`             | `0 0 0 1 * * *`        |
| `@weekly`              | `0 0 0 * * 0 *`        |
| `@daily`               | `0 0 0 * * * *`        |
| `@hourly`              | `0 0 * * * * *`        |

## Limitations

- Year range is 1970-2099
- `@reboot` is not supported
- `W` (nearest weekday) only accepts a single day value, not a range or list, and does not cross month boundaries

## License

[Apache License 2.0](http://www.apache.org/licenses/LICENSE-2.0)
