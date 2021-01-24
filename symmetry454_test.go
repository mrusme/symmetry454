package symmetry454

import (
  "testing"
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
  res := SymToFixed(2009, 4, 5)
  if res != 733500 {
    t.Fatalf(`SymToFixed(2009, 4, 5) = %v, want 733500`, res)
  }
}

func TestFixedToSymYear(t *testing.T) {
  res, _ := FixedToSymYear(733774)
  if res != 2009 {
    t.Fatalf(`FixedToSymYear(733774) = %v, want 2009`, res)
  }
}

