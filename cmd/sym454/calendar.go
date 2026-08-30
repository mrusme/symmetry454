package main

import (
	"fmt"
	"strings"

	"github.com/mrusme/symmetry454"
)

const (
	innerWidth    = 22
	boxWidth      = innerWidth + 2
	weekdayHeader = " Mo Tu We Th Fr Sa Su "
	maxWeekRows   = 5
	gutter        = " "
)

type styler struct {
	enabled bool
}

func (style styler) highlight(text string) string {
	if !style.enabled {
		return text
	}

	return "\x1b[7m" + text + "\x1b[0m"
}

func centerText(text string, width int) string {
	runes := []rune(text)
	if len(runes) >= width {
		return string(runes[:width])
	}

	left := (width - len(runes)) / 2
	right := width - len(runes) - left

	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}

func shiftMonth(symYear int, symMonth int, offset int) (int, int) {
	index := symYear*12 + symMonth - 1 + offset
	year := index / 12
	month := index % 12
	if month < 0 {
		month += 12
		year--
	}

	return year, month + 1
}

func monthTitle(symYear int, symMonth int, withYear bool) string {
	if withYear {
		return fmt.Sprintf("%s %d", symmetry454.MonthName(symMonth), symYear)
	}

	return symmetry454.MonthName(symMonth)
}

func renderMonth(
	symYear int,
	symMonth int,
	today symmetry454.Sym,
	style styler,
	title string,
	weekRows int,
) []string {
	daysInMonth := symmetry454.SymDaysInMonth(symYear, symMonth)

	lines := make([]string, 0, weekRows+5)
	lines = append(lines, "┌"+strings.Repeat("─", innerWidth)+"┐")
	lines = append(lines, "│"+centerText(title, innerWidth)+"│")
	lines = append(lines, "├"+strings.Repeat("─", innerWidth)+"┤")
	lines = append(lines, "│"+weekdayHeader+"│")

	for day := 1; day <= daysInMonth; day += 7 {
		var row strings.Builder
		row.WriteString("│")
		for offset := 0; offset < 7; offset++ {
			cell := fmt.Sprintf("%2d", day+offset)
			if today.Year == symYear &&
				today.Month == symMonth &&
				today.Day == day+offset {
				cell = style.highlight(cell)
			}
			row.WriteString(" ")
			row.WriteString(cell)
		}
		row.WriteString(" │")
		lines = append(lines, row.String())
	}

	for filled := daysInMonth / 7; filled < weekRows; filled++ {
		lines = append(lines, "│"+strings.Repeat(" ", innerWidth)+"│")
	}

	return append(lines, "└"+strings.Repeat("─", innerWidth)+"┘")
}

func joinBlocks(blocks [][]string, perRow int) string {
	if perRow < 1 {
		perRow = 1
	}

	var out strings.Builder
	for start := 0; start < len(blocks); start += perRow {
		end := start + perRow
		if end > len(blocks) {
			end = len(blocks)
		}

		height := 0
		for _, block := range blocks[start:end] {
			if len(block) > height {
				height = len(block)
			}
		}

		for line := 0; line < height; line++ {
			parts := make([]string, 0, end-start)
			for _, block := range blocks[start:end] {
				if line < len(block) {
					parts = append(parts, block[line])
					continue
				}
				parts = append(parts, strings.Repeat(" ", boxWidth))
			}
			out.WriteString(strings.TrimRight(strings.Join(parts, gutter), " "))
			out.WriteString("\n")
		}

		if end < len(blocks) {
			out.WriteString("\n")
		}
	}

	return out.String()
}

func gregorianSpan(symYear int, symMonth int) string {
	first := symmetry454.Sym{Year: symYear, Month: symMonth, Day: 1}
	last := symmetry454.Sym{
		Year:  symYear,
		Month: symMonth,
		Day:   symmetry454.SymDaysInMonth(symYear, symMonth),
	}

	return fmt.Sprintf("%s to %s Gregorian",
		first.ToTime().Format("2006-01-02"),
		last.ToTime().Format("2006-01-02"))
}

func gregorianYearSpan(symYear int) string {
	first := symmetry454.Sym{Year: symYear, Month: 1, Day: 1}
	last := symmetry454.Sym{
		Year:  symYear,
		Month: 12,
		Day:   symmetry454.SymDaysInMonth(symYear, 12),
	}

	return fmt.Sprintf("%s to %s Gregorian",
		first.ToTime().Format("2006-01-02"),
		last.ToTime().Format("2006-01-02"))
}

func yearSummary(symYear int) string {
	if symmetry454.IsSymLeapYear(symYear) {
		return fmt.Sprintf("%d, leap year of 371 days in 53 weeks", symYear)
	}

	return fmt.Sprintf("%d, common year of 364 days in 52 weeks", symYear)
}

func renderSingleMonth(
	symYear int,
	symMonth int,
	today symmetry454.Sym,
	style styler,
	showGregorian bool,
) string {
	weekRows := symmetry454.SymDaysInMonth(symYear, symMonth) / 7
	block := renderMonth(symYear, symMonth, today, style,
		monthTitle(symYear, symMonth, true), weekRows)

	out := strings.Join(block, "\n") + "\n"
	if showGregorian {
		out += gregorianSpan(symYear, symMonth) + "\n"
	}

	return out
}

func renderThreeMonths(
	symYear int,
	symMonth int,
	today symmetry454.Sym,
	style styler,
	showGregorian bool,
) string {
	blocks := make([][]string, 0, 3)
	for offset := -1; offset <= 1; offset++ {
		year, month := shiftMonth(symYear, symMonth, offset)
		blocks = append(blocks, renderMonth(year, month, today, style,
			monthTitle(year, month, true), maxWeekRows))
	}

	out := joinBlocks(blocks, 3)
	if showGregorian {
		out += gregorianSpan(symYear, symMonth) + "\n"
	}

	return out
}

func renderYear(
	symYear int,
	today symmetry454.Sym,
	style styler,
	perRow int,
	showGregorian bool,
) string {
	if perRow < 1 {
		perRow = 1
	}

	blocks := make([][]string, 0, 12)
	for symMonth := 1; symMonth <= 12; symMonth++ {
		blocks = append(blocks, renderMonth(symYear, symMonth, today, style,
			monthTitle(symYear, symMonth, false), maxWeekRows))
	}

	width := perRow*boxWidth + perRow - 1
	out := strings.TrimRight(centerText(yearSummary(symYear), width), " ") + "\n\n"
	out += joinBlocks(blocks, perRow)
	if showGregorian {
		out += "\n" + gregorianYearSpan(symYear) + "\n"
	}

	return out
}
