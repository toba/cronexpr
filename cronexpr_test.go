package cronexpr

import (
	"testing"
	"time"
)

type crontimes struct {
	from string
	next string
}

type crontest struct {
	name   string
	expr   string
	layout string
	times  []crontimes
}

var crontests = []crontest{
	// Seconds
	{
		name:   "Seconds",
		expr:   "* * * * * * *",
		layout: "2006-01-02 15:04:05",
		times: []crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:00:01"},
			{"2013-01-01 00:00:59", "2013-01-01 00:01:00"},
			{"2013-01-01 00:59:59", "2013-01-01 01:00:00"},
			{"2013-01-01 23:59:59", "2013-01-02 00:00:00"},
			{"2013-02-28 23:59:59", "2013-03-01 00:00:00"},
			{"2016-02-28 23:59:59", "2016-02-29 00:00:00"},
			{"2012-12-31 23:59:59", "2013-01-01 00:00:00"},
		},
	},

	// every 5 Second
	{
		name:   "Every5thSecond",
		expr:   "*/5 * * * * * *",
		layout: "2006-01-02 15:04:05",
		times: []crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:00:05"},
			{"2013-01-01 00:00:59", "2013-01-01 00:01:00"},
			{"2013-01-01 00:59:59", "2013-01-01 01:00:00"},
			{"2013-01-01 23:59:59", "2013-01-02 00:00:00"},
			{"2013-02-28 23:59:59", "2013-03-01 00:00:00"},
			{"2016-02-28 23:59:59", "2016-02-29 00:00:00"},
			{"2012-12-31 23:59:59", "2013-01-01 00:00:00"},
		},
	},

	// Minutes
	{
		name:   "Minutes",
		expr:   "* * * * *",
		layout: "2006-01-02 15:04:05",
		times: []crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:01:00"},
			{"2013-01-01 00:00:59", "2013-01-01 00:01:00"},
			{"2013-01-01 00:59:00", "2013-01-01 01:00:00"},
			{"2013-01-01 23:59:00", "2013-01-02 00:00:00"},
			{"2013-02-28 23:59:00", "2013-03-01 00:00:00"},
			{"2016-02-28 23:59:00", "2016-02-29 00:00:00"},
			{"2012-12-31 23:59:00", "2013-01-01 00:00:00"},
		},
	},

	// Minutes with interval
	{
		name:   "MinutesWithInterval",
		expr:   "17-43/5 * * * *",
		layout: "2006-01-02 15:04:05",
		times: []crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:17:00"},
			{"2013-01-01 00:16:59", "2013-01-01 00:17:00"},
			{"2013-01-01 00:30:00", "2013-01-01 00:32:00"},
			{"2013-01-01 00:50:00", "2013-01-01 01:17:00"},
			{"2013-01-01 23:50:00", "2013-01-02 00:17:00"},
			{"2013-02-28 23:50:00", "2013-03-01 00:17:00"},
			{"2016-02-28 23:50:00", "2016-02-29 00:17:00"},
			{"2012-12-31 23:50:00", "2013-01-01 00:17:00"},
		},
	},

	// Minutes interval, list
	{
		name:   "MinutesIntervalList",
		expr:   "15-30/4,55 * * * *",
		layout: "2006-01-02 15:04:05",
		times: []crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:15:00"},
			{"2013-01-01 00:16:00", "2013-01-01 00:19:00"},
			{"2013-01-01 00:30:00", "2013-01-01 00:55:00"},
			{"2013-01-01 00:55:00", "2013-01-01 01:15:00"},
			{"2013-01-01 23:55:00", "2013-01-02 00:15:00"},
			{"2013-02-28 23:55:00", "2013-03-01 00:15:00"},
			{"2016-02-28 23:55:00", "2016-02-29 00:15:00"},
			{"2012-12-31 23:54:00", "2012-12-31 23:55:00"},
			{"2012-12-31 23:55:00", "2013-01-01 00:15:00"},
		},
	},

	// Days of week
	{
		name:   "DaysOfWeek_MON",
		expr:   "0 0 * * MON",
		layout: "Mon 2006-01-02 15:04",
		times: []crontimes{
			{"2013-01-01 00:00:00", "Mon 2013-01-07 00:00"},
			{"2013-01-28 00:00:00", "Mon 2013-02-04 00:00"},
			{"2013-12-30 00:30:00", "Mon 2014-01-06 00:00"},
		},
	},
	{
		name:   "DaysOfWeek_friday",
		expr:   "0 0 * * friday",
		layout: "Mon 2006-01-02 15:04",
		times: []crontimes{
			{"2013-01-01 00:00:00", "Fri 2013-01-04 00:00"},
			{"2013-01-28 00:00:00", "Fri 2013-02-01 00:00"},
			{"2013-12-30 00:30:00", "Fri 2014-01-03 00:00"},
		},
	},
	{
		name:   "DaysOfWeek_6and7",
		expr:   "0 0 * * 6,7",
		layout: "Mon 2006-01-02 15:04",
		times: []crontimes{
			{"2013-01-01 00:00:00", "Sat 2013-01-05 00:00"},
			{"2013-01-28 00:00:00", "Sat 2013-02-02 00:00"},
			{"2013-12-30 00:30:00", "Sat 2014-01-04 00:00"},
		},
	},

	// Specific days of week
	{
		name:   "SpecificDayOfWeek_6hash5",
		expr:   "0 0 * * 6#5",
		layout: "Mon 2006-01-02 15:04",
		times: []crontimes{
			{"2013-09-02 00:00:00", "Sat 2013-11-30 00:00"},
		},
	},

	// Work day of month
	{
		name:   "WorkDayOfMonth_14W",
		expr:   "0 0 14W * *",
		layout: "Mon 2006-01-02 15:04",
		times: []crontimes{
			{"2013-03-31 00:00:00", "Mon 2013-04-15 00:00"},
			{"2013-08-31 00:00:00", "Fri 2013-09-13 00:00"},
		},
	},

	// Work day of month -- end of month
	{
		name:   "WorkDayOfMonth_30W",
		expr:   "0 0 30W * *",
		layout: "Mon 2006-01-02 15:04",
		times: []crontimes{
			{"2013-03-02 00:00:00", "Fri 2013-03-29 00:00"},
			{"2013-06-02 00:00:00", "Fri 2013-06-28 00:00"},
			{"2013-09-02 00:00:00", "Mon 2013-09-30 00:00"},
			{"2013-11-02 00:00:00", "Fri 2013-11-29 00:00"},
		},
	},

	// Last day of month
	{
		name:   "LastDayOfMonth",
		expr:   "0 0 L * *",
		layout: "Mon 2006-01-02 15:04",
		times: []crontimes{
			{"2013-09-02 00:00:00", "Mon 2013-09-30 00:00"},
			{"2014-01-01 00:00:00", "Fri 2014-01-31 00:00"},
			{"2014-02-01 00:00:00", "Fri 2014-02-28 00:00"},
			{"2016-02-15 00:00:00", "Mon 2016-02-29 00:00"},
		},
	},

	// Last work day of month
	{
		name:   "LastWorkDayOfMonth",
		expr:   "0 0 LW * *",
		layout: "Mon 2006-01-02 15:04",
		times: []crontimes{
			{"2013-09-02 00:00:00", "Mon 2013-09-30 00:00"},
			{"2013-11-02 00:00:00", "Fri 2013-11-29 00:00"},
			{"2014-08-15 00:00:00", "Fri 2014-08-29 00:00"},
		},
	},

	// Wrap-around hour range (14-3 means 14:00 through 03:00, wrapping past midnight)
	{
		name:   "WrapHourRange",
		expr:   "0 14-3 * * *",
		layout: "2006-01-02 15:04:05",
		times: []crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 01:00:00"},
			{"2013-01-01 03:00:00", "2013-01-01 14:00:00"},
			{"2013-01-01 13:00:00", "2013-01-01 14:00:00"},
			{"2013-01-01 14:00:00", "2013-01-01 15:00:00"},
			{"2013-01-01 23:00:00", "2013-01-02 00:00:00"},
			{"2013-01-01 02:00:00", "2013-01-01 03:00:00"},
		},
	},

	// Wrap-around hour range with step (22-4/2 means 22,0,2,4)
	{
		name:   "WrapHourRangeWithStep",
		expr:   "0 22-4/2 * * *",
		layout: "2006-01-02 15:04:05",
		times: []crontimes{
			{"2013-01-01 21:00:00", "2013-01-01 22:00:00"},
			{"2013-01-01 22:00:00", "2013-01-02 00:00:00"},
			{"2013-01-02 00:00:00", "2013-01-02 02:00:00"},
			{"2013-01-02 02:00:00", "2013-01-02 04:00:00"},
			{"2013-01-02 04:00:00", "2013-01-02 22:00:00"},
		},
	},

	// Wrap-around minute range (45-15 means :45 through :15, wrapping past the hour)
	{
		name:   "WrapMinuteRange",
		expr:   "45-15 * * * *",
		layout: "2006-01-02 15:04:05",
		times: []crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:01:00"},
			{"2013-01-01 00:15:00", "2013-01-01 00:45:00"},
			{"2013-01-01 00:44:00", "2013-01-01 00:45:00"},
			{"2013-01-01 00:59:00", "2013-01-01 01:00:00"},
		},
	},

	// Wrap-around day-of-week range (FRI-MON means Fri, Sat, Sun, Mon)
	{
		name:   "WrapDayOfWeekRange",
		expr:   "0 0 * * 5-1",
		layout: "Mon 2006-01-02 15:04",
		times: []crontimes{
			{"2013-01-01 00:00:00", "Fri 2013-01-04 00:00"},
			{"2013-01-04 00:00:00", "Sat 2013-01-05 00:00"},
			{"2013-01-05 00:00:00", "Sun 2013-01-06 00:00"},
			{"2013-01-06 00:00:00", "Mon 2013-01-07 00:00"},
			{"2013-01-07 00:00:00", "Fri 2013-01-11 00:00"},
		},
	},

	// Wrap-around month range (OCT-FEB means Oct, Nov, Dec, Jan, Feb)
	{
		name:   "WrapMonthRange",
		expr:   "0 0 1 10-2 *",
		layout: "2006-01-02 15:04:05",
		times: []crontimes{
			{"2013-01-01 00:00:00", "2013-02-01 00:00:00"},
			{"2013-02-01 00:00:00", "2013-10-01 00:00:00"},
			{"2013-10-01 00:00:00", "2013-11-01 00:00:00"},
			{"2013-12-01 00:00:00", "2014-01-01 00:00:00"},
		},
	},

	// TODO: more tests
}

func TestExpressions(t *testing.T) {
	for _, test := range crontests {
		t.Run(test.name, func(t *testing.T) {
			expr, err := Parse(test.expr)
			if err != nil {
				t.Fatalf(`Parse("%s") returned "%s"`, test.expr, err.Error())
			}
			for _, times := range test.times {
				t.Run(times.from, func(t *testing.T) {
					from, _ := time.Parse("2006-01-02 15:04:05", times.from)
					next := expr.Next(from)
					nextstr := next.Format(test.layout)
					if nextstr != times.next {
						t.Errorf(`("%s").Next("%s") = "%s", got "%s"`, test.expr, times.from, times.next, nextstr)
					}
				})
			}
		})
	}
}

func TestZero(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		from     string
		wantZero bool
	}{
		{"PastYear", "* * * * * 1980", "2013-08-31", true},
		{"FutureYear", "* * * * * 2050", "2013-08-31", false},
		{"ZeroTime", "* * * * * 2099", "0001-01-01", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var from time.Time
			if tt.from != "0001-01-01" {
				from, _ = time.Parse("2006-01-02", tt.from)
			}
			next := MustParse(tt.expr).Next(from)
			if next.IsZero() != tt.wantZero {
				t.Errorf(`("%s").Next("%s").IsZero() = %v, want %v`, tt.expr, tt.from, next.IsZero(), tt.wantZero)
			}
		})
	}
}

func TestNextN(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		from     string
		layout   string
		expected []string
	}{
		{
			name:   "FifthSaturday",
			expr:   "0 0 * * 6#5",
			from:   "2013-09-02 08:44:30",
			layout: "Mon, 2 Jan 2006 15:04:15",
			expected: []string{
				"Sat, 30 Nov 2013 00:00:00",
				"Sat, 29 Mar 2014 00:00:00",
				"Sat, 31 May 2014 00:00:00",
				"Sat, 30 Aug 2014 00:00:00",
				"Sat, 29 Nov 2014 00:00:00",
			},
		},
		{
			name:   "Every5Min",
			expr:   "*/5 * * * *",
			from:   "2013-09-02 08:44:32",
			layout: "Mon, 2 Jan 2006 15:04:05",
			expected: []string{
				"Mon, 2 Sep 2013 08:45:00",
				"Mon, 2 Sep 2013 08:50:00",
				"Mon, 2 Sep 2013 08:55:00",
				"Mon, 2 Sep 2013 09:00:00",
				"Mon, 2 Sep 2013 09:05:00",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, _ := time.Parse("2006-01-02 15:04:05", tt.from)
			result := MustParse(tt.expr).NextN(from, uint(len(tt.expected)))
			if len(result) != len(tt.expected) {
				t.Fatalf("got %d results, want %d", len(result), len(tt.expected))
			}
			for i, next := range result {
				nextStr := next.Format(tt.layout)
				if nextStr != tt.expected[i] {
					t.Errorf("result[%d] = %q, want %q", i, nextStr, tt.expected[i])
				}
			}
		})
	}
}

// Issue: https://github.com/toba/cronexpr/issues/16
func TestInterval_Interval60Issue(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{"StarSlash60", "*/60 * * * * *"},
		{"StarSlash61", "*/61 * * * * *"},
		{"2Slash60", "2/60 * * * * *"},
		{"RangeSlash61", "2-20/61 * * * * *"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.expr)
			if err == nil {
				t.Errorf("Parse(%q) should return error for invalid interval", tt.expr)
			}
		})
	}
}

var benchmarkExpressions = []string{
	"* * * * *",
	"@hourly",
	"@weekly",
	"@yearly",
	"30 3 15W 3/3 *",
	"30 0 0 1-31/5 Oct-Dec * 2000,2006,2008,2013-2015",
	"0 0 0 * Feb-Nov/2 thu#3 2000-2050",
}
var benchmarkExpressionsLen = len(benchmarkExpressions)

func BenchmarkParse(b *testing.B) {
	for i := range b.N {
		_ = MustParse(benchmarkExpressions[i%benchmarkExpressionsLen])
	}
}

func BenchmarkNext(b *testing.B) {
	exprs := make([]*Expression, benchmarkExpressionsLen)
	for i := range benchmarkExpressionsLen {
		exprs[i] = MustParse(benchmarkExpressions[i])
	}
	from := time.Now()
	b.ResetTimer()
	for i := range b.N {
		expr := exprs[i%benchmarkExpressionsLen]
		next := expr.Next(from)
		next = expr.Next(next)
		next = expr.Next(next)
		next = expr.Next(next)
		_ = expr.Next(next)
	}
}
