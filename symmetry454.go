package symmetry454

import (
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

func modulus(x float64, y float64) (float64) {
  return x - y * math.Floor(x / y)
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
  if gregYear % 4 != 0 {
    return false
  }

  if gregYear % 100 != 0 {
    return true
  }

  return gregYear % 400 == 0
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

func FixedToGregorian(fixedDate int) (int, int, int) {
  d0 := float64(fixedDate - GregorianEpoch)
  n400 := math.Floor(d0 / 146097.0)
  d1 := modulus(d0, 146097.0)
  n100 := math.Floor(d1 / 36524.0)
  d2 := modulus(d1, 36524.0)
  n4 := math.Floor(d2 / 1461.0)
  d3 := modulus(d2, 1461.0)
  n1 := math.Floor(d3 / 365.0)

  gregYear := int(400.0 * n400 + 100.0 * n100 + 4.0 * n4 + n1)
  if n100 != 4.0 && n1 != 4.0 {
    gregYear += 1
  }

  priorDays := fixedDate - GregorianToFixedDate(gregYear, 1, 1)

  correction := 2
  if fixedDate < GregorianToFixedDate(gregYear, 3, 1) {
    correction = 0
  } else if IsGregorianLeapYear(gregYear) {
    correction = 1
  }

  gregMonth :=
    int(math.Floor((12.0 * float64(priorDays + correction) + 373.0) / 367.0))
  gregDay := fixedDate - GregorianToFixedDate(gregYear, gregMonth, 1) + 1

  return gregYear, gregMonth, gregDay
}

func IsSymLeapYear(symYear int) (bool) {
  C := 293.0
  L := 52.0
  K := (C - 1.0) / 2.0
  return modulus(L * float64(symYear) + K, C) < L
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

func SymDaysInMonth(symYear int, symMonth int) (int) {
  daysInMonth :=
    28 + 7 * int(math.Floor(modulus(float64(symMonth), 3.0) / 2.0))

  if symMonth == 12 && IsSymLeapYear(symYear) {
    daysInMonth += 7
  }

  return daysInMonth
}

func SymToFixed(symYear int, symMonth int, symDay int) (int) {
  return (SymNewYearDay(symYear) + SymDayOfYear(symMonth, symDay) - 1)
}

func FixedToSymYear(fixedDate int) (int, int) {
  cycleMeanYear := 365.0 + 71.0 / 293.0
  symYear := int(math.Ceil((float64(fixedDate) - SymEpoch) / cycleMeanYear))
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
  weekOfYear := int(math.Ceil(float64(dayOfYear) / 7.0))
  quarter := int(math.Ceil((4.0 / 53.0) * float64(weekOfYear)))
  dayOfQuarter := dayOfYear - 91 * (quarter - 1)
  weekOfQuarter := int(math.Ceil(float64(dayOfQuarter) / 7.0))
  monthOfQuarter := int(math.Ceil((2.0 / 9.0) * float64(weekOfQuarter)))
  if monthOfQuarter > 3 {
    monthOfQuarter = 3
  }
  symMonth := 3 * quarter + monthOfQuarter - 3
  symDay := dayOfYear - SymDaysBeforeMonth(symMonth)

  return symYear, symMonth, symDay
}

func FixedToWeekdayNum(fixedDate int) (int) {
  weekdayAdjust := modulus(SymEpoch - 1, 7.0)
  return int(modulus(float64(fixedDate) - weekdayAdjust, 7.0))
}

func FromTime(t time.Time) (Sym) {
  fixedDate := GregorianToFixedDate(t.Year(), int(t.Month()), t.Day())
  symYear, symMonth, symDay := FixedToSym(fixedDate)

  sym := Sym{
    Year: symYear,
    Month: symMonth,
    Day: symDay,
  }

  return sym
}

func (sym Sym) ToTime() (time.Time) {
  fixedDate := SymToFixed(sym.Year, sym.Month, sym.Day)
  gregYear, gregMonth, gregDay := FixedToGregorian(fixedDate)

  return time.Date(
    gregYear, time.Month(gregMonth), gregDay, 0, 0, 0, 0, time.UTC)
}
