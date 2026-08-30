package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/mrusme/symmetry454"
)

const defaultFormat = "{{.Gregorian}} -> {{.Sym}} " +
	"{{.Weekday}}, {{.MonthName}} {{.Day}}, {{.Year}}"

const defaultSymFormat = "{{.Sym}} -> {{.Gregorian}} {{.Weekday}}"

var epochSecondsPattern = regexp.MustCompile(`^[+-]?[0-9]+$`)

var symDatePattern = regexp.MustCompile(`^(-?[0-9]+)-([0-9]{1,2})-([0-9]{1,2})$`)

var instantLayouts = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC3339Nano,
	time.RFC3339,
	time.UnixDate,
	time.RubyDate,
	time.ANSIC,
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"02 Jan 2006 15:04:05 -0700",
	"2 Jan 2006 15:04:05 -0700",
	"2006-01-02",
	"2006/01/02",
}

type convertOptions struct {
	location *time.Location
	format   string
	long     bool
	asJSON   bool
	fromSym  bool
}

type conversion struct {
	Gregorian   string `json:"gregorian"`
	Unix        int64  `json:"unix"`
	Fixed       int    `json:"fixed"`
	Sym         string `json:"sym454"`
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	MonthName   string `json:"month_name"`
	Day         int    `json:"day"`
	Weekday     string `json:"weekday"`
	DayOfYear   int    `json:"day_of_year"`
	WeekOfYear  int    `json:"week_of_year"`
	Quarter     int    `json:"quarter"`
	DaysInMonth int    `json:"days_in_month"`
	DaysInYear  int    `json:"days_in_year"`
	WeeksInYear int    `json:"weeks_in_year"`
	LeapYear    bool   `json:"leap_year"`
}

func parseInstant(input string, loc *time.Location) (time.Time, error) {
	text := strings.TrimSpace(input)
	if text == "" {
		return time.Time{}, fmt.Errorf("empty input")
	}

	if epochSecondsPattern.MatchString(text) {
		seconds, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("epoch seconds %q: %w", text, err)
		}
		return time.Unix(seconds, 0).In(loc), nil
	}

	for _, layout := range instantLayouts {
		if parsed, err := time.ParseInLocation(layout, text, loc); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized date format %q", text)
}

func parseSym(input string) (symmetry454.Sym, error) {
	text := strings.TrimSpace(input)

	fields := symDatePattern.FindStringSubmatch(text)
	if fields == nil {
		return symmetry454.Sym{},
			fmt.Errorf("unrecognized Symmetry454 date %q, want YYYY-MM-DD", text)
	}

	year, err := strconv.Atoi(fields[1])
	if err != nil {
		return symmetry454.Sym{}, fmt.Errorf("year in %q: %w", text, err)
	}

	month, err := strconv.Atoi(fields[2])
	if err != nil {
		return symmetry454.Sym{}, fmt.Errorf("month in %q: %w", text, err)
	}

	day, err := strconv.Atoi(fields[3])
	if err != nil {
		return symmetry454.Sym{}, fmt.Errorf("day in %q: %w", text, err)
	}

	if month < 1 || month > 12 {
		return symmetry454.Sym{}, fmt.Errorf(
			"invalid Symmetry454 date %q, month %d is not in 1 to 12", text, month)
	}

	sym := symmetry454.Sym{Year: year, Month: month, Day: day}
	if !sym.Valid() {
		return symmetry454.Sym{}, fmt.Errorf(
			"invalid Symmetry454 date %q, %s %d has %d days",
			text, symmetry454.MonthName(month), year,
			symmetry454.SymDaysInMonth(year, month))
	}

	return sym, nil
}

func convertTime(t time.Time) conversion {
	return newConversion(symmetry454.FromTime(t), t.Format(time.RFC3339), t.Unix())
}

func convertSym(sym symmetry454.Sym, loc *time.Location) conversion {
	t := sym.ToTimeIn(loc)
	return newConversion(sym, t.Format("2006-01-02"), t.Unix())
}

func newConversion(
	sym symmetry454.Sym,
	gregorian string,
	unix int64,
) conversion {
	return conversion{
		Gregorian:   gregorian,
		Unix:        unix,
		Fixed:       sym.Fixed(),
		Sym:         sym.String(),
		Year:        sym.Year,
		Month:       sym.Month,
		MonthName:   sym.MonthName(),
		Day:         sym.Day,
		Weekday:     sym.Weekday().String(),
		DayOfYear:   sym.DayOfYear(),
		WeekOfYear:  sym.WeekOfYear(),
		Quarter:     sym.Quarter(),
		DaysInMonth: sym.DaysInMonth(),
		DaysInYear:  sym.DaysInYear(),
		WeeksInYear: sym.WeeksInYear(),
		LeapYear:    sym.IsLeapYear(),
	}
}

func longOutput(result conversion) string {
	leapYear := "no"
	if result.LeapYear {
		leapYear = "yes"
	}

	var out strings.Builder
	fmt.Fprintf(&out, "Gregorian     %s\n", result.Gregorian)
	fmt.Fprintf(&out, "Unix          %d\n", result.Unix)
	fmt.Fprintf(&out, "Fixed day     %d\n", result.Fixed)
	fmt.Fprintf(&out, "Symmetry454   %s\n", result.Sym)
	fmt.Fprintf(&out, "Long form     %s, %s %d, %d\n",
		result.Weekday, result.MonthName, result.Day, result.Year)
	fmt.Fprintf(&out, "Day of year   %d of %d\n", result.DayOfYear, result.DaysInYear)
	fmt.Fprintf(&out, "Week of year  %d of %d\n", result.WeekOfYear, result.WeeksInYear)
	fmt.Fprintf(&out, "Quarter       %d\n", result.Quarter)
	fmt.Fprintf(&out, "Days in month %d\n", result.DaysInMonth)
	fmt.Fprintf(&out, "Leap year     %s\n", leapYear)

	return out.String()
}

func runConvert(
	in io.Reader,
	out io.Writer,
	errOut io.Writer,
	opts convertOptions,
) int {
	format := opts.format
	if format == "" {
		format = defaultFormat
		if opts.fromSym {
			format = defaultSymFormat
		}
	}

	tmpl, err := template.New("line").Parse(format)
	if err != nil {
		fmt.Fprintf(errOut, "sym454: invalid format: %v\n", err)
		return 2
	}

	encoder := json.NewEncoder(out)
	status := 0
	blocks := 0

	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		result, err := convertLine(line, opts)
		if err != nil {
			fmt.Fprintf(errOut, "sym454: %v\n", err)
			status = 1
			continue
		}

		switch {
		case opts.asJSON:
			if err := encoder.Encode(result); err != nil {
				fmt.Fprintf(errOut, "sym454: %v\n", err)
				return 1
			}
		case opts.long:
			if blocks > 0 {
				fmt.Fprintln(out)
			}
			fmt.Fprint(out, longOutput(result))
			blocks++
		default:
			if err := tmpl.Execute(out, result); err != nil {
				fmt.Fprintf(errOut, "sym454: %v\n", err)
				return 1
			}
			fmt.Fprintln(out)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(errOut, "sym454: reading input: %v\n", err)
		return 1
	}

	return status
}

func convertLine(line string, opts convertOptions) (conversion, error) {
	if opts.fromSym {
		sym, err := parseSym(line)
		if err != nil {
			return conversion{}, err
		}
		return convertSym(sym, opts.location), nil
	}

	instant, err := parseInstant(line, opts.location)
	if err != nil {
		return conversion{}, err
	}

	return convertTime(instant), nil
}
