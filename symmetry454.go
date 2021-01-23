package symmetry454

import (
  // "fmt"
  "time"
  "math"
)

const GregorianEpoch = 1
const SymEpoch = GregorianEpoch

func PriorElapsedDays(gregYear int) (int) {
  priorYear := float64(gregYear) - 1.0
  priorElapsedDays :=
    GregorianEpoch +
    priorYear * 365 +
    math.Floor(priorYear / 4) -
    math.Floor(priorYear / 100) +
    math.Floor(priorYear / 400) -
    1
  return int(priorElapsedDays)
}

func IsGregorianLeapYear(gregYear int) (bool) {
  year := time.Date(gregYear, time.December, 31, 0, 0, 0, 0, time.Local)
  days := year.YearDay()

  if days > 365 {
    return true
  }

  return false
}

func GregorianOrdinalDay(gregYear int, gregMonth int, gregDay int) (int) {
  gregorianOrdinalDay := math.Floor((367.0 * float64(gregMonth) - 362.0) / 12.0) +
    float64(gregDay)

  if gregMonth > 2 {
    if IsGregorianLeapYear(gregYear) == true {
      gregorianOrdinalDay -= 1.0
    } else {
      gregorianOrdinalDay -= 2.0
    }
  }

  return int(gregorianOrdinalDay)
}

func IsSymLeapYear(symYear int) (bool) {
  C := 293.0
  L := 52.0
  K := (C - 1.0) / 2.0
  return math.Mod(L * float64(symYear) + K, C) < L
}

func SymNewYearDay(symYear int) (int) {
  E := float64(symYear - 1)
  fixedDayNumber := SymEpoch + 364.0 * E + 7 * math.Floor((52.0 * E + 146.0) / 293.0)
  return int(fixedDayNumber)
}

func SymDaysBeforeMonth(symMonth int) (int) {
  symMonthF := float64(symMonth)
  symDaysBeforeMonth := 28 * (symMonthF - 1) + 7 * math.Floor(symMonthF / 3)
  return int(symDaysBeforeMonth)
}
