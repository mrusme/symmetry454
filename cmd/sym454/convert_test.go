package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mrusme/symmetry454"
)

type tableRow struct {
	gregorian string
	fixedDate int
	sym       string
	weekday   string
}

var verificationTable = []tableRow{
	{"0122-09-07", 44444, "0122-09-08", "Monday"},
	{"1776-07-04", 648491, "1776-07-04", "Thursday"},
	{"1867-07-01", 681724, "1867-07-01", "Monday"},
	{"1947-10-24", 711058, "1947-10-26", "Friday"},
	{"1995-08-10", 728515, "1995-08-11", "Thursday"},
	{"2000-02-29", 730179, "2000-02-30", "Tuesday"},
	{"2004-05-02", 731703, "2004-05-07", "Sunday"},
	{"2004-12-31", 731946, "2004-12-33", "Friday"},
	{"2020-02-20", 737475, "2020-02-25", "Thursday"},
	{"2222-02-02", 811236, "2222-02-06", "Saturday"},
	{"3333-03-01", 1217048, "3333-02-35", "Sunday"},
}

func TestParseInstantEpochSeconds(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"0", 0},
		{"1788094843", 1788094843},
		{"-1", -1},
		{"+42", 42},
		{"  1788094843  ", 1788094843},
	}

	for _, test := range tests {
		parsed, err := parseInstant(test.input, time.UTC)
		if err != nil {
			t.Fatalf(`parseInstant(%q) returned %v`, test.input, err)
		}
		if parsed.Unix() != test.want {
			t.Fatalf(`parseInstant(%q) = %v, want unix %v`,
				test.input, parsed.Unix(), test.want)
		}
	}
}

func TestParseInstantLayouts(t *testing.T) {
	want := time.Date(2026, 8, 30, 23, 30, 0, 0, time.FixedZone("", 2*3600))

	inputs := []string{
		"Sun, 30 Aug 2026 23:30:00 +0200",
		"2026-08-30T23:30:00+02:00",
		"2026-08-30 23:30:00 +0200",
		"30 Aug 2026 23:30:00 +0200",
	}

	for _, input := range inputs {
		parsed, err := parseInstant(input, time.UTC)
		if err != nil {
			t.Fatalf(`parseInstant(%q) returned %v`, input, err)
		}
		if !parsed.Equal(want) {
			t.Fatalf(`parseInstant(%q) = %v, want %v`,
				input, parsed.Format(time.RFC3339), want.Format(time.RFC3339))
		}
	}
}

func TestParseInstantDateOnly(t *testing.T) {
	parsed, err := parseInstant("2026-08-30", time.UTC)
	if err != nil {
		t.Fatalf(`parseInstant returned %v`, err)
	}

	if parsed.Year() != 2026 || parsed.Month() != time.August || parsed.Day() != 30 {
		t.Fatalf(`parseInstant = %v, want 2026-08-30`, parsed.Format(time.RFC3339))
	}
	if parsed.Location() != time.UTC {
		t.Fatalf(`parseInstant location = %v, want UTC`, parsed.Location())
	}
}

func TestParseInstantUnixDate(t *testing.T) {
	parsed, err := parseInstant("Sun Aug 30 23:30:00 UTC 2026", time.UTC)
	if err != nil {
		t.Fatalf(`parseInstant returned %v`, err)
	}

	if parsed.Year() != 2026 || parsed.Month() != time.August || parsed.Day() != 30 {
		t.Fatalf(`parseInstant = %v, want 2026-08-30`, parsed.Format(time.RFC3339))
	}
}

func TestParseInstantErrors(t *testing.T) {
	for _, input := range []string{"", "   ", "not-a-date", "2026-13-45", "tomorrow"} {
		if _, err := parseInstant(input, time.UTC); err == nil {
			t.Fatalf(`parseInstant(%q) returned no error`, input)
		}
	}
}

func TestParseSym(t *testing.T) {
	tests := []struct {
		input string
		want  symmetry454.Sym
	}{
		{"2026-08-35", symmetry454.Sym{Year: 2026, Month: 8, Day: 35}},
		{"2009-12-33", symmetry454.Sym{Year: 2009, Month: 12, Day: 33}},
		{"-0121-04-27", symmetry454.Sym{Year: -121, Month: 4, Day: 27}},
		{"-121-4-27", symmetry454.Sym{Year: -121, Month: 4, Day: 27}},
		{"  2026-01-01  ", symmetry454.Sym{Year: 2026, Month: 1, Day: 1}},
	}

	for _, test := range tests {
		parsed, err := parseSym(test.input)
		if err != nil {
			t.Fatalf(`parseSym(%q) returned %v`, test.input, err)
		}
		if parsed != test.want {
			t.Fatalf(`parseSym(%q) = %v, want %v`, test.input, parsed, test.want)
		}
	}
}

func TestParseSymErrors(t *testing.T) {
	tests := []string{
		"2026-01-29",
		"2026-13-01",
		"2026-00-01",
		"2026-01-00",
		"2027-12-33",
		"2026-08-36",
		"not-a-date",
		"2026-08",
		"",
	}

	for _, input := range tests {
		if _, err := parseSym(input); err == nil {
			t.Fatalf(`parseSym(%q) returned no error`, input)
		}
	}
}

func TestParseSymRejectsWhatRollsOver(t *testing.T) {
	for symYear := 2020; symYear <= 2030; symYear++ {
		for symMonth := 1; symMonth <= 12; symMonth++ {
			for symDay := 1; symDay <= 38; symDay++ {
				input := symmetry454.Sym{
					Year:  symYear,
					Month: symMonth,
					Day:   symDay,
				}
				parsed, err := parseSym(input.String())

				if err == nil && parsed != input {
					t.Fatalf(`parseSym(%q) = %v, which is a different date`,
						input.String(), parsed)
				}
				if err != nil && input.Valid() {
					t.Fatalf(`parseSym(%q) rejected a valid date: %v`,
						input.String(), err)
				}
			}
		}
	}
}

func runLines(t *testing.T, input string, opts convertOptions) (string, string, int) {
	t.Helper()

	var out, errOut strings.Builder
	status := runConvert(strings.NewReader(input), &out, &errOut, opts)

	return out.String(), errOut.String(), status
}

func TestRunConvertVerificationTable(t *testing.T) {
	input := make([]string, 0, len(verificationTable))
	for _, row := range verificationTable {
		input = append(input, row.gregorian)
	}

	out, errOut, status := runLines(t, strings.Join(input, "\n"), convertOptions{
		location: time.UTC,
		format:   "{{.Sym}} {{.Weekday}} {{.Fixed}}",
	})

	if status != 0 {
		t.Fatalf(`runConvert status = %v, stderr %q`, status, errOut)
	}

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != len(verificationTable) {
		t.Fatalf(`got %v lines, want %v`, len(lines), len(verificationTable))
	}

	for index, row := range verificationTable {
		want := row.sym + " " + row.weekday + " " + strconv.Itoa(row.fixedDate)
		if lines[index] != want {
			t.Fatalf(`line %v = %q, want %q`, index, lines[index], want)
		}
	}
}

func TestRunConvertFixedDayNumbers(t *testing.T) {
	for _, row := range verificationTable {
		out, errOut, status := runLines(t, row.gregorian, convertOptions{
			location: time.UTC,
			format:   "{{.Fixed}}",
		})

		if status != 0 {
			t.Fatalf(`runConvert(%q) status = %v, stderr %q`,
				row.gregorian, status, errOut)
		}
		if strings.TrimSpace(out) != strconv.Itoa(row.fixedDate) {
			t.Fatalf(`runConvert(%q) fixed day = %q, want %v`,
				row.gregorian, strings.TrimSpace(out), row.fixedDate)
		}
	}
}

func TestRunConvertEpochSeconds(t *testing.T) {
	out, errOut, status := runLines(t, "0\n86400", convertOptions{
		location: time.UTC,
		format:   "{{.Sym}}",
	})

	if status != 0 {
		t.Fatalf(`runConvert status = %v, stderr %q`, status, errOut)
	}

	want := "1970-01-04\n1970-01-05\n"
	if out != want {
		t.Fatalf(`runConvert = %q, want %q`, out, want)
	}
}

func TestRunConvertFromSym(t *testing.T) {
	input := make([]string, 0, len(verificationTable))
	for _, row := range verificationTable {
		input = append(input, row.sym)
	}

	out, errOut, status := runLines(t, strings.Join(input, "\n"), convertOptions{
		location: time.UTC,
		format:   "{{.Gregorian}}",
		fromSym:  true,
	})

	if status != 0 {
		t.Fatalf(`runConvert status = %v, stderr %q`, status, errOut)
	}

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	for index, row := range verificationTable {
		if lines[index] != row.gregorian {
			t.Fatalf(`line %v = %q, want %q`, index, lines[index], row.gregorian)
		}
	}
}

func TestRunConvertRoundTrip(t *testing.T) {
	for _, row := range verificationTable {
		forward, _, status := runLines(t, row.gregorian, convertOptions{
			location: time.UTC,
			format:   "{{.Sym}}",
		})
		if status != 0 {
			t.Fatalf(`converting %q failed`, row.gregorian)
		}

		back, _, status := runLines(t, forward, convertOptions{
			location: time.UTC,
			format:   "{{.Gregorian}}",
			fromSym:  true,
		})
		if status != 0 {
			t.Fatalf(`converting %q back failed`, forward)
		}

		if strings.TrimSpace(back) != row.gregorian {
			t.Fatalf(`%q round tripped to %q`, row.gregorian, strings.TrimSpace(back))
		}
	}
}

func TestRunConvertJSON(t *testing.T) {
	out, errOut, status := runLines(t, "2004-12-31", convertOptions{
		location: time.UTC,
		asJSON:   true,
	})

	if status != 0 {
		t.Fatalf(`runConvert status = %v, stderr %q`, status, errOut)
	}

	var result conversion
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf(`decoding %q returned %v`, out, err)
	}

	if result.Sym != "2004-12-33" {
		t.Fatalf(`Sym = %q, want "2004-12-33"`, result.Sym)
	}
	if result.Fixed != 731946 {
		t.Fatalf(`Fixed = %v, want 731946`, result.Fixed)
	}
	if result.Weekday != "Friday" {
		t.Fatalf(`Weekday = %q, want "Friday"`, result.Weekday)
	}
	if !result.LeapYear || result.DaysInYear != 371 {
		t.Fatalf(`LeapYear = %v, DaysInYear = %v, want true and 371`,
			result.LeapYear, result.DaysInYear)
	}
	if result.WeekOfYear != 53 || result.Quarter != 4 {
		t.Fatalf(`WeekOfYear = %v, Quarter = %v, want 53 and 4`,
			result.WeekOfYear, result.Quarter)
	}
}

func TestRunConvertLong(t *testing.T) {
	out, _, status := runLines(t, "2004-12-31", convertOptions{
		location: time.UTC,
		long:     true,
	})

	if status != 0 {
		t.Fatalf(`runConvert status = %v`, status)
	}

	for _, want := range []string{
		"Symmetry454   2004-12-33",
		"Fixed day     731946",
		"Long form     Friday, December 33, 2004",
		"Day of year   369 of 371",
		"Week of year  53 of 53",
		"Quarter       4",
		"Leap year     yes",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf(`long output is missing %q, got:\n%s`, want, out)
		}
	}
}

func TestRunConvertSkipsBlankLines(t *testing.T) {
	out, _, status := runLines(t, "\n\n2004-12-31\n\n   \n", convertOptions{
		location: time.UTC,
		format:   "{{.Sym}}",
	})

	if status != 0 {
		t.Fatalf(`runConvert status = %v`, status)
	}
	if out != "2004-12-33\n" {
		t.Fatalf(`runConvert = %q, want "2004-12-33\n"`, out)
	}
}

func TestRunConvertReportsBadLines(t *testing.T) {
	out, errOut, status := runLines(t, "2004-12-31\nnot-a-date\n1776-07-04",
		convertOptions{location: time.UTC, format: "{{.Sym}}"})

	if status != 1 {
		t.Fatalf(`runConvert status = %v, want 1`, status)
	}
	if out != "2004-12-33\n1776-07-04\n" {
		t.Fatalf(`runConvert = %q, want the two good lines`, out)
	}
	if !strings.Contains(errOut, "not-a-date") {
		t.Fatalf(`stderr = %q, want it to name the bad line`, errOut)
	}
}

func TestRunConvertBadFormat(t *testing.T) {
	_, errOut, status := runLines(t, "2004-12-31", convertOptions{
		location: time.UTC,
		format:   "{{.Nope",
	})

	if status != 2 {
		t.Fatalf(`runConvert status = %v, want 2`, status)
	}
	if !strings.Contains(errOut, "invalid format") {
		t.Fatalf(`stderr = %q, want an invalid format message`, errOut)
	}
}

func TestRunConvertDefaultFormats(t *testing.T) {
	out, _, status := runLines(t, "2004-12-31", convertOptions{location: time.UTC})
	if status != 0 {
		t.Fatalf(`runConvert status = %v`, status)
	}
	if !strings.Contains(out, "2004-12-33 Friday, December 33, 2004") {
		t.Fatalf(`default format = %q`, out)
	}

	out, _, status = runLines(t, "2004-12-33", convertOptions{
		location: time.UTC,
		fromSym:  true,
	})
	if status != 0 {
		t.Fatalf(`runConvert status = %v`, status)
	}
	if strings.TrimSpace(out) != "2004-12-33 -> 2004-12-31 Friday" {
		t.Fatalf(`default from-sym format = %q`, out)
	}
}
