package symmetry454

import (
  "testing"
  "time"
)

func TestPriorElapsedDays(t *testing.T) {
  res := PriorElapsedDays(2009)
  if res != 733407 {
    t.Fatalf(`PriorElapsedDays(2009) = %v, want 733407`, res)
  }
}

func TestGregorianOrdinalDay(t *testing.T) {
  res := GregorianOrdinalDay(2012, 7, 14)
  if res != 196 {
    t.Fatalf(`GregorianOrdinalDay(2012, 7, 14) = %v, want 196`, res)
  }
}

func TestGregorianToFixed(t *testing.T) {
  res := GregorianToFixedDate(2004, 12, 31)
  if res != 731946 {
    t.Fatalf(`TestGregorianToFixed(2004, 12, 31) = %v, want 731946`, res)
  }
}

func TestIsSymLeapYear(t *testing.T) {
  res := IsSymLeapYear(2009)
  if res != true {
    t.Fatalf(`IsSymLeapYear(2009) = %v, want true`, res)
  }
}

func TestSymNewYearDay(t *testing.T) {
  res := SymNewYearDay(2010)
  if res != 733776 {
    t.Fatalf(`SymNewYearDay(2010) = %v, want 733776`, res)
  }
}

func TestSymDaysBeforeMonth(t *testing.T) {
  res := SymDaysBeforeMonth(6)
  if res != 154 {
    t.Fatalf(`SymDaysBeforeMonth(6) = %v, want 154`, res)
  }
}

func TestSymDayOfYear(t *testing.T) {
  res := SymDayOfYear(6, 17)
  if res != 171 {
    t.Fatalf(`SymDayOfYear(6, 17) = %v, want 171`, res)
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
}

func TestFixedToSymYear(t *testing.T) {
  res, _ := FixedToSymYear(733774)
  if res != 2009 {
    t.Fatalf(`FixedToSymYear(733774) = %v, want 2009`, res)
  }
}

func TestFixedToSym(t *testing.T) {
  fixed := SymToFixed(2009, 4, 5)
  year, month, day := FixedToSym(fixed)
  if year != 2009 || month != 4 || day != 5 {
    t.Fatalf(`FixedToSym(%v) = (%v, %v, %v), want (2009, 4, 5)`,
      fixed, year, month, day)
  }
}

func TestFixedToWeekdayNum(t *testing.T) {
  res := FixedToWeekdayNum(1461)
  if res != 5 {
    t.Fatalf(`FixedToWeekdayNum(1461) = %v, want 5`, res)
  }
}

func TestFromTime(t *testing.T) {
  tm := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
  res := FromTime(tm)
  if res.Year != 1970 || res.Month != 1 || res.Day != 1 {
    t.Fatalf(`FromTime(t) = (%v, %v, %v), want (1970, 1, 1)`, res.Year, res.Month, res.Day)
  }
}
