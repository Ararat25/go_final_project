package model

import (
	"github.com/Ararat25/go_final_project/customError"
	"strconv"
	"strings"
	"time"
)

var period = []string{"d", "w", "m", "y"}

type Repeat struct {
	Period      string
	FirstSlice  []int
	SecondSlice []int
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", customError.ErrRepeatNotSpecified
	}

	parseDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	parseRepeatRule, err := parseRepeat(repeat)
	if err != nil {
		return "", err
	}

	switch parseRepeatRule.Period {
	case period[0]:
		return dParse(now, parseDate, parseRepeatRule)
	case period[1]:
		return wParse(now, parseDate, parseRepeatRule)
	case period[2]:
		return mParse(now, parseDate, parseRepeatRule)
	case period[3]:
		return yParse(now, parseDate, parseRepeatRule)
	default:
		return "", nil
	}
}

func parseRepeat(repeat string) (Repeat, error) {
	repeatSlice := strings.Split(repeat, " ")

	repeatPeriod := repeatSlice[0]

	validRepeatRule := false
	for i := 0; i < len(period); i++ {
		if repeatPeriod == period[i] {
			validRepeatRule = true
		}
	}

	if !validRepeatRule {
		return Repeat{}, customError.ErrNotValidRepeat
	}

	if len(repeatSlice) < 2 {
		return Repeat{
			Period: repeatPeriod,
		}, nil
	}

	firstStrSlice := strings.Split(repeatSlice[1], ",")

	firstIntSlice, err := sliceStringToInt(firstStrSlice)
	if err != nil {
		return Repeat{}, err
	}

	if len(repeatSlice) < 3 {
		return Repeat{
			Period:     repeatPeriod,
			FirstSlice: firstIntSlice,
		}, nil
	}

	secondStrSlice := strings.Split(repeatSlice[2], ",")

	secondIntSlice, err := sliceStringToInt(secondStrSlice)
	if err != nil {
		return Repeat{}, err
	}

	return Repeat{
		Period:      repeatPeriod,
		FirstSlice:  firstIntSlice,
		SecondSlice: secondIntSlice,
	}, nil
}

func sliceStringToInt(strSlice []string) ([]int, error) {
	var intSlice []int
	for _, str := range strSlice {
		num, err := strconv.Atoi(str)
		if err != nil {
			return []int{}, err
		}

		intSlice = append(intSlice, num)
	}
	return intSlice, nil
}

func isValidDay(nextDate time.Time, dayMap map[int]bool) bool {
	dInMonth := daysInMonth(nextDate)
	day := nextDate.Day()
	return dayMap[day] || (day == dInMonth && dayMap[-1]) || (day == dInMonth-1 && dayMap[-2])
}

// daysInMonth возвращает количество дней в месяце для заданной даты
func daysInMonth(t time.Time) int {
	year, month := t.Year(), t.Month()
	return time.Date(year, month+1, 0, 0, 0, 0, 0, t.Location()).Day()
}

func dParse(now time.Time, date time.Time, repeat Repeat) (string, error) {
	if len(repeat.FirstSlice) != 1 {
		return "", customError.ErrNotValidRepeat
	}

	interval := repeat.FirstSlice[0]
	if interval > 400 || interval < 1 {
		return "", customError.ErrNotValidRepeat
	}

	nextDate := date.AddDate(0, 0, interval)
	for nextDate.Before(now) {
		nextDate = nextDate.AddDate(0, 0, interval)
	}

	return nextDate.Format("20060102"), nil
}

func wParse(now time.Time, date time.Time, repeat Repeat) (string, error) {
	if len(repeat.FirstSlice) < 1 || len(repeat.FirstSlice) > 7 {
		return "", customError.ErrNotValidRepeat
	}

	interval := repeat.FirstSlice
	for _, elem := range interval {
		if elem < 1 || elem > 7 {
			return "", customError.ErrNotValidRepeat
		}
	}

	weekDays := make(map[int]bool, 7)
	for i := 0; i < 7; i++ {
		weekDays[i] = false
	}

	for _, day := range interval {
		if day == 7 {
			weekDays[0] = true
			continue
		}
		weekDays[day] = true
	}

	nextDate := date.AddDate(0, 0, 1)

	for {
		if weekDays[int(nextDate.Weekday())] {
			if nextDate.After(now) {
				break
			}
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}

	return nextDate.Format("20060102"), nil
}

func mParse(now time.Time, date time.Time, repeat Repeat) (string, error) {
	if len(repeat.FirstSlice) == 0 {
		return "", customError.ErrNotValidRepeat
	}

	additionalParameters := false
	if len(repeat.SecondSlice) != 0 {
		additionalParameters = true
	}

	days := repeat.FirstSlice
	for _, day := range days {
		if day < -2 || day > 31 || day == 0 {
			return "", customError.ErrNotValidRepeat
		}
	}

	var months []int
	if additionalParameters {
		months = repeat.SecondSlice
		for _, m := range months {
			if m < 1 || m > 12 {
				return "", customError.ErrNotValidRepeat
			}
		}
	}

	dayMap := make(map[int]bool, 31)
	for i := -2; i <= 31; i++ {
		if i == 0 {
			continue
		}
		dayMap[i] = false
	}

	for _, day := range days {
		dayMap[day] = true
	}

	nextDate := date.AddDate(0, 0, 1)

	if !additionalParameters {
		for {
			if isValidDay(nextDate, dayMap) {
				if nextDate.After(now) {
					break
				}
			}
			nextDate = nextDate.AddDate(0, 0, 1)
		}
	} else {
	loop:
		for {
			flag := false
			if isValidDay(nextDate, dayMap) {
				if nextDate.After(now) {
					flag = true
				}
			}

			for _, month := range months {
				if int(nextDate.Month()) == month && flag {
					break loop
				}
			}

			nextDate = nextDate.AddDate(0, 0, 1)
		}
	}

	return nextDate.Format("20060102"), nil
}

func yParse(now time.Time, date time.Time, repeat Repeat) (string, error) {
	if len(repeat.FirstSlice) != 0 {
		return "", customError.ErrNotValidRepeat
	}

	nextDate := date.AddDate(1, 0, 0)
	for nextDate.Before(now) {
		nextDate = nextDate.AddDate(1, 0, 0)
	}

	return nextDate.Format("20060102"), nil
}
