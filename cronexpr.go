// Package cronexpr parses cron time expressions.
package cronexpr

import (
	"errors"
	"slices"
	"time"
)

// Expression represents a parsed cron expression. Use Parse or MustParse to create one.
type Expression struct {
	secondList             []int
	minuteList             []int
	hourList               []int
	daysOfMonth            map[int]bool
	workdaysOfMonth        map[int]bool
	lastDayOfMonth         bool
	lastWorkdayOfMonth     bool
	daysOfMonthRestricted  bool
	actualDaysOfMonthList  []int
	monthList              []int
	daysOfWeek             map[int]bool
	specificWeekDaysOfWeek map[int]bool
	lastWeekDaysOfWeek     map[int]bool
	daysOfWeekRestricted   bool
	yearList               []int
}

// MustParse returns a new Expression pointer. It expects a well-formed cron
// expression. If a malformed cron expression is supplied, it will panic.
func MustParse(cronLine string) *Expression {
	expr, err := Parse(cronLine)
	if err != nil {
		panic(err)
	}
	return expr
}

// Parse returns a new Expression pointer. An error is returned if a malformed
// cron expression is supplied.
func Parse(cronLine string) (*Expression, error) {

	// Maybe one of the built-in aliases is being used
	cron := cronNormalizer.Replace(cronLine)

	const (
		minCronFields = 5
		maxCronFields = 7
	)

	indices := fieldFinder.FindAllStringIndex(cron, -1)
	fieldCount := len(indices)
	if fieldCount < minCronFields {
		return nil, errors.New("missing field(s)")
	}
	// ignore fields beyond 7th
	if fieldCount > maxCronFields {
		fieldCount = maxCronFields
	}

	var expr Expression
	var field = 0
	var err error

	// second field (optional)
	if fieldCount == maxCronFields {
		err = parseField(cron[indices[field][0]:indices[field][1]], secondDescriptor, &expr.secondList)
		if err != nil {
			return nil, err
		}
		field++
	} else {
		expr.secondList = []int{0}
	}

	// minute field
	err = parseField(cron[indices[field][0]:indices[field][1]], minuteDescriptor, &expr.minuteList)
	if err != nil {
		return nil, err
	}
	field++

	// hour field
	err = parseField(cron[indices[field][0]:indices[field][1]], hourDescriptor, &expr.hourList)
	if err != nil {
		return nil, err
	}
	field++

	// day of month field
	err = expr.domFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field++

	// month field
	err = parseField(cron[indices[field][0]:indices[field][1]], monthDescriptor, &expr.monthList)
	if err != nil {
		return nil, err
	}
	field++

	// day of week field
	err = expr.dowFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field++

	// year field
	if field < fieldCount {
		err = parseField(cron[indices[field][0]:indices[field][1]], yearDescriptor, &expr.yearList)
		if err != nil {
			return nil, err
		}
	} else {
		expr.yearList = yearDescriptor.defaultList
	}

	return &expr, nil
}

// Next returns the closest time instant immediately following fromTime which
// matches the cron expression.
//
// The time.Location of the returned time instant is the same as that of
// fromTime.
//
// The zero value of time.Time is returned if no matching time instant exists
// or if fromTime is itself a zero value.
func (expr *Expression) Next(fromTime time.Time) time.Time {
	// Special case
	if fromTime.IsZero() {
		return fromTime
	}

	// Walk each field from year down to second. If any field doesn't match,
	// advance to the next matching time for that field.
	// year
	v := fromTime.Year()
	i, _ := slices.BinarySearch(expr.yearList, v)
	if i == len(expr.yearList) {
		return time.Time{}
	}
	if v != expr.yearList[i] {
		return expr.nextYear(fromTime)
	}
	// month
	v = int(fromTime.Month())
	i, _ = slices.BinarySearch(expr.monthList, v)
	if i == len(expr.monthList) {
		return expr.nextYear(fromTime)
	}
	if v != expr.monthList[i] {
		return expr.nextMonth(fromTime)
	}

	expr.actualDaysOfMonthList = expr.calculateActualDaysOfMonth(fromTime.Year(), int(fromTime.Month()))
	if len(expr.actualDaysOfMonthList) == 0 {
		return expr.nextMonth(fromTime)
	}

	// day of month
	v = fromTime.Day()
	i, _ = slices.BinarySearch(expr.actualDaysOfMonthList, v)
	if i == len(expr.actualDaysOfMonthList) {
		return expr.nextMonth(fromTime)
	}
	if v != expr.actualDaysOfMonthList[i] {
		return expr.nextDayOfMonth(fromTime)
	}
	// hour
	v = fromTime.Hour()
	i, _ = slices.BinarySearch(expr.hourList, v)
	if i == len(expr.hourList) {
		return expr.nextDayOfMonth(fromTime)
	}
	if v != expr.hourList[i] {
		return expr.nextHour(fromTime)
	}
	// minute
	v = fromTime.Minute()
	i, _ = slices.BinarySearch(expr.minuteList, v)
	if i == len(expr.minuteList) {
		return expr.nextHour(fromTime)
	}
	if v != expr.minuteList[i] {
		return expr.nextMinute(fromTime)
	}
	// second
	v = fromTime.Second()
	i, _ = slices.BinarySearch(expr.secondList, v)
	if i == len(expr.secondList) {
		return expr.nextMinute(fromTime)
	}

	return expr.nextSecond(fromTime)
}

// NextN returns a slice of the n closest time instants immediately following
// fromTime which match the cron expression.
//
// The time instants in the returned slice are in chronological ascending order.
// The time.Location of the returned time instants is the same as that of
// fromTime.
//
// A slice with length between 0 and n is returned; if not enough matching
// time instants exist, the number of returned entries will be less than n.
func (expr *Expression) NextN(fromTime time.Time, n uint) []time.Time {
	nextTimes := make([]time.Time, 0, n)
	if n > 0 {
		fromTime = expr.Next(fromTime)
		for !fromTime.IsZero() {
			nextTimes = append(nextTimes, fromTime)
			n--
			if n == 0 {
				break
			}
			fromTime = expr.nextSecond(fromTime)
		}
	}
	return nextTimes
}
