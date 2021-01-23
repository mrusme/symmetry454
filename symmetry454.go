package symmetry454

import (
  // "fmt"
  "time"
  "math"
)

const GregorianEpoch = 1

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
