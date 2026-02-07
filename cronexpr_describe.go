package cronexpr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DescribeOptions controls how a cron expression is described.
type DescribeOptions struct {
	// Short uses abbreviated names ("Mon" vs "Monday", "Jan" vs "January").
	Short bool
	// SourceLocation is the cron schedule's timezone (nil = UTC).
	SourceLocation *time.Location
	// TargetLocation is the display timezone (nil = UTC).
	TargetLocation *time.Location
}

var (
	descDayNames      = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	descDayShortNames = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	descMonthNames    = []string{"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}
	descMonthShortNames = []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
)

// Describe returns a human-readable description of the cron expression.
// If opts is nil, defaults are used (long names, UTC timezone).
func (expr *Expression) Describe(opts *DescribeOptions) string {
	if opts == nil {
		opts = &DescribeOptions{}
	}
	srcLoc := opts.SourceLocation
	if srcLoc == nil {
		srcLoc = time.UTC
	}
	targetLoc := opts.TargetLocation
	if targetLoc == nil {
		targetLoc = time.UTC
	}

	fields := descParseFields(expr.normalized)
	if fields == nil {
		return expr.normalized
	}

	dNames := descDayNames
	mNames := descMonthNames
	if opts.Short {
		dNames = descDayShortNames
		mNames = descMonthShortNames
	}

	var parts []string

	timeDesc, dayOffset := describeTime(fields, srcLoc, targetLoc)
	if timeDesc != "" {
		parts = append(parts, timeDesc)
	}

	dateDesc := describeDate(fields, dayOffset, dNames, mNames)
	if dateDesc != "" {
		parts = append(parts, dateDesc)
	}

	if len(parts) == 0 {
		return "Every minute"
	}

	return strings.Join(parts, ", ")
}

// descFields holds the raw field strings from the normalized cron expression.
type descFields struct {
	seconds    string
	minutes    string
	hours      string
	dayOfMonth string
	month      string
	dayOfWeek  string
}

// descParseFields splits the normalized cron string into fields.
// Returns nil if the field count is unexpected.
func descParseFields(normalized string) *descFields {
	parts := strings.Fields(normalized)

	var f descFields
	switch len(parts) {
	case 5: // minute hour dom month dow
		f.seconds = "0"
		f.minutes = parts[0]
		f.hours = parts[1]
		f.dayOfMonth = parts[2]
		f.month = parts[3]
		f.dayOfWeek = parts[4]
	case 6: // second minute hour dom month dow
		f.seconds = parts[0]
		f.minutes = parts[1]
		f.hours = parts[2]
		f.dayOfMonth = parts[3]
		f.month = parts[4]
		f.dayOfWeek = parts[5]
	case 7: // second minute hour dom month dow year
		f.seconds = parts[0]
		f.minutes = parts[1]
		f.hours = parts[2]
		f.dayOfMonth = parts[3]
		f.month = parts[4]
		f.dayOfWeek = parts[5]
		// year ignored for description
	default:
		return nil
	}

	// Normalize ? to *
	f.dayOfMonth = strings.ReplaceAll(f.dayOfMonth, "?", "*")
	f.dayOfWeek = descNormalizeDow(f.dayOfWeek)
	f.month = descNormalizeMonth(f.month)
	f.minutes = strings.ReplaceAll(f.minutes, "?", "*")
	f.hours = strings.ReplaceAll(f.hours, "?", "*")

	return &f
}

var descDayMap = map[string]string{
	"sun": "0", "mon": "1", "tue": "2", "wed": "3",
	"thu": "4", "fri": "5", "sat": "6",
}
var descMonthMap = map[string]string{
	"jan": "1", "feb": "2", "mar": "3", "apr": "4",
	"may": "5", "jun": "6", "jul": "7", "aug": "8",
	"sep": "9", "oct": "10", "nov": "11", "dec": "12",
}

func descNormalizeDow(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "?", "*")
	for name, num := range descDayMap {
		s = strings.ReplaceAll(s, name, num)
	}
	if s == "7" {
		s = "0"
	}
	return s
}

func descNormalizeMonth(s string) string {
	s = strings.ToLower(s)
	for name, num := range descMonthMap {
		s = strings.ReplaceAll(s, name, num)
	}
	return strings.ReplaceAll(s, "?", "*")
}

// describeTime returns the time description and a day offset (-1, 0, or 1) for TZ conversion.
func describeTime(f *descFields, srcLoc, targetLoc *time.Location) (string, int) {
	// Intervals are timezone-agnostic
	if f.hours == "*" && f.minutes != "*" {
		if interval, ok := descParseInterval(f.minutes); ok {
			if interval == 1 {
				return "Every minute", 0
			}
			return fmt.Sprintf("Every %d minutes", interval), 0
		}
		// Specific minute, every hour (e.g. @hourly = "0 * * * *")
		minDesc := describeMinutes(f.minutes)
		return minDesc + ", every hour", 0
	}

	if f.minutes == "*" && f.hours != "*" {
		if interval, ok := descParseInterval(f.hours); ok {
			if interval == 1 {
				return "Every minute, every hour", 0
			}
			return fmt.Sprintf("Every minute, every %d hours", interval), 0
		}
	}

	if f.hours == "*" && f.minutes == "*" {
		return "Every minute", 0
	}

	if interval, ok := descParseInterval(f.hours); ok {
		minDesc := describeMinutes(f.minutes)
		if interval == 1 {
			return minDesc + ", every hour", 0
		}
		return fmt.Sprintf("%s, every %d hours", minDesc, interval), 0
	}

	// Specific times — need timezone conversion
	if descIsList(f.hours) {
		hours := descSplitList(f.hours)
		var times []string
		var dayOffset int
		for _, h := range hours {
			t, offset := descFormatTimeWithTZ(h, f.minutes, srcLoc, targetLoc)
			times = append(times, t)
			dayOffset = offset
		}
		return "At " + descJoinWithAnd(times), dayOffset
	}

	if descIsRange(f.hours) {
		start, end := descParseRange(f.hours)
		minDesc := describeMinutes(f.minutes)
		startFmt, dayOffset := descFormatHourWithTZ(start, srcLoc, targetLoc)
		endFmt, _ := descFormatHourWithTZ(end, srcLoc, targetLoc)
		return fmt.Sprintf("%s, %s–%s", minDesc, startFmt, endFmt), dayOffset
	}

	t, dayOffset := descFormatTimeWithTZ(f.hours, f.minutes, srcLoc, targetLoc)
	return "At " + t, dayOffset
}

func describeMinutes(minutes string) string {
	if minutes == "*" || minutes == "0" {
		return "At minute 0"
	}
	if interval, ok := descParseInterval(minutes); ok {
		if interval == 1 {
			return "Every minute"
		}
		return fmt.Sprintf("Every %d minutes", interval)
	}
	min, _ := strconv.Atoi(minutes)
	return fmt.Sprintf("At minute %d", min)
}

// describeDate generates date/day description, adjusting DOW by dayOffset for TZ conversion.
func describeDate(f *descFields, dayOffset int, dNames, mNames []string) string {
	var parts []string

	dowDesc := describeDayOfWeek(f.dayOfWeek, dayOffset, dNames)
	domDesc := describeDayOfMonth(f.dayOfMonth)
	monthDesc := describeMonth(f.month, mNames)

	if domDesc != "" && dowDesc != "" {
		if monthDesc != "" {
			parts = append(parts, domDesc, "and "+dowDesc, monthDesc)
		} else {
			parts = append(parts, domDesc, "and "+dowDesc)
		}
	} else if domDesc != "" {
		if monthDesc != "" {
			parts = append(parts, domDesc, monthDesc)
		} else {
			parts = append(parts, domDesc)
		}
	} else if dowDesc != "" {
		parts = append(parts, dowDesc)
		if monthDesc != "" {
			parts = append(parts, monthDesc)
		}
	} else if monthDesc != "" {
		parts = append(parts, monthDesc)
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, " ")
}

func describeDayOfWeek(dow string, dayOffset int, names []string) string {
	if dow == "*" {
		return ""
	}

	// Last DOW pattern (e.g., 5L = last Friday)
	if strings.HasSuffix(dow, "l") || strings.HasSuffix(dow, "L") {
		day := strings.TrimSuffix(strings.TrimSuffix(dow, "l"), "L")
		d, _ := strconv.Atoi(day)
		if d >= 0 && d <= 6 {
			d = descAdjustDay(d, dayOffset)
			return fmt.Sprintf("on the last %s of the month", names[d])
		}
	}

	// Nth DOW pattern (e.g., 1#2 = second Monday)
	if strings.Contains(dow, "#") {
		parts := strings.Split(dow, "#")
		if len(parts) == 2 {
			day, _ := strconv.Atoi(parts[0])
			nth, _ := strconv.Atoi(parts[1])
			if day >= 0 && day <= 6 && nth >= 1 && nth <= 5 {
				day = descAdjustDay(day, dayOffset)
				ordinal := []string{"", "first", "second", "third", "fourth", "fifth"}[nth]
				return fmt.Sprintf("on the %s %s of the month", ordinal, names[day])
			}
		}
	}

	if descIsRange(dow) {
		start, end := descParseRange(dow)
		startDay, _ := strconv.Atoi(start)
		endDay, _ := strconv.Atoi(end)
		if startDay >= 0 && startDay <= 6 && endDay >= 0 && endDay <= 6 {
			startDay = descAdjustDay(startDay, dayOffset)
			endDay = descAdjustDay(endDay, dayOffset)
			return fmt.Sprintf("%s–%s", names[startDay], names[endDay])
		}
	}

	if descIsList(dow) {
		days := descSplitList(dow)
		var dayNamesList []string
		for _, d := range days {
			if n, err := strconv.Atoi(d); err == nil && n >= 0 && n <= 6 {
				n = descAdjustDay(n, dayOffset)
				dayNamesList = append(dayNamesList, names[n])
			}
		}
		return "only on " + descJoinWithAnd(dayNamesList)
	}

	d, err := strconv.Atoi(dow)
	if err == nil && d >= 0 && d <= 6 {
		d = descAdjustDay(d, dayOffset)
		return "only on " + names[d]
	}

	return ""
}

func describeDayOfMonth(dom string) string {
	if dom == "*" {
		return ""
	}

	if strings.ToLower(dom) == "l" {
		return "on the last day of the month"
	}

	if strings.HasSuffix(strings.ToUpper(dom), "W") {
		day := strings.TrimSuffix(strings.TrimSuffix(dom, "w"), "W")
		return fmt.Sprintf("on the weekday nearest day %s of the month", day)
	}

	if descIsRange(dom) {
		start, end := descParseRange(dom)
		return fmt.Sprintf("on days %s–%s of the month", start, end)
	}

	if descIsList(dom) {
		days := descSplitList(dom)
		return "on days " + descJoinWithAnd(days) + " of the month"
	}

	if interval, ok := descParseInterval(dom); ok {
		if interval == 1 {
			return "every day"
		}
		return fmt.Sprintf("every %d days", interval)
	}

	return fmt.Sprintf("on day %s of the month", dom)
}

func describeMonth(month string, names []string) string {
	if month == "*" {
		return ""
	}

	if descIsRange(month) {
		start, end := descParseRange(month)
		startMonth, _ := strconv.Atoi(start)
		endMonth, _ := strconv.Atoi(end)
		if startMonth >= 1 && startMonth <= 12 && endMonth >= 1 && endMonth <= 12 {
			return fmt.Sprintf("%s–%s", names[startMonth], names[endMonth])
		}
	}

	if descIsList(month) {
		months := descSplitList(month)
		var monthNamesList []string
		for _, m := range months {
			if n, err := strconv.Atoi(m); err == nil && n >= 1 && n <= 12 {
				monthNamesList = append(monthNamesList, names[n])
			}
		}
		return "in " + descJoinWithAnd(monthNamesList)
	}

	m, err := strconv.Atoi(month)
	if err == nil && m >= 1 && m <= 12 {
		return "only in " + names[m]
	}

	return ""
}

// Helpers

var descIntervalRe = regexp.MustCompile(`^\*/(\d+)$`)

func descParseInterval(s string) (int, bool) {
	matches := descIntervalRe.FindStringSubmatch(s)
	if len(matches) == 2 {
		n, _ := strconv.Atoi(matches[1])
		return n, true
	}
	return 0, false
}

func descIsRange(s string) bool {
	return strings.Contains(s, "-") && !strings.HasPrefix(s, "-")
}

func descParseRange(s string) (string, string) {
	if before, after, ok := strings.Cut(s, "-"); ok {
		return before, after
	}
	return s, s
}

func descIsList(s string) bool {
	return strings.Contains(s, ",")
}

func descSplitList(s string) []string {
	return strings.Split(s, ",")
}

func descJoinWithAnd(items []string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	case 2:
		return items[0] + " and " + items[1]
	default:
		return strings.Join(items[:len(items)-1], ", ") + ", and " + items[len(items)-1]
	}
}

// descAdjustDay shifts a day-of-week (0-6) by offset, wrapping around.
func descAdjustDay(day, offset int) int {
	day = (day + offset) % 7
	if day < 0 {
		day += 7
	}
	return day
}

// descFormatTimeWithTZ converts hour:minute from srcLoc to targetLoc in 12-hour format.
// Returns formatted time and day offset (-1, 0, or 1) if conversion crossed a day boundary.
func descFormatTimeWithTZ(hour, minute string, srcLoc, targetLoc *time.Location) (string, int) {
	h, err := strconv.Atoi(hour)
	if err != nil {
		return hour + ":" + minute, 0
	}
	m, _ := strconv.Atoi(minute)

	now := time.Now()
	srcTime := time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, srcLoc)
	targetTime := srcTime.In(targetLoc)

	dayOffset := targetTime.Day() - srcTime.Day()
	if dayOffset > 1 {
		dayOffset = -1
	} else if dayOffset < -1 {
		dayOffset = 1
	}

	targetH := targetTime.Hour()
	period := "AM"
	displayHour := targetH
	if targetH >= 12 {
		period = "PM"
		if targetH > 12 {
			displayHour = targetH - 12
		}
	}
	if displayHour == 0 {
		displayHour = 12
	}

	return fmt.Sprintf("%d:%02d %s", displayHour, m, period), dayOffset
}

// descFormatHourWithTZ converts hour from srcLoc to targetLoc in 12-hour format.
func descFormatHourWithTZ(hour string, srcLoc, targetLoc *time.Location) (string, int) {
	h, err := strconv.Atoi(hour)
	if err != nil {
		return hour + ":00", 0
	}

	now := time.Now()
	srcTime := time.Date(now.Year(), now.Month(), now.Day(), h, 0, 0, 0, srcLoc)
	targetTime := srcTime.In(targetLoc)

	dayOffset := targetTime.Day() - srcTime.Day()
	if dayOffset > 1 {
		dayOffset = -1
	} else if dayOffset < -1 {
		dayOffset = 1
	}

	targetH := targetTime.Hour()
	period := "AM"
	displayHour := targetH
	if targetH >= 12 {
		period = "PM"
		if targetH > 12 {
			displayHour = targetH - 12
		}
	}
	if displayHour == 0 {
		displayHour = 12
	}

	return fmt.Sprintf("%d:00 %s", displayHour, period), dayOffset
}
