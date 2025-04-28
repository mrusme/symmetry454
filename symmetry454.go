package symmetry454

import (
  // "fmt"
  "time"
  "math"
)

const GregorianEpoch = 1
const SymEpoch = GregorianEpoch

type Sym struct {
  Year      int
  Month     int
  Day       int
}

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
  gregorianOrdinalDay :=
    math.Floor((367.0 * float64(gregMonth) - 362.0) / 12.0) + float64(gregDay)

  if gregMonth > 2 {
    if IsGregorianLeapYear(gregYear) == true {
      gregorianOrdinalDay -= 1.0
    } else {
      gregorianOrdinalDay -= 2.0
    }
  }

  return int(gregorianOrdinalDay)
}

func GregorianToFixedDate(gregYear int, gregMonth int, gregDay int) (int) {
  year := PriorElapsedDays(gregYear)
  day := GregorianOrdinalDay(gregYear, gregMonth, gregDay)
  return (year + day)
}

func IsSymLeapYear(symYear int) (bool) {
  C := 293.0
  L := 52.0
  K := (C - 1.0) / 2.0
  return math.Mod(L * float64(symYear) + K, C) < L
}

func SymNewYearDay(symYear int) (int) {
  E := float64(symYear - 1)
  fixedDayNumber :=
    SymEpoch + 364.0 * E + 7 * math.Floor((52.0 * E + 146.0) / 293.0)
  return int(fixedDayNumber)
}

func SymDaysBeforeMonth(symMonth int) (int) {
  symMonthF := float64(symMonth)
  symDaysBeforeMonth := 28 * (symMonthF - 1) + 7 * math.Floor(symMonthF / 3)
  return int(symDaysBeforeMonth)
}

func SymDayOfYear(symMonth int, symDay int) (int) {
  return (SymDaysBeforeMonth(symMonth) + symDay);
}

func SymToFixed(symYear int, symMonth int, symDay int) (int) {
  return (SymNewYearDay(symYear) + SymDayOfYear(symMonth, symDay) - 1)
}

func FixedToSymYear(fixedDate int) (int, int) {
  cycleMeanYear := 365.0 + 71.0 / 293.0
  symYear := int((float64(fixedDate) - SymEpoch) / cycleMeanYear)
  startOfYear := SymNewYearDay(symYear)
  if startOfYear < fixedDate {
    if fixedDate - startOfYear >= 364 {
      startOfNextYear := SymNewYearDay(symYear + 1)
      if fixedDate >= startOfNextYear {
        symYear += 1
        startOfYear = startOfNextYear
      }
    }
  } else if startOfYear > fixedDate {
    symYear -= 1
    startOfYear = SymNewYearDay(symYear)
  }

  return symYear, startOfYear
}

func FixedToSym(fixedDate int) (int, int, int) {
  symYear, startOfYear := FixedToSymYear(fixedDate)
  dayOfYear := fixedDate - startOfYear + 1
  weekOfYear := int(math.Ceil(float64(dayOfYear) / 7))
  quarter := int(math.Ceil((4.0 / 53.0) * float64(weekOfYear)))
  dayOfQuarter := dayOfYear - 91 * (quarter - 1)
  weekOfQuarter := int(math.Ceil(float64(dayOfQuarter) / 7.0))
  monthOfQuarter := int(math.Ceil((2.0 / 9.0) * float64(weekOfQuarter)))
  symMonth := 3 * quarter + monthOfQuarter - 3
  daysInYear := 364
  if IsSymLeapYear(symYear) {
    daysInYear = 371
  }
  /*weeksInYear*/ _ = daysInYear / 7
  daysInMonth :=
    int(28 + 7 * math.Floor((math.Mod(float64(symMonth), 3.0) / 2.0)))
  if symMonth == 12 {
    if IsSymLeapYear(symYear) {
      daysInMonth += 7
    }
  }
  /*weeksInMonth*/_ = daysInMonth / 7
  symDay := dayOfYear - SymDaysBeforeMonth(symMonth)

  // fmt.Printf("dayOfYear: %v\nweekOfYear: %v\nquarter: %v\ndayOfQuarter: %v\nweekOfQuarter: %v\nmonthOfQuarter: %v\ndaysInYear: %v\ndaysInMonth: %v\n",
  //   dayOfYear,
  //   weekOfYear,
  //   quarter,
  //   dayOfQuarter,
  //   weekOfQuarter,
  //   monthOfQuarter,
  //   daysInYear,
  //   daysInMonth,
  // )
  return symYear, symMonth, symDay
}

func FixedToWeekdayNum(fixedDate int) (int) {
  _, _, symDay := FixedToSym(fixedDate)
  weekday := int(math.Mod(float64(symDay), 7.0))
  return weekday
}

func FromTime(t time.Time) (Sym) {
  sym := Sym{
    Year: t.Year(),
    Month: int(t.Month()),
    Day: t.Day(),
  }

  return sym
}
