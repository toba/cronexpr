package cronexpr_test

import (
	"testing"
	"time"

	"github.com/toba/cronexpr"
)

func TestDescribe(t *testing.T) {
	utc := time.UTC
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		// Minute intervals
		{"every 5 minutes", "*/5 * * * *", "Every 5 minutes"},
		{"every 20 minutes", "*/20 * * * *", "Every 20 minutes"},
		{"every minute", "* * * * *", "Every minute"},

		// Specific times
		{"daily at noon", "0 12 * * *", "At 12:00 PM"},
		{"daily at 11pm", "0 23 * * *", "At 11:00 PM"},
		{"at midnight", "0 0 * * *", "At 12:00 AM"},

		// Multiple times
		{"twice daily", "0 12,19 * * *", "At 12:00 PM and 7:00 PM"},

		// Day of week patterns
		{"weekdays at 11pm", "0 23 * * 1-5", "At 11:00 PM, Monday–Friday"},
		{"sunday at 9am", "0 9 * * 0", "At 9:00 AM, Sunday only"},
		{"tue and thu at 2am", "0 2 * * 2,4", "At 2:00 AM, Tuesday and Thursday only"},

		// Day of month patterns
		{"first of month", "0 9 1 * *", "At 9:00 AM, on the 1st of the month"},

		// Hour intervals
		{"every 4 hours", "0 */4 * * *", "At minute 0, every 4 hours"},

		// Minute interval with hour range
		{"every 20 min 7am-9pm", "*/20 7-20 * * *", "Every 20 minutes, 7:00 AM–8:00 PM"},

		// Aliases
		{"@daily", "@daily", "At 12:00 AM"},
		{"@hourly", "@hourly", "At minute 0, every hour"},
		{"@monthly", "@monthly", "At 12:00 AM, on the 1st of the month"},
		{"@yearly", "@yearly", "At 12:00 AM, on the 1st of the month only in January"},
		{"@weekly", "@weekly", "At 12:00 AM, Sunday only"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := cronexpr.MustParse(tc.expr).Describe(&cronexpr.DescribeOptions{
				SourceLocation: utc, TargetLocation: utc,
			})
			if result != tc.expected {
				t.Errorf("Describe(%q) = %q, want %q", tc.expr, result, tc.expected)
			}
		})
	}
}

func TestDescribe_NilOpts(t *testing.T) {
	result := cronexpr.MustParse("0 12 * * *").Describe(nil)
	if result != "At 12:00 PM" {
		t.Errorf("Describe(nil) = %q, want %q", result, "At 12:00 PM")
	}
}

func TestDescribeShort(t *testing.T) {
	utc := time.UTC
	mst := time.FixedZone("MST", -7*60*60)

	tests := []struct {
		name      string
		expr      string
		srcLoc    *time.Location
		targetLoc *time.Location
		expected  string
	}{
		{"weekdays short names", "0 23 * * 1-5", utc, utc, "At 11PM, Mon–Fri"},
		{"specific day short", "0 9 * * 0", utc, utc, "At 9AM, Sunday only"},
		{"multiple days short", "0 2 * * 2,4", utc, utc, "At 2AM, Tue and Thu only"},
		{"day of month short", "0 9 1 * *", utc, utc, "At 9AM, on the 1st"},
		{"day of month with month short", "0 9 1 3 *", utc, utc, "At 9AM, on the 1st only in Mar"},
		{"day list short", "0 9 1,15 * *", utc, utc, "At 9AM, on the 1st and 15th"},
		{"day range short", "0 9 1-15 * *", utc, utc, "At 9AM, days 1–15th"},
		{"last day short", "0 9 L * *", utc, utc, "At 9AM, last day of month"},
		{"weekday nearest short", "0 9 5W * *", utc, utc, "At 9AM, weekday nearest the 5th"},
		{"month range jan-jun short", "0 9 * 1-6 *", utc, utc, "At 9AM, Jan–Jun"},
		{"timezone with short day names", "0 2 * * 2,4", utc, mst, "At 7PM, Mon and Wed only"},
		{"interval unchanged", "*/5 * * * *", utc, utc, "Every 5 minutes"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := cronexpr.MustParse(tc.expr).Describe(&cronexpr.DescribeOptions{
				Short:          true,
				SourceLocation: tc.srcLoc,
				TargetLocation: tc.targetLoc,
			})
			if result != tc.expected {
				t.Errorf("Describe(%q, short) = %q, want %q", tc.expr, result, tc.expected)
			}
		})
	}
}

func TestDescribeTimezone(t *testing.T) {
	utc := time.UTC
	mst := time.FixedZone("MST", -7*60*60)

	tests := []struct {
		name      string
		expr      string
		srcLoc    *time.Location
		targetLoc *time.Location
		expected  string
	}{
		{"9am UTC to MST", "0 9 * * *", utc, mst, "At 2:00 AM"},
		{"2am UTC to MST with day change", "0 2 * * 2,4", utc, mst, "At 7:00 PM, Monday and Wednesday only"},
		{"midnight UTC to MST", "0 0 * * *", utc, mst, "At 5:00 PM"},
		{"interval unchanged", "*/20 * * * *", utc, mst, "Every 20 minutes"},
		{"same timezone", "0 9 * * *", utc, utc, "At 9:00 AM"},
		{"hour range UTC to MST", "0 9-17 * * *", utc, mst, "At minute 0, 2:00 AM–10:00 AM"},
		{"minute interval with hour range same tz", "*/20 7-20 * * *", mst, mst, "Every 20 minutes, 7:00 AM–8:00 PM"},
		{"multiple hours UTC to MST", "0 9,15 * * *", utc, mst, "At 2:00 AM and 8:00 AM"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := cronexpr.MustParse(tc.expr).Describe(&cronexpr.DescribeOptions{
				SourceLocation: tc.srcLoc, TargetLocation: tc.targetLoc,
			})
			if result != tc.expected {
				t.Errorf("Describe(%q, %v→%v) = %q, want %q",
					tc.expr, tc.srcLoc, tc.targetLoc, result, tc.expected)
			}
		})
	}
}
