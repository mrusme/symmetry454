package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/mrusme/symmetry454"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func displayWidth(line string) int {
	return len([]rune(ansiPattern.ReplaceAllString(line, "")))
}

func TestCenterText(t *testing.T) {
	tests := []struct {
		text  string
		width int
		want  string
	}{
		{"May", 9, "   May   "},
		{"June", 9, "  June   "},
		{"", 4, "    "},
		{"truncated", 5, "trunc"},
	}

	for _, test := range tests {
		if res := centerText(test.text, test.width); res != test.want {
			t.Fatalf(`centerText(%q, %v) = %q, want %q`,
				test.text, test.width, res, test.want)
		}
	}
}

func TestCenterTextWidth(t *testing.T) {
	for symMonth := 1; symMonth <= 12; symMonth++ {
		title := monthTitle(-100000, symMonth, true)
		if res := len([]rune(centerText(title, innerWidth))); res != innerWidth {
			t.Fatalf(`centerText(%q, %v) is %v wide, want %v`,
				title, innerWidth, res, innerWidth)
		}
	}
}

func TestShiftMonth(t *testing.T) {
	tests := []struct {
		year   int
		month  int
		offset int
		wantY  int
		wantM  int
	}{
		{2026, 8, 0, 2026, 8},
		{2026, 8, 1, 2026, 9},
		{2026, 8, -1, 2026, 7},
		{2026, 12, 1, 2027, 1},
		{2026, 1, -1, 2025, 12},
		{1, 1, -1, 0, 12},
		{0, 1, -1, -1, 12},
		{-1, 12, 1, 0, 1},
		{-121, 1, -1, -122, 12},
		{2026, 1, -13, 2024, 12},
		{2026, 12, 13, 2028, 1},
	}

	for _, test := range tests {
		year, month := shiftMonth(test.year, test.month, test.offset)
		if year != test.wantY || month != test.wantM {
			t.Fatalf(`shiftMonth(%v, %v, %v) = (%v, %v), want (%v, %v)`,
				test.year, test.month, test.offset, year, month,
				test.wantY, test.wantM)
		}
	}
}

func TestShiftMonthRoundTrip(t *testing.T) {
	for year := -50; year <= 50; year++ {
		for month := 1; month <= 12; month++ {
			for offset := -30; offset <= 30; offset++ {
				shiftedY, shiftedM := shiftMonth(year, month, offset)
				backY, backM := shiftMonth(shiftedY, shiftedM, -offset)
				if backY != year || backM != month {
					t.Fatalf(`shiftMonth(%v, %v, %v) does not reverse, got (%v, %v)`,
						year, month, offset, backY, backM)
				}
			}
		}
	}
}

func TestRenderMonthLayout(t *testing.T) {
	style := styler{enabled: false}
	none := symmetry454.Sym{}

	for _, symMonth := range []int{1, 2, 12} {
		for _, symYear := range []int{2009, 2010} {
			block := renderMonth(symYear, symMonth, none, style,
				monthTitle(symYear, symMonth, true), maxWeekRows)

			if len(block) != maxWeekRows+5 {
				t.Fatalf(`renderMonth(%v, %v) has %v lines, want %v`,
					symYear, symMonth, len(block), maxWeekRows+5)
			}

			for index, line := range block {
				if res := displayWidth(line); res != boxWidth {
					t.Fatalf(`renderMonth(%v, %v) line %v is %v wide, want %v: %q`,
						symYear, symMonth, index, res, boxWidth, line)
				}
			}

			if !strings.HasPrefix(block[0], "┌") ||
				!strings.HasSuffix(block[0], "┐") {
				t.Fatalf(`renderMonth(%v, %v) top border is %q`,
					symYear, symMonth, block[0])
			}
			if !strings.HasPrefix(block[len(block)-1], "└") ||
				!strings.HasSuffix(block[len(block)-1], "┘") {
				t.Fatalf(`renderMonth(%v, %v) bottom border is %q`,
					symYear, symMonth, block[len(block)-1])
			}
			if block[3] != "│"+weekdayHeader+"│" {
				t.Fatalf(`renderMonth(%v, %v) weekday row is %q`,
					symYear, symMonth, block[3])
			}
		}
	}
}

func TestRenderMonthDays(t *testing.T) {
	style := styler{enabled: false}
	none := symmetry454.Sym{}

	for symMonth := 1; symMonth <= 12; symMonth++ {
		daysInMonth := symmetry454.SymDaysInMonth(2009, symMonth)
		weekRows := daysInMonth / 7

		block := renderMonth(2009, symMonth, none, style,
			monthTitle(2009, symMonth, true), weekRows)

		cells := make([]string, 0, daysInMonth)
		for _, line := range block[4 : 4+weekRows] {
			cells = append(cells, strings.Fields(strings.Trim(line, "\u2502"))...)
		}

		if len(cells) != daysInMonth {
			t.Fatalf(`renderMonth(2009, %v) rendered %v days, want %v`,
				symMonth, len(cells), daysInMonth)
		}

		for index, cell := range cells {
			if cell != strconv.Itoa(index+1) {
				t.Fatalf(`renderMonth(2009, %v) cell %v is %q, want %q`,
					symMonth, index, cell, strconv.Itoa(index+1))
			}
		}
	}
}

func TestRenderMonthHighlight(t *testing.T) {
	today := symmetry454.Sym{Year: 2026, Month: 8, Day: 35}

	block := renderMonth(2026, 8, today, styler{enabled: true},
		monthTitle(2026, 8, true), 5)
	joined := strings.Join(block, "\n")

	if strings.Count(joined, "\x1b[7m") != 1 {
		t.Fatalf(`want exactly one highlighted cell, got %v`,
			strings.Count(joined, "\x1b[7m"))
	}
	if !strings.Contains(joined, "\x1b[7m35\x1b[0m") {
		t.Fatalf(`want day 35 highlighted, got %q`, joined)
	}

	plain := renderMonth(2026, 8, today, styler{enabled: false},
		monthTitle(2026, 8, true), 5)
	if strings.Contains(strings.Join(plain, "\n"), "\x1b[") {
		t.Fatalf(`disabled styler still emitted an escape sequence`)
	}

	other := renderMonth(2026, 9, today, styler{enabled: true},
		monthTitle(2026, 9, true), 5)
	if strings.Contains(strings.Join(other, "\n"), "\x1b[") {
		t.Fatalf(`a month that is not today's was highlighted`)
	}

	otherYear := renderMonth(2025, 8, today, styler{enabled: true},
		monthTitle(2025, 8, true), 5)
	if strings.Contains(strings.Join(otherYear, "\n"), "\x1b[") {
		t.Fatalf(`a year that is not today's was highlighted`)
	}
}

func TestRenderYearWidth(t *testing.T) {
	today := symmetry454.Sym{Year: 2026, Month: 8, Day: 35}

	for _, perRow := range []int{1, 2, 3, 4, 6, 12} {
		out := renderYear(2026, today, styler{enabled: false}, perRow, false)
		want := perRow*boxWidth + perRow - 1

		for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
			if res := displayWidth(line); res > want {
				t.Fatalf(`renderYear(perRow %v) line is %v wide, want at most %v: %q`,
					perRow, res, want, line)
			}
		}
	}
}

func TestRenderYearHasEveryMonth(t *testing.T) {
	today := symmetry454.Sym{Year: 2026, Month: 8, Day: 35}
	out := renderYear(2026, today, styler{enabled: false}, 3, false)

	for symMonth := 1; symMonth <= 12; symMonth++ {
		if !strings.Contains(out, symmetry454.MonthName(symMonth)) {
			t.Fatalf(`renderYear is missing %v`, symmetry454.MonthName(symMonth))
		}
	}

	if strings.Count(out, "Mo Tu We Th Fr Sa Su") != 12 {
		t.Fatalf(`renderYear has %v weekday rows, want 12`,
			strings.Count(out, "Mo Tu We Th Fr Sa Su"))
	}
}

func TestRenderYearLeapWeek(t *testing.T) {
	none := symmetry454.Sym{}
	style := styler{enabled: false}

	leap := renderYear(2026, none, style, 3, false)
	if !strings.Contains(leap, "leap year of 371 days in 53 weeks") {
		t.Fatalf(`renderYear(2026) is missing the leap year summary`)
	}

	common := renderYear(2027, none, style, 3, false)
	if !strings.Contains(common, "common year of 364 days in 52 weeks") {
		t.Fatalf(`renderYear(2027) is missing the common year summary`)
	}
}

func TestRenderSingleMonthWeekRows(t *testing.T) {
	none := symmetry454.Sym{}
	style := styler{enabled: false}

	short := renderSingleMonth(2010, 1, none, style, false)
	if res := len(strings.Split(strings.TrimRight(short, "\n"), "\n")); res != 9 {
		t.Fatalf(`a 28 day month rendered %v lines, want 9`, res)
	}

	long := renderSingleMonth(2010, 2, none, style, false)
	if res := len(strings.Split(strings.TrimRight(long, "\n"), "\n")); res != 10 {
		t.Fatalf(`a 35 day month rendered %v lines, want 10`, res)
	}

	leapDecember := renderSingleMonth(2026, 12, none, style, false)
	if !strings.Contains(leapDecember, "29 30 31 32 33 34 35") {
		t.Fatalf(`December of a leap year is missing its leap week`)
	}

	commonDecember := renderSingleMonth(2027, 12, none, style, false)
	if strings.Contains(commonDecember, "29 30 31 32 33 34 35") {
		t.Fatalf(`December of a common year has a leap week`)
	}
}

func TestRenderThreeMonths(t *testing.T) {
	none := symmetry454.Sym{}
	out := renderThreeMonths(2027, 1, none, styler{enabled: false}, false)

	for _, want := range []string{"December 2026", "January 2027", "February 2027"} {
		if !strings.Contains(out, want) {
			t.Fatalf(`renderThreeMonths is missing %q`, want)
		}
	}
}

func TestGregorianSpan(t *testing.T) {
	if res := gregorianSpan(2026, 8); res != "2026-07-27 to 2026-08-30 Gregorian" {
		t.Fatalf(`gregorianSpan(2026, 8) = %q`, res)
	}

	if res := gregorianYearSpan(2026); res != "2025-12-29 to 2027-01-03 Gregorian" {
		t.Fatalf(`gregorianYearSpan(2026) = %q`, res)
	}
}

func TestColorEnabled(t *testing.T) {
	if !colorEnabled("always", nil) {
		t.Fatalf(`colorEnabled("always") = false, want true`)
	}
	if colorEnabled("never", nil) {
		t.Fatalf(`colorEnabled("never") = true, want false`)
	}

	previous, had := os.LookupEnv("NO_COLOR")
	defer func() {
		if had {
			os.Setenv("NO_COLOR", previous)
			return
		}
		os.Unsetenv("NO_COLOR")
	}()

	os.Setenv("NO_COLOR", "1")
	if colorEnabled("auto", nil) {
		t.Fatalf(`colorEnabled("auto") with NO_COLOR = true, want false`)
	}
	if !colorEnabled("always", nil) {
		t.Fatalf(`colorEnabled("always") with NO_COLOR = false, want true`)
	}

	os.Unsetenv("NO_COLOR")
	if colorEnabled("auto", nil) {
		t.Fatalf(`colorEnabled("auto") on a non terminal = true, want false`)
	}
}
