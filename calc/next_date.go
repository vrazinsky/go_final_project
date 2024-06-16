package calc

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	TypeYear  = "y"
	TypeDay   = "d"
	TypeWeek  = "w"
	TypeMonth = "m"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}
	t, arr1, arr2, err := parseRepeat(repeat)
	if err != nil {
		return "", err
	}
	if t == TypeYear {
		cntYears := 1
		for {
			newDate := startDate.AddDate(cntYears, 0, 0)
			if newDate.After(now) {
				return newDate.Format("20060102"), nil
			}
			cntYears++
		}
	}
	if t == TypeDay {
		daysDiff := arr1[0]
		cnt := 1
		for {
			newDate := startDate.AddDate(0, 0, daysDiff*cnt)
			if newDate.After(now) || newDate.Format("20060102") == now.Format("20060102") {
				return newDate.Format("20060102"), nil
			}
			cnt++
		}
	}
	if t == TypeWeek {
		cnt := 1
		for {
			newDate := now.AddDate(0, 0, cnt)
			if slices.Contains(arr1, int(newDate.Weekday())) {
				return newDate.Format("20060102"), nil
			}
			cnt++
			if cnt > 100000 {
				return "", errors.New("error with w")
			}
		}
	}
	if t == TypeMonth {
		cnt := 1
		var date time.Time
		if now.After(startDate) {
			date = now
		} else {
			date = startDate
		}
		for {
			newDate := date.AddDate(0, 0, cnt)
			firstDay := time.Date(newDate.Year(), newDate.Month(), 1, 0, 0, 0, 0, time.UTC)
			maxDays := firstDay.AddDate(0, 1, 0).Add(-time.Nanosecond).Day()

			if arr2 == nil || slices.Contains(arr2, int(newDate.Month())) {
				if slices.Contains(arr1, int(newDate.Day())) {
					return newDate.Format("20060102"), nil
				}
				if slices.Contains(arr1, -1) && newDate.Day() == maxDays {
					return newDate.Format("20060102"), nil
				}
				if slices.Contains(arr1, -2) && newDate.Day() == maxDays-1 {
					return newDate.Format("20060102"), nil
				}
			}
			cnt++
			if cnt > 100000 {
				return "", errors.New("error with month")
			}
		}
	}
	return "", nil
}

func parseRepeat(repeat string) (string, []int, []int, error) {
	arr := strings.Split(repeat, " ")
	t := arr[0]
	if t == TypeYear {
		return t, nil, nil, nil
	}

	if t == TypeDay {
		if len(arr) != 2 || len(arr[1]) == 0 {
			return "", nil, nil, errors.New("incorrect days repeat format")
		}
		n, err := strconv.Atoi(arr[1])
		if err != nil {
			return "", nil, nil, err
		}
		if n > 400 || n <= 0 {
			return "", nil, nil, errors.New("incorrect days value")
		}
		return t, []int{n}, nil, nil
	}
	if t == TypeWeek {
		if len(arr) != 2 || len(arr[1]) == 0 {
			return "", nil, nil, errors.New("incorrect week repeat format")
		}
		weekDays := strings.Split(arr[1], ",")
		weekSlice := make([]int, 0)
		for _, weekDay := range weekDays {
			weekDayInt, err := strconv.Atoi(weekDay)
			if err != nil {
				return "", nil, nil, err
			}
			if weekDayInt < 1 || weekDayInt > 7 {
				return "", nil, nil, errors.New("incorrect week  value")
			}
			if weekDayInt == 7 {
				weekDayInt = 0
			}
			if !slices.Contains(weekSlice, weekDayInt) {
				weekSlice = append(weekSlice, weekDayInt)
			}
		}
		return t, weekSlice, nil, nil
	}
	if t == TypeMonth {
		if !(len(arr) == 2 || len(arr) == 3) {
			return "", nil, nil, errors.New("incorrect month repeat format")
		}
		monthDays := strings.Split(arr[1], ",")
		monthDaysSlice := make([]int, 0)
		for _, monthDay := range monthDays {
			monthDayInt, err := strconv.Atoi(monthDay)
			if err != nil {
				return "", nil, nil, err
			}
			if (monthDayInt < -2 || monthDayInt > 31) || monthDayInt == 0 {
				return "", nil, nil, errors.New("incorrect month day value")
			}
			if !slices.Contains(monthDaysSlice, monthDayInt) {
				monthDaysSlice = append(monthDaysSlice, monthDayInt)
			}
		}
		if len(arr) == 2 {
			return t, monthDaysSlice, nil, nil
		}
		months := strings.Split(arr[2], ",")
		monthsSlice := make([]int, 0)
		for _, month := range months {
			monthInt, err := strconv.Atoi(month)
			if err != nil {
				return "", nil, nil, err
			}
			if monthInt < 1 || monthInt > 12 {
				return "", nil, nil, errors.New("incorrect month value")
			}
			if !slices.Contains(monthsSlice, monthInt) {
				monthsSlice = append(monthsSlice, monthInt)
			}
		}
		return t, monthDaysSlice, monthsSlice, nil
	}
	return "", nil, nil, errors.New("incorrect year repeat format")
}
