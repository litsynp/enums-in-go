package enum

import (
	"fmt"
	"strings"
)

type IntWeekday int32

const (
	IntMonday    IntWeekday = iota
	IntTuesday   IntWeekday = iota
	IntWednesday IntWeekday = iota
	IntThursday  IntWeekday = iota
	IntFriday    IntWeekday = iota
	IntSaturday  IntWeekday = iota
	IntSunday    IntWeekday = iota
)

type StrWeekday string

const (
	StrMonday    StrWeekday = "Monday"
	StrTuesday   StrWeekday = "Tuesday"
	StrWednesday StrWeekday = "Wednesday"
	StrThursday  StrWeekday = "Thursday"
	StrFriday    StrWeekday = "Friday"
	StrSaturday  StrWeekday = "Saturday"
	StrSunday    StrWeekday = "Sunday"
)

// Run `go generate ./...` to generate the stringer code
//
//go:generate stringer -type=DayOfWeek -trimprefix=DayOfWeek -output=generated_day_of_week.go
type DayOfWeek int32

const (
	DayOfWeekMonday    DayOfWeek = iota
	DayOfWeekTuesday   DayOfWeek = iota
	DayOfWeekWednesday DayOfWeek = iota
	DayOfWeekThursday  DayOfWeek = iota
	DayOfWeekFriday    DayOfWeek = iota
	DayOfWeekSaturday  DayOfWeek = iota
	DayOfWeekSunday    DayOfWeek = iota
)

func DayOfWeekFromString(s string) (DayOfWeek, error) {
	switch strings.ToLower(s) {
	case "monday":
		return DayOfWeekMonday, nil
	case "tuesday":
		return DayOfWeekTuesday, nil
	case "wednesday":
		return DayOfWeekWednesday, nil
	case "thursday":
		return DayOfWeekThursday, nil
	case "friday":
		return DayOfWeekFriday, nil
	case "saturday":
		return DayOfWeekSaturday, nil
	case "sunday":
		return DayOfWeekSunday, nil
	default:
		return -1, fmt.Errorf("invalid day of week: %s", s)
	}
}

func (d *DayOfWeek) UnmarshalJSON(data []byte) error {
	dayStr := strings.Trim(string(data), "\"")
	day, err := DayOfWeekFromString(dayStr)
	if err != nil {
		return err
	}

	*d = day
	return nil
}
