package symmetry454

import (
	"testing"
	"time"
)

type verificationRow struct {
	gregYear   int
	gregMonth  int
	gregDay    int
	fixedDate  int
	weekdayNum int
	symYear    int
	symMonth   int
	symDay     int
}

var verificationTable = []verificationRow{
	{-121, 4, 26, -44444, 6, -121, 4, 27},
	{-91, 9, 27, -33333, 1, -91, 9, 22},
	{122, 9, 7, 44444, 1, 122, 9, 8},
	{1776, 7, 4, 648491, 4, 1776, 7, 4},
	{1867, 7, 1, 681724, 1, 1867, 7, 1},
	{1947, 10, 24, 711058, 5, 1947, 10, 26},
	{1995, 8, 10, 728515, 4, 1995, 8, 11},
	{2000, 2, 29, 730179, 2, 2000, 2, 30},
	{2004, 5, 2, 731703, 0, 2004, 5, 7},
	{2004, 12, 31, 731946, 5, 2004, 12, 33},
	{2020, 2, 20, 737475, 4, 2020, 2, 25},
	{2222, 2, 2, 811236, 6, 2222, 2, 6},
	{3333, 3, 1, 1217048, 0, 3333, 2, 35},
}

func TestPriorElapsedDays(t *testing.T) {
	res := PriorElapsedDays(2009)
	if res != 733407 {
		t.Fatalf(`PriorElapsedDays(2009) = %v, want 733407`, res)
	}
}

func TestIsGregorianLeapYear(t *testing.T) {
	for _, gregYear := range []int{2000, 2004, 1996, -4, 400, -400} {
		if IsGregorianLeapYear(gregYear) != true {
			t.Fatalf(`IsGregorianLeapYear(%v) = false, want true`, gregYear)
		}
	}

	for _, gregYear := range []int{1900, 2001, 2100, -1, -100, 300} {
		if IsGregorianLeapYear(gregYear) != false {
			t.Fatalf(`IsGregorianLeapYear(%v) = true, want false`, gregYear)
		}
	}
}

func TestGregorianOrdinalDay(t *testing.T) {
	res := GregorianOrdinalDay(2012, 7, 14)
	if res != 196 {
		t.Fatalf(`GregorianOrdinalDay(2012, 7, 14) = %v, want 196`, res)
	}

	res = GregorianOrdinalDay(2011, 7, 14)
	if res != 195 {
		t.Fatalf(`GregorianOrdinalDay(2011, 7, 14) = %v, want 195`, res)
	}
}

func TestGregorianToFixed(t *testing.T) {
	res := GregorianToFixedDate(2004, 12, 31)
	if res != 731946 {
		t.Fatalf(`TestGregorianToFixed(2004, 12, 31) = %v, want 731946`, res)
	}

	for _, row := range verificationTable {
		res := GregorianToFixedDate(row.gregYear, row.gregMonth, row.gregDay)
		if res != row.fixedDate {
			t.Fatalf(`GregorianToFixedDate(%v, %v, %v) = %v, want %v`,
				row.gregYear, row.gregMonth, row.gregDay, res, row.fixedDate)
		}
	}
}

func TestFixedToGregorian(t *testing.T) {
	for _, row := range verificationTable {
		gregYear, gregMonth, gregDay := FixedToGregorian(row.fixedDate)
		if gregYear != row.gregYear ||
			gregMonth != row.gregMonth ||
			gregDay != row.gregDay {
			t.Fatalf(`FixedToGregorian(%v) = (%v, %v, %v), want (%v, %v, %v)`,
				row.fixedDate, gregYear, gregMonth, gregDay,
				row.gregYear, row.gregMonth, row.gregDay)
		}
	}
}

func TestGregorianRoundTrip(t *testing.T) {
	for fixedDate := -400000; fixedDate <= 1300000; fixedDate++ {
		gregYear, gregMonth, gregDay := FixedToGregorian(fixedDate)
		back := GregorianToFixedDate(gregYear, gregMonth, gregDay)
		if back != fixedDate {
			t.Fatalf(`FixedToGregorian(%v) = (%v, %v, %v), converts back to %v`,
				fixedDate, gregYear, gregMonth, gregDay, back)
		}
	}
}

func TestIsSymLeapYear(t *testing.T) {
	res := IsSymLeapYear(2009)
	if res != true {
		t.Fatalf(`IsSymLeapYear(2009) = %v, want true`, res)
	}

	res = IsSymLeapYear(2010)
	if res != false {
		t.Fatalf(`IsSymLeapYear(2010) = %v, want false`, res)
	}
}

func TestIsSymLeapYearBeforeEpoch(t *testing.T) {
	for _, symYear := range []int{-2, -8, -14, -19, -25, -30, -36} {
		if IsSymLeapYear(symYear) != true {
			t.Fatalf(`IsSymLeapYear(%v) = false, want true`, symYear)
		}
	}

	for _, symYear := range []int{-1, -3, -4, -5, -6, -7, -9, -13, -15} {
		if IsSymLeapYear(symYear) != false {
			t.Fatalf(`IsSymLeapYear(%v) = true, want false`, symYear)
		}
	}
}

func TestIsSymLeapYearCycle(t *testing.T) {
	for _, start := range []int{1, 294, -292, -1000, -5000} {
		count := 0
		for symYear := start; symYear < start+293; symYear++ {
			if IsSymLeapYear(symYear) {
				count++
			}
		}
		if count != 52 {
			t.Fatalf(`years [%v, %v) contain %v leap years, want 52`,
				start, start+293, count)
		}
	}
}

func TestSymNewYearDay(t *testing.T) {
	res := SymNewYearDay(2010)
	if res != 733776 {
		t.Fatalf(`SymNewYearDay(2010) = %v, want 733776`, res)
	}

	res = SymNewYearDay(2009)
	if res != 733405 {
		t.Fatalf(`SymNewYearDay(2009) = %v, want 733405`, res)
	}
}

func TestSymYearLength(t *testing.T) {
	for symYear := -3000; symYear <= 4000; symYear++ {
		length := SymNewYearDay(symYear+1) - SymNewYearDay(symYear)
		want := 364
		if IsSymLeapYear(symYear) {
			want = 371
		}
		if length != want {
			t.Fatalf(`year %v is %v days long, want %v`, symYear, length, want)
		}
	}
}

func TestSymDaysBeforeMonth(t *testing.T) {
	res := SymDaysBeforeMonth(6)
	if res != 154 {
		t.Fatalf(`SymDaysBeforeMonth(6) = %v, want 154`, res)
	}

	want := []int{0, 28, 63, 91, 119, 154, 182, 210, 245, 273, 301, 336}
	for symMonth := 1; symMonth <= 12; symMonth++ {
		res := SymDaysBeforeMonth(symMonth)
		if res != want[symMonth-1] {
			t.Fatalf(`SymDaysBeforeMonth(%v) = %v, want %v`,
				symMonth, res, want[symMonth-1])
		}
	}
}

func TestSymDayOfYear(t *testing.T) {
	res := SymDayOfYear(6, 17)
	if res != 171 {
		t.Fatalf(`SymDayOfYear(6, 17) = %v, want 171`, res)
	}
}

func TestSymDaysInMonth(t *testing.T) {
	want := []int{28, 35, 28, 28, 35, 28, 28, 35, 28, 28, 35, 28}
	for symMonth := 1; symMonth <= 12; symMonth++ {
		res := SymDaysInMonth(2010, symMonth)
		if res != want[symMonth-1] {
			t.Fatalf(`SymDaysInMonth(2010, %v) = %v, want %v`,
				symMonth, res, want[symMonth-1])
		}
	}

	res := SymDaysInMonth(2009, 12)
	if res != 35 {
		t.Fatalf(`SymDaysInMonth(2009, 12) = %v, want 35`, res)
	}
}

func TestSymToFixed(t *testing.T) {
	res := SymToFixed(-121, 4, 27)
	if res != -44444 {
		t.Fatalf(`SymToFixed(-121, 4, 27) = %v, want -44444`, res)
	}

	res = SymToFixed(1776, 7, 4)
	if res != 648491 {
		t.Fatalf(`SymToFixed(1776, 7, 4) = %v, want 648491`, res)
	}

	res = SymToFixed(2009, 4, 5)
	if res != 733500 {
		t.Fatalf(`SymToFixed(2009, 4, 5) = %v, want 733500`, res)
	}

	res = SymToFixed(3333, 2, 35)
	if res != 1217048 {
		t.Fatalf(`SymToFixed(3333, 2, 35) = %v, want 1217048`, res)
	}

	for _, row := range verificationTable {
		res := SymToFixed(row.symYear, row.symMonth, row.symDay)
		if res != row.fixedDate {
			t.Fatalf(`SymToFixed(%v, %v, %v) = %v, want %v`,
				row.symYear, row.symMonth, row.symDay, res, row.fixedDate)
		}
	}
}

func TestFixedToSymYear(t *testing.T) {
	res, _ := FixedToSymYear(733774)
	if res != 2009 {
		t.Fatalf(`FixedToSymYear(733774) = %v, want 2009`, res)
	}

	res, _ = FixedToSymYear(733406)
	if res != 2009 {
		t.Fatalf(`FixedToSymYear(733406) = %v, want 2009`, res)
	}

	res, startOfYear := FixedToSymYear(733649)
	if res != 2009 || startOfYear != 733405 {
		t.Fatalf(`FixedToSymYear(733649) = (%v, %v), want (2009, 733405)`,
			res, startOfYear)
	}
}

func TestFixedToSymYearBrackets(t *testing.T) {
	for fixedDate := -900000; fixedDate <= 1400000; fixedDate++ {
		symYear, startOfYear := FixedToSymYear(fixedDate)
		if startOfYear != SymNewYearDay(symYear) {
			t.Fatalf(`FixedToSymYear(%v) = (%v, %v), but SymNewYearDay(%v) = %v`,
				fixedDate, symYear, startOfYear, symYear, SymNewYearDay(symYear))
		}
		if fixedDate < startOfYear || fixedDate >= SymNewYearDay(symYear+1) {
			t.Fatalf(`FixedToSymYear(%v) = %v, whose year spans [%v, %v)`,
				fixedDate, symYear, startOfYear, SymNewYearDay(symYear+1))
		}
	}
}

func TestFixedToSym(t *testing.T) {
	fixed := SymToFixed(2009, 4, 5)
	year, month, day := FixedToSym(fixed)
	if year != 2009 || month != 4 || day != 5 {
		t.Fatalf(`FixedToSym(%v) = (%v, %v, %v), want (2009, 4, 5)`,
			fixed, year, month, day)
	}

	for _, row := range verificationTable {
		year, month, day := FixedToSym(row.fixedDate)
		if year != row.symYear || month != row.symMonth || day != row.symDay {
			t.Fatalf(`FixedToSym(%v) = (%v, %v, %v), want (%v, %v, %v)`,
				row.fixedDate, year, month, day,
				row.symYear, row.symMonth, row.symDay)
		}
	}
}

func TestFixedToSymLeapWeek(t *testing.T) {
	for fixedDate := -900000; fixedDate <= 1400000; fixedDate++ {
		symYear, symMonth, symDay := FixedToSym(fixedDate)
		if symMonth < 1 || symMonth > 12 {
			t.Fatalf(`FixedToSym(%v) = (%v, %v, %v), month out of range`,
				fixedDate, symYear, symMonth, symDay)
		}
		if symDay < 1 || symDay > SymDaysInMonth(symYear, symMonth) {
			t.Fatalf(`FixedToSym(%v) = (%v, %v, %v), day out of range 1..%v`,
				fixedDate, symYear, symMonth, symDay,
				SymDaysInMonth(symYear, symMonth))
		}
	}
}

func TestFixedToSymRoundTrip(t *testing.T) {
	for fixedDate := -900000; fixedDate <= 1400000; fixedDate++ {
		symYear, symMonth, symDay := FixedToSym(fixedDate)
		back := SymToFixed(symYear, symMonth, symDay)
		if back != fixedDate {
			t.Fatalf(`FixedToSym(%v) = (%v, %v, %v), converts back to %v`,
				fixedDate, symYear, symMonth, symDay, back)
		}
	}
}

func TestSymToFixedRoundTrip(t *testing.T) {
	for symYear := -2500; symYear <= 4000; symYear++ {
		for symMonth := 1; symMonth <= 12; symMonth++ {
			for symDay := 1; symDay <= SymDaysInMonth(symYear, symMonth); symDay++ {
				fixedDate := SymToFixed(symYear, symMonth, symDay)
				year, month, day := FixedToSym(fixedDate)
				if year != symYear || month != symMonth || day != symDay {
					t.Fatalf(`(%v, %v, %v) -> %v -> (%v, %v, %v)`,
						symYear, symMonth, symDay, fixedDate, year, month, day)
				}
			}
		}
	}
}

func TestFixedToWeekdayNum(t *testing.T) {
	res := FixedToWeekdayNum(1461)
	if res != 5 {
		t.Fatalf(`FixedToWeekdayNum(1461) = %v, want 5`, res)
	}

	for _, row := range verificationTable {
		res := FixedToWeekdayNum(row.fixedDate)
		if res != row.weekdayNum {
			t.Fatalf(`FixedToWeekdayNum(%v) = %v, want %v`,
				row.fixedDate, res, row.weekdayNum)
		}
	}
}

func TestFixedToWeekdayNumRange(t *testing.T) {
	for fixedDate := -900000; fixedDate <= 1400000; fixedDate++ {
		res := FixedToWeekdayNum(fixedDate)
		if res < 0 || res > 6 {
			t.Fatalf(`FixedToWeekdayNum(%v) = %v, out of range 0..6`, fixedDate, res)
		}
	}
}

func TestSymNewYearDayIsMonday(t *testing.T) {
	for symYear := -3000; symYear <= 4000; symYear++ {
		res := FixedToWeekdayNum(SymNewYearDay(symYear))
		if res != 1 {
			t.Fatalf(`New Year Day of %v falls on weekday %v, want 1`, symYear, res)
		}
	}
}

func TestFromTime(t *testing.T) {
	tm := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	res := FromTime(tm)
	if res.Year != 1970 || res.Month != 1 || res.Day != 4 {
		t.Fatalf(`FromTime(t) = (%v, %v, %v), want (1970, 1, 4)`,
			res.Year, res.Month, res.Day)
	}

	tm = time.Date(2004, 12, 31, 0, 0, 0, 0, time.UTC)
	res = FromTime(tm)
	if res.Year != 2004 || res.Month != 12 || res.Day != 33 {
		t.Fatalf(`FromTime(t) = (%v, %v, %v), want (2004, 12, 33)`,
			res.Year, res.Month, res.Day)
	}
}

func TestToTime(t *testing.T) {
	for _, row := range verificationTable {
		sym := Sym{Year: row.symYear, Month: row.symMonth, Day: row.symDay}
		tm := sym.ToTime()
		if tm.Year() != row.gregYear ||
			int(tm.Month()) != row.gregMonth ||
			tm.Day() != row.gregDay {
			t.Fatalf(`Sym{%v, %v, %v}.ToTime() = %v, want %v-%v-%v`,
				row.symYear, row.symMonth, row.symDay, tm.Format("2006-01-02"),
				row.gregYear, row.gregMonth, row.gregDay)
		}
	}
}

func TestTimeRoundTrip(t *testing.T) {
	tm := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 80000; i++ {
		res := FromTime(tm)
		back := res.ToTime()
		if !back.Equal(tm) {
			t.Fatalf(`FromTime(%v) = (%v, %v, %v), converts back to %v`,
				tm.Format("2006-01-02"), res.Year, res.Month, res.Day,
				back.Format("2006-01-02"))
		}
		tm = tm.AddDate(0, 0, 1)
	}
}

func TestSymDaysInYear(t *testing.T) {
	if res := SymDaysInYear(2009); res != 371 {
		t.Fatalf(`SymDaysInYear(2009) = %v, want 371`, res)
	}

	if res := SymDaysInYear(2010); res != 364 {
		t.Fatalf(`SymDaysInYear(2010) = %v, want 364`, res)
	}
}

func TestMonthName(t *testing.T) {
	want := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	for symMonth := 1; symMonth <= 12; symMonth++ {
		if res := MonthName(symMonth); res != want[symMonth-1] {
			t.Fatalf(`MonthName(%v) = %q, want %q`, symMonth, res, want[symMonth-1])
		}
	}

	for _, symMonth := range []int{-1, 0, 13, 14} {
		if res := MonthName(symMonth); res != "" {
			t.Fatalf(`MonthName(%v) = %q, want ""`, symMonth, res)
		}
	}
}

func TestFromFixed(t *testing.T) {
	for _, row := range verificationTable {
		sym := FromFixed(row.fixedDate)
		if sym.Year != row.symYear ||
			sym.Month != row.symMonth ||
			sym.Day != row.symDay {
			t.Fatalf(`FromFixed(%v) = (%v, %v, %v), want (%v, %v, %v)`,
				row.fixedDate, sym.Year, sym.Month, sym.Day,
				row.symYear, row.symMonth, row.symDay)
		}
	}
}

func TestSymFixed(t *testing.T) {
	for _, row := range verificationTable {
		sym := Sym{Year: row.symYear, Month: row.symMonth, Day: row.symDay}
		if res := sym.Fixed(); res != row.fixedDate {
			t.Fatalf(`Sym{%v, %v, %v}.Fixed() = %v, want %v`,
				row.symYear, row.symMonth, row.symDay, res, row.fixedDate)
		}
	}
}

func TestSymWeekday(t *testing.T) {
	for _, row := range verificationTable {
		sym := Sym{Year: row.symYear, Month: row.symMonth, Day: row.symDay}
		if res := int(sym.Weekday()); res != row.weekdayNum {
			t.Fatalf(`Sym{%v, %v, %v}.Weekday() = %v, want %v`,
				row.symYear, row.symMonth, row.symDay, res, row.weekdayNum)
		}
	}
}

func TestSymWeekdayMatchesDayOfMonth(t *testing.T) {
	for fixedDate := -900000; fixedDate <= 1400000; fixedDate++ {
		sym := FromFixed(fixedDate)
		if int(sym.Weekday()) != sym.Day%7 {
			t.Fatalf(`Sym{%v, %v, %v}.Weekday() = %v, want %v`,
				sym.Year, sym.Month, sym.Day, sym.Weekday(), sym.Day%7)
		}
	}
}

func TestSymValid(t *testing.T) {
	valid := []Sym{
		{2009, 12, 35},
		{2010, 12, 28},
		{2010, 1, 28},
		{2010, 2, 35},
		{2010, 11, 35},
		{-121, 4, 27},
	}
	for _, sym := range valid {
		if !sym.Valid() {
			t.Fatalf(`Sym{%v, %v, %v}.Valid() = false, want true`,
				sym.Year, sym.Month, sym.Day)
		}
	}

	invalid := []Sym{
		{2010, 1, 29},
		{2010, 2, 36},
		{2010, 12, 29},
		{2009, 12, 36},
		{2010, 0, 1},
		{2010, 13, 1},
		{2010, 1, 0},
		{2010, 1, -1},
	}
	for _, sym := range invalid {
		if sym.Valid() {
			t.Fatalf(`Sym{%v, %v, %v}.Valid() = true, want false`,
				sym.Year, sym.Month, sym.Day)
		}
	}
}

func TestSymValidMatchesRoundTrip(t *testing.T) {
	for symYear := 1990; symYear <= 2050; symYear++ {
		for symMonth := 1; symMonth <= 12; symMonth++ {
			for symDay := 1; symDay <= 40; symDay++ {
				sym := Sym{Year: symYear, Month: symMonth, Day: symDay}
				want := FromFixed(sym.Fixed()) == sym
				if sym.Valid() != want {
					t.Fatalf(`Sym{%v, %v, %v}.Valid() = %v, round trip says %v`,
						symYear, symMonth, symDay, sym.Valid(), want)
				}
			}
		}
	}
}

func TestSymAccessors(t *testing.T) {
	sym := Sym{Year: 2009, Month: 12, Day: 33}

	if res := sym.DayOfYear(); res != 369 {
		t.Fatalf(`DayOfYear() = %v, want 369`, res)
	}
	if res := sym.WeekOfYear(); res != 53 {
		t.Fatalf(`WeekOfYear() = %v, want 53`, res)
	}
	if res := sym.Quarter(); res != 4 {
		t.Fatalf(`Quarter() = %v, want 4`, res)
	}
	if res := sym.DayOfQuarter(); res != 96 {
		t.Fatalf(`DayOfQuarter() = %v, want 96`, res)
	}
	if res := sym.WeekOfQuarter(); res != 14 {
		t.Fatalf(`WeekOfQuarter() = %v, want 14`, res)
	}
	if res := sym.MonthOfQuarter(); res != 3 {
		t.Fatalf(`MonthOfQuarter() = %v, want 3`, res)
	}
	if res := sym.WeekOfMonth(); res != 5 {
		t.Fatalf(`WeekOfMonth() = %v, want 5`, res)
	}
	if res := sym.DaysInMonth(); res != 35 {
		t.Fatalf(`DaysInMonth() = %v, want 35`, res)
	}
	if res := sym.WeeksInMonth(); res != 5 {
		t.Fatalf(`WeeksInMonth() = %v, want 5`, res)
	}
	if res := sym.DaysInYear(); res != 371 {
		t.Fatalf(`DaysInYear() = %v, want 371`, res)
	}
	if res := sym.WeeksInYear(); res != 53 {
		t.Fatalf(`WeeksInYear() = %v, want 53`, res)
	}
	if res := sym.StartOfYear(); res != 733405 {
		t.Fatalf(`StartOfYear() = %v, want 733405`, res)
	}
	if !sym.IsLeapYear() {
		t.Fatalf(`IsLeapYear() = false, want true`)
	}
	if res := sym.MonthName(); res != "December" {
		t.Fatalf(`MonthName() = %q, want "December"`, res)
	}
}

func TestSymAccessorsConsistency(t *testing.T) {
	for fixedDate := -900000; fixedDate <= 1400000; fixedDate++ {
		sym := FromFixed(fixedDate)

		if res := sym.DayOfYear(); res != fixedDate-sym.StartOfYear()+1 {
			t.Fatalf(`Sym{%v, %v, %v}.DayOfYear() = %v, want %v`,
				sym.Year, sym.Month, sym.Day, res, fixedDate-sym.StartOfYear()+1)
		}

		if res := sym.Quarter(); res != (sym.Month-1)/3+1 {
			t.Fatalf(`Sym{%v, %v, %v}.Quarter() = %v, want %v`,
				sym.Year, sym.Month, sym.Day, res, (sym.Month-1)/3+1)
		}

		if res := sym.MonthOfQuarter(); res != (sym.Month-1)%3+1 {
			t.Fatalf(`Sym{%v, %v, %v}.MonthOfQuarter() = %v, want %v`,
				sym.Year, sym.Month, sym.Day, res, (sym.Month-1)%3+1)
		}

		if res := 3*sym.Quarter() + sym.MonthOfQuarter() - 3; res != sym.Month {
			t.Fatalf(`Sym{%v, %v, %v} quarter/month mismatch, got month %v`,
				sym.Year, sym.Month, sym.Day, res)
		}

		if res := sym.DayOfQuarter(); res < 1 || res > 98 {
			t.Fatalf(`Sym{%v, %v, %v}.DayOfQuarter() = %v, out of range 1..98`,
				sym.Year, sym.Month, sym.Day, res)
		}

		if res := sym.WeekOfYear(); res < 1 || res > sym.WeeksInYear() {
			t.Fatalf(`Sym{%v, %v, %v}.WeekOfYear() = %v, out of range 1..%v`,
				sym.Year, sym.Month, sym.Day, res, sym.WeeksInYear())
		}

		if res := sym.WeekOfMonth(); res < 1 || res > sym.WeeksInMonth() {
			t.Fatalf(`Sym{%v, %v, %v}.WeekOfMonth() = %v, out of range 1..%v`,
				sym.Year, sym.Month, sym.Day, res, sym.WeeksInMonth())
		}

		if !sym.Valid() {
			t.Fatalf(`FromFixed(%v) = Sym{%v, %v, %v}, which is not valid`,
				fixedDate, sym.Year, sym.Month, sym.Day)
		}
	}
}

func TestSymDaysInMonthSumsToYear(t *testing.T) {
	for symYear := -3000; symYear <= 4000; symYear++ {
		total := 0
		for symMonth := 1; symMonth <= 12; symMonth++ {
			total += SymDaysInMonth(symYear, symMonth)
		}
		if total != SymDaysInYear(symYear) {
			t.Fatalf(`year %v months total %v days, want %v`,
				symYear, total, SymDaysInYear(symYear))
		}
	}
}

func TestSymString(t *testing.T) {
	tests := []struct {
		sym  Sym
		want string
	}{
		{Sym{2009, 12, 33}, "2009-12-33"},
		{Sym{2026, 8, 35}, "2026-08-35"},
		{Sym{-121, 4, 27}, "-0121-04-27"},
		{Sym{1, 1, 1}, "0001-01-01"},
	}
	for _, test := range tests {
		if res := test.sym.String(); res != test.want {
			t.Fatalf(`Sym{%v, %v, %v}.String() = %q, want %q`,
				test.sym.Year, test.sym.Month, test.sym.Day, res, test.want)
		}
	}
}

func TestToTimeIn(t *testing.T) {
	loc := time.FixedZone("CEST", 2*3600)
	sym := Sym{Year: 2026, Month: 8, Day: 35}
	tm := sym.ToTimeIn(loc)

	if tm.Year() != 2026 || tm.Month() != time.August || tm.Day() != 30 {
		t.Fatalf(`ToTimeIn() = %v, want 2026-08-30`, tm.Format("2006-01-02"))
	}
	if tm.Location() != loc {
		t.Fatalf(`ToTimeIn() location = %v, want %v`, tm.Location(), loc)
	}

	if !sym.ToTime().Equal(sym.ToTimeIn(time.UTC)) {
		t.Fatalf(`ToTime() and ToTimeIn(time.UTC) disagree`)
	}
}

func TestFromTimeZoned(t *testing.T) {
	loc := time.FixedZone("CEST", 2*3600)
	tm := time.Date(2026, 8, 30, 23, 30, 0, 0, loc)

	sym := FromTime(tm)
	if sym.Year != 2026 || sym.Month != 8 || sym.Day != 35 {
		t.Fatalf(`FromTime(%v) = (%v, %v, %v), want (2026, 8, 35)`,
			tm.Format(time.RFC3339), sym.Year, sym.Month, sym.Day)
	}

	back := sym.ToTimeIn(loc)
	if back.Year() != tm.Year() ||
		back.Month() != tm.Month() ||
		back.Day() != tm.Day() {
		t.Fatalf(`round trip = %v, want the calendar date of %v`,
			back.Format(time.RFC3339), tm.Format(time.RFC3339))
	}
}
