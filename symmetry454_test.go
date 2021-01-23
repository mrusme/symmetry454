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
