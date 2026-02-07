package cronexpr

import (
	"slices"
	"time"
)

const daysPerWeek = 7

// dowNormalizedOffsets maps each weekday (Sun=0 .. Sat=6) to the possible
// day-of-month values for that weekday across the five weeks of a month.
var dowNormalizedOffsets = [][]int{
	{1, 8, 15, 22, 29},
	{2, 9, 16, 23, 30},
	{3, 10, 17, 24, 31},
	{4, 11, 18, 25},
	{5, 12, 19, 26},
	{6, 13, 20, 27},
	{7, 14, 21, 28},
}

// nextYear advances to the first matching instant in the next eligible year.
func (expr *Expression) nextYear(t time.Time) time.Time {
	i, _ := slices.BinarySearch(expr.yearList, t.Year()+1)
	if i == len(expr.yearList) {
		return time.Time{}
	}
	expr.actualDaysOfMonthList = expr.calculateActualDaysOfMonth(expr.yearList[i], expr.monthList[0])
	if len(expr.actualDaysOfMonthList) == 0 {
		return expr.nextMonth(time.Date(
			expr.yearList[i],
			time.Month(expr.monthList[0]),
			1,
			expr.hourList[0],
			expr.minuteList[0],
			expr.secondList[0],
			0,
			t.Location()))
	}
	return time.Date(
		expr.yearList[i],
		time.Month(expr.monthList[0]),
		expr.actualDaysOfMonthList[0],
		expr.hourList[0],
		expr.minuteList[0],
		expr.secondList[0],
		0,
		t.Location())
}

// nextMonth advances to the first matching instant in the next eligible month,
// cascading to nextYear if no remaining months match in the current year.
func (expr *Expression) nextMonth(t time.Time) time.Time {
	i, _ := slices.BinarySearch(expr.monthList, int(t.Month())+1)
	if i == len(expr.monthList) {
		return expr.nextYear(t)
	}
	expr.actualDaysOfMonthList = expr.calculateActualDaysOfMonth(t.Year(), expr.monthList[i])
	if len(expr.actualDaysOfMonthList) == 0 {
		return expr.nextMonth(time.Date(
			t.Year(),
			time.Month(expr.monthList[i]),
			1,
			expr.hourList[0],
			expr.minuteList[0],
			expr.secondList[0],
			0,
			t.Location()))
	}

	return time.Date(
		t.Year(),
		time.Month(expr.monthList[i]),
		expr.actualDaysOfMonthList[0],
		expr.hourList[0],
		expr.minuteList[0],
		expr.secondList[0],
		0,
		t.Location())
}

// nextDayOfMonth advances to the next eligible day within the current month,
// cascading to nextMonth if no remaining days match.
func (expr *Expression) nextDayOfMonth(t time.Time) time.Time {
	i, _ := slices.BinarySearch(expr.actualDaysOfMonthList, t.Day()+1)
	if i == len(expr.actualDaysOfMonthList) {
		return expr.nextMonth(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		expr.actualDaysOfMonthList[i],
		expr.hourList[0],
		expr.minuteList[0],
		expr.secondList[0],
		0,
		t.Location())
}

// nextHour advances to the next eligible hour within the current day,
// cascading to nextDayOfMonth if no remaining hours match.
func (expr *Expression) nextHour(t time.Time) time.Time {
	i, _ := slices.BinarySearch(expr.hourList, t.Hour()+1)
	if i == len(expr.hourList) {
		return expr.nextDayOfMonth(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		expr.hourList[i],
		expr.minuteList[0],
		expr.secondList[0],
		0,
		t.Location())
}

// nextMinute advances to the next eligible minute within the current hour,
// cascading to nextHour if no remaining minutes match.
func (expr *Expression) nextMinute(t time.Time) time.Time {
	i, _ := slices.BinarySearch(expr.minuteList, t.Minute()+1)
	if i == len(expr.minuteList) {
		return expr.nextHour(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		expr.minuteList[i],
		expr.secondList[0],
		0,
		t.Location())
}

// nextSecond assumes all other fields already match the cron expression.
func (expr *Expression) nextSecond(t time.Time) time.Time {
	i, _ := slices.BinarySearch(expr.secondList, t.Second()+1)
	if i == len(expr.secondList) {
		return expr.nextMinute(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		expr.secondList[i],
		0,
		t.Location())
}

// calculateActualDaysOfMonth computes the set of valid days for the given
// year/month by merging day-of-month and day-of-week constraints. Per the
// crontab spec, if both fields are restricted, a day matches if either matches.
func (expr *Expression) calculateActualDaysOfMonth(year, month int) []int {
	actualDaysOfMonthMap := make(map[int]bool)
	firstDayOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)

	// If both fields are unrestricted, all days of the month match.
	if !expr.daysOfMonthRestricted && !expr.daysOfWeekRestricted {
		return genericDefaultList[1 : lastDayOfMonth.Day()+1]
	}

	if expr.daysOfMonthRestricted {
		if expr.lastDayOfMonth {
			actualDaysOfMonthMap[lastDayOfMonth.Day()] = true
		}
		if expr.lastWorkdayOfMonth {
			actualDaysOfMonthMap[workdayOfMonth(lastDayOfMonth, lastDayOfMonth)] = true
		}
		for v := range expr.daysOfMonth {
			if v <= lastDayOfMonth.Day() {
				actualDaysOfMonthMap[v] = true
			}
		}
		// W (nearest weekday) does not cross month boundaries.
		for v := range expr.workdaysOfMonth {
			if v <= lastDayOfMonth.Day() {
				actualDaysOfMonthMap[workdayOfMonth(firstDayOfMonth.AddDate(0, 0, v-1), lastDayOfMonth)] = true
			}
		}
	}

	if expr.daysOfWeekRestricted {
		offset := daysPerWeek - int(firstDayOfMonth.Weekday())
		for v := range expr.daysOfWeek {
			w := dowNormalizedOffsets[(offset+v)%daysPerWeek]
			for _, day := range w {
				if day <= lastDayOfMonth.Day() {
					actualDaysOfMonthMap[day] = true
				}
			}
		}
		for v := range expr.specificWeekDaysOfWeek {
			v = 1 + daysPerWeek*(v/daysPerWeek) + (offset+v)%daysPerWeek
			if v <= lastDayOfMonth.Day() {
				actualDaysOfMonthMap[v] = true
			}
		}
		lastWeekOrigin := firstDayOfMonth.AddDate(0, 1, -daysPerWeek)
		offset = daysPerWeek - int(lastWeekOrigin.Weekday())
		for v := range expr.lastWeekDaysOfWeek {
			v = lastWeekOrigin.Day() + (offset+v)%daysPerWeek
			if v <= lastDayOfMonth.Day() {
				actualDaysOfMonthMap[v] = true
			}
		}
	}

	return toList(actualDaysOfMonthMap)
}

// workdayOfMonth returns the nearest weekday to targetDom that does not cross
// the month boundary defined by lastDom.
func workdayOfMonth(targetDom, lastDom time.Time) int {
	dom := targetDom.Day()
	switch dow := targetDom.Weekday(); dow {
	case time.Saturday:
		if dom > 1 {
			dom--
		} else {
			dom += 2
		}
	case time.Sunday:
		if dom < lastDom.Day() {
			dom++
		} else {
			dom -= 2
		}
	}
	return dom
}
