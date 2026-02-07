package cronexpr

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
)

const (
	minYear = 1970
	maxYear = 2099
)

// makeIntRange returns a slice of consecutive integers from lo to hi inclusive.
func makeIntRange(lo, hi int) []int {
	s := make([]int, hi-lo+1)
	for i := range s {
		s[i] = lo + i
	}
	return s
}

var (
	genericDefaultList = makeIntRange(0, 59)
	yearDefaultList    = makeIntRange(minYear, maxYear)
)

var (
	numberTokens = func() map[string]int {
		m := make(map[string]int, 60+10+maxYear-minYear+1) // bare, zero-padded, years
		for i := range 60 {
			m[strconv.Itoa(i)] = i
			if i < 10 {
				m[fmt.Sprintf("%02d", i)] = i
			}
		}
		for y := minYear; y <= maxYear; y++ {
			m[strconv.Itoa(y)] = y
		}
		return m
	}()
	monthTokens = map[string]int{
		`1`: 1, `jan`: 1, `january`: 1,
		`2`: 2, `feb`: 2, `february`: 2,
		`3`: 3, `mar`: 3, `march`: 3,
		`4`: 4, `apr`: 4, `april`: 4,
		`5`: 5, `may`: 5,
		`6`: 6, `jun`: 6, `june`: 6,
		`7`: 7, `jul`: 7, `july`: 7,
		`8`: 8, `aug`: 8, `august`: 8,
		`9`: 9, `sep`: 9, `september`: 9,
		`10`: 10, `oct`: 10, `october`: 10,
		`11`: 11, `nov`: 11, `november`: 11,
		`12`: 12, `dec`: 12, `december`: 12,
	}
	dowTokens = map[string]int{
		`0`: 0, `sun`: 0, `sunday`: 0,
		`1`: 1, `mon`: 1, `monday`: 1,
		`2`: 2, `tue`: 2, `tuesday`: 2,
		`3`: 3, `wed`: 3, `wednesday`: 3,
		`4`: 4, `thu`: 4, `thursday`: 4,
		`5`: 5, `fri`: 5, `friday`: 5,
		`6`: 6, `sat`: 6, `saturday`: 6,
		`7`: 0,
	}
)

// atoi converts a numeric string to int using the pre-built numberTokens lookup table.
func atoi(s string) int {
	return numberTokens[s]
}

type fieldDescriptor struct {
	name         string
	min, max     int
	defaultList  []int
	valuePattern string
	atoi         func(string) int
}

var (
	secondDescriptor = fieldDescriptor{
		name:         "second",
		min:          0,
		max:          59,
		defaultList:  genericDefaultList[0:60],
		valuePattern: `0?[0-9]|[1-5][0-9]`,
		atoi:         atoi,
	}
	minuteDescriptor = fieldDescriptor{
		name:         "minute",
		min:          0,
		max:          59,
		defaultList:  genericDefaultList[0:60],
		valuePattern: `0?[0-9]|[1-5][0-9]`,
		atoi:         atoi,
	}
	hourDescriptor = fieldDescriptor{
		name:         "hour",
		min:          0,
		max:          23,
		defaultList:  genericDefaultList[0:24],
		valuePattern: `0?[0-9]|1[0-9]|2[0-3]`,
		atoi:         atoi,
	}
	domDescriptor = fieldDescriptor{
		name:         "day-of-month",
		min:          1,
		max:          31,
		defaultList:  genericDefaultList[1:32],
		valuePattern: `0?[1-9]|[12][0-9]|3[01]`,
		atoi:         atoi,
	}
	monthDescriptor = fieldDescriptor{
		name:         "month",
		min:          1,
		max:          12,
		defaultList:  genericDefaultList[1:13],
		valuePattern: `0?[1-9]|1[012]|jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|january|february|march|april|may|june|july|august|september|october|november|december`,
		atoi: func(s string) int {
			return monthTokens[s]
		},
	}
	dowDescriptor = fieldDescriptor{
		name:         "day-of-week",
		min:          0,
		max:          6,
		defaultList:  genericDefaultList[0:7],
		valuePattern: `0?[0-7]|sun|mon|tue|wed|thu|fri|sat|sunday|monday|tuesday|wednesday|thursday|friday|saturday`,
		atoi: func(s string) int {
			return dowTokens[s]
		},
	}
	yearDescriptor = fieldDescriptor{
		name:         "year",
		min:          minYear,
		max:          maxYear,
		defaultList:  yearDefaultList,
		valuePattern: `19[789][0-9]|20[0-9]{2}`,
		atoi:         atoi,
	}
)

// Layout patterns for matching cron field entries. %value% is replaced with
// the field-specific value pattern before compilation.
var (
	layoutWildcard            = `^\*$|^\?$`
	layoutValue               = `^(%value%)$`
	layoutRange               = `^(%value%)-(%value%)$`
	layoutWildcardAndInterval = `^\*/(\d+)$`
	layoutValueAndInterval    = `^(%value%)/(\d+)$`
	layoutRangeAndInterval    = `^(%value%)-(%value%)/(\d+)$`
	layoutLastDom             = `^l$`
	layoutWorkdom             = `^(%value%)w$`
	layoutLastWorkdom         = `^lw$`
	layoutDowOfLastWeek       = `^(%value%)l$`
	layoutDowOfSpecificWeek   = `^(%value%)#([1-5])$`
	fieldFinder               = regexp.MustCompile(`\S+`)
	entryFinder               = regexp.MustCompile(`[^,]+`)
	layoutRegexp              sync.Map
)

// cronNormalizer expands predefined cron aliases into 7-field expressions.
var cronNormalizer = strings.NewReplacer(
	"@yearly", "0 0 0 1 1 * *",
	"@annually", "0 0 0 1 1 * *",
	"@monthly", "0 0 0 1 * * *",
	"@weekly", "0 0 0 * * 0 *",
	"@daily", "0 0 0 * * * *",
	"@hourly", "0 0 * * * * *")

// parseField parses a single cron field string into a sorted list of matching
// integer values using the given field descriptor.
func parseField(s string, desc fieldDescriptor, target *[]int) error {
	var err error
	*target, err = genericFieldHandler(s, desc)
	return err
}

// Directive kinds returned by genericFieldParse.
const (
	none = 0
	one  = 1
	span = 2
	all  = 3
)

type cronDirective struct {
	kind  int
	first int
	last  int
	step  int
	sbeg  int
	send  int
}

// genericFieldHandler converts parsed directives into a sorted list of matching
// values for a standard cron field (one without special modifiers like L or W).
func genericFieldHandler(s string, desc fieldDescriptor) ([]int, error) {
	directives, err := genericFieldParse(s, desc)
	if err != nil {
		return nil, err
	}
	values := make(map[int]bool)
	for _, directive := range directives {
		switch directive.kind {
		case none:
			return nil, fmt.Errorf("syntax error in %s field: '%s'", desc.name, s[directive.sbeg:directive.send])
		case one:
			populateOne(values, directive.first)
		case span:
			populateMany(values, directive.first, directive.last, directive.step, desc.min, desc.max)
		case all:
			return desc.defaultList, nil
		}
	}
	return toList(values), nil
}

// dowFieldHandler parses the day-of-week field, handling standard values plus
// special modifiers like L (last week of month) and # (specific week number).
func (expr *Expression) dowFieldHandler(s string) error {
	expr.daysOfWeekRestricted = true
	expr.daysOfWeek = make(map[int]bool)
	expr.lastWeekDaysOfWeek = make(map[int]bool)
	expr.specificWeekDaysOfWeek = make(map[int]bool)

	directives, err := genericFieldParse(s, dowDescriptor)
	if err != nil {
		return err
	}

	for _, directive := range directives {
		switch directive.kind {
		case none:
			sdirective := s[directive.sbeg:directive.send]
			snormal := strings.ToLower(sdirective)
			// `5L`
			pairs := makeLayoutRegexp(layoutDowOfLastWeek, dowDescriptor.valuePattern).FindStringSubmatchIndex(snormal)
			if len(pairs) > 0 {
				populateOne(expr.lastWeekDaysOfWeek, dowDescriptor.atoi(snormal[pairs[2]:pairs[3]]))
			} else {
				// `5#3`
				pairs := makeLayoutRegexp(layoutDowOfSpecificWeek, dowDescriptor.valuePattern).FindStringSubmatchIndex(snormal)
				if len(pairs) > 0 {
					populateOne(expr.specificWeekDaysOfWeek, (dowDescriptor.atoi(snormal[pairs[4]:pairs[5]])-1)*7+(dowDescriptor.atoi(snormal[pairs[2]:pairs[3]])%7))
				} else {
					return fmt.Errorf("syntax error in day-of-week field: '%s'", sdirective)
				}
			}
		case one:
			populateOne(expr.daysOfWeek, directive.first)
		case span:
			populateMany(expr.daysOfWeek, directive.first, directive.last, directive.step, dowDescriptor.min, dowDescriptor.max)
		case all:
			populateMany(expr.daysOfWeek, directive.first, directive.last, directive.step, dowDescriptor.min, dowDescriptor.max)
			expr.daysOfWeekRestricted = false
		}
	}
	return nil
}

// domFieldHandler parses the day-of-month field, handling standard values plus
// special modifiers like L (last day), W (nearest weekday), and LW (last weekday).
func (expr *Expression) domFieldHandler(s string) error {
	expr.daysOfMonthRestricted = true
	expr.lastDayOfMonth = false
	expr.lastWorkdayOfMonth = false
	expr.daysOfMonth = make(map[int]bool)
	expr.workdaysOfMonth = make(map[int]bool)

	directives, err := genericFieldParse(s, domDescriptor)
	if err != nil {
		return err
	}

	for _, directive := range directives {
		switch directive.kind {
		case none:
			sdirective := s[directive.sbeg:directive.send]
			snormal := strings.ToLower(sdirective)
			// `L`
			if makeLayoutRegexp(layoutLastDom, domDescriptor.valuePattern).MatchString(snormal) {
				expr.lastDayOfMonth = true
			} else {
				// `LW`
				if makeLayoutRegexp(layoutLastWorkdom, domDescriptor.valuePattern).MatchString(snormal) {
					expr.lastWorkdayOfMonth = true
				} else {
					// `15W`
					pairs := makeLayoutRegexp(layoutWorkdom, domDescriptor.valuePattern).FindStringSubmatchIndex(snormal)
					if len(pairs) > 0 {
						populateOne(expr.workdaysOfMonth, domDescriptor.atoi(snormal[pairs[2]:pairs[3]]))
					} else {
						return fmt.Errorf("syntax error in day-of-month field: '%s'", sdirective)
					}
				}
			}
		case one:
			populateOne(expr.daysOfMonth, directive.first)
		case span:
			populateMany(expr.daysOfMonth, directive.first, directive.last, directive.step, domDescriptor.min, domDescriptor.max)
		case all:
			populateMany(expr.daysOfMonth, directive.first, directive.last, directive.step, domDescriptor.min, domDescriptor.max)
			expr.daysOfMonthRestricted = false
		}
	}
	return nil
}

// populateOne adds a single value to the set.
func populateOne(values map[int]bool, v int) {
	values[v] = true
}

// populateMany fills values for a range [lo..hi] with the given step.
// If lo > hi, the range wraps around through fieldMax back to fieldMin
// (e.g. hours 22-3 becomes 22,23,0,1,2,3).
func populateMany(values map[int]bool, lo, hi, step int, bounds ...int) {
	if lo <= hi {
		for i := lo; i <= hi; i += step {
			values[i] = true
		}
		return
	}
	// Wrap-around: lo > hi requires field bounds.
	if len(bounds) < 2 {
		return
	}
	fieldMin, fieldMax := bounds[0], bounds[1]
	for i := lo; i <= fieldMax; i += step {
		values[i] = true
	}
	for i := fieldMin; i <= hi; i += step {
		values[i] = true
	}
}

// toList converts a set of integers into a sorted slice.
func toList(set map[int]bool) []int {
	return slices.Sorted(maps.Keys(set))
}

// validateStep checks that a step/interval value is between 1 and the field's max.
func validateStep(step, maxVal int, raw string) error {
	if step < 1 || step > maxVal {
		return fmt.Errorf("invalid interval %s", raw)
	}
	return nil
}

// genericFieldParse tokenizes a cron field string into directives by matching
// each comma-separated entry against the layout patterns.
func genericFieldParse(s string, desc fieldDescriptor) ([]*cronDirective, error) {
	indices := entryFinder.FindAllStringIndex(s, -1)
	if len(indices) == 0 {
		return nil, fmt.Errorf("%s field: missing directive", desc.name)
	}

	directives := make([]*cronDirective, 0, len(indices))

	for i := range indices {
		directive := cronDirective{
			sbeg: indices[i][0],
			send: indices[i][1],
		}
		snormal := strings.ToLower(s[indices[i][0]:indices[i][1]])

		// `*`
		if makeLayoutRegexp(layoutWildcard, desc.valuePattern).MatchString(snormal) {
			directive.kind = all
			directive.first = desc.min
			directive.last = desc.max
			directive.step = 1
			directives = append(directives, &directive)
			continue
		}
		// `5`
		if makeLayoutRegexp(layoutValue, desc.valuePattern).MatchString(snormal) {
			directive.kind = one
			directive.first = desc.atoi(snormal)
			directives = append(directives, &directive)
			continue
		}
		// `5-20`
		pairs := makeLayoutRegexp(layoutRange, desc.valuePattern).FindStringSubmatchIndex(snormal)
		if len(pairs) > 0 {
			directive.kind = span
			directive.first = desc.atoi(snormal[pairs[2]:pairs[3]])
			directive.last = desc.atoi(snormal[pairs[4]:pairs[5]])
			directive.step = 1
			directives = append(directives, &directive)
			continue
		}
		// `*/2`
		pairs = makeLayoutRegexp(layoutWildcardAndInterval, desc.valuePattern).FindStringSubmatchIndex(snormal)
		if len(pairs) > 0 {
			directive.kind = span
			directive.first = desc.min
			directive.last = desc.max
			directive.step = atoi(snormal[pairs[2]:pairs[3]])
			if err := validateStep(directive.step, desc.max, snormal); err != nil {
				return nil, err
			}
			directives = append(directives, &directive)
			continue
		}
		// `5/2`
		pairs = makeLayoutRegexp(layoutValueAndInterval, desc.valuePattern).FindStringSubmatchIndex(snormal)
		if len(pairs) > 0 {
			directive.kind = span
			directive.first = desc.atoi(snormal[pairs[2]:pairs[3]])
			directive.last = desc.max
			directive.step = atoi(snormal[pairs[4]:pairs[5]])
			if err := validateStep(directive.step, desc.max, snormal); err != nil {
				return nil, err
			}
			directives = append(directives, &directive)
			continue
		}
		// `5-20/2`
		pairs = makeLayoutRegexp(layoutRangeAndInterval, desc.valuePattern).FindStringSubmatchIndex(snormal)
		if len(pairs) > 0 {
			directive.kind = span
			directive.first = desc.atoi(snormal[pairs[2]:pairs[3]])
			directive.last = desc.atoi(snormal[pairs[4]:pairs[5]])
			directive.step = atoi(snormal[pairs[6]:pairs[7]])
			if err := validateStep(directive.step, desc.max, snormal); err != nil {
				return nil, err
			}
			directives = append(directives, &directive)
			continue
		}
		directive.kind = none
		directives = append(directives, &directive)
	}
	return directives, nil
}

// makeLayoutRegexp compiles a layout pattern with the field's value pattern
// substituted in, caching the result for reuse.
func makeLayoutRegexp(layout, value string) *regexp.Regexp {
	layout = strings.ReplaceAll(layout, `%value%`, value)
	if re, ok := layoutRegexp.Load(layout); ok {
		return re.(*regexp.Regexp)
	}
	re := regexp.MustCompile(layout)
	layoutRegexp.Store(layout, re)
	return re
}
