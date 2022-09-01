package helper

import (
	"math"
	"strconv"
	"strings"
)

func divide(dividend int64, divisor int64) (int64, int64) {
	quotient := dividend / divisor
	remainder := dividend % divisor
	return quotient, remainder
}

func SecondToTime(seconds float64) string {
	remainder := int64(0)
	quotient := int64(math.Round(seconds))
	time := ""
	if quotient >= 60 {
		for quotient >= 60 {
			quotient, remainder = divide(quotient, 60)
			strRemainder := strconv.FormatInt(remainder, 10)
			if remainder < 10 {
				strRemainder = "0" + strRemainder
			}
			time = strRemainder + time
		}
	} else {
		time = strconv.FormatInt(quotient, 10)
		if quotient < 10 {
			time = "0" + time
		}
	}
	return time
}

func TimeToSecond(time string) float64 {
	splits := strings.Split(time, ":")
	seconds := 0.0
	for _, split := range splits {
		depart, _ := strconv.ParseFloat(split, 64)
		seconds = seconds + depart
	}
	return seconds
}
