package helper

import "strings"

func FirstToUpper(input string) string {
	if input == "" {
		return ""
	}
	tmp := []byte(input)
	first := tmp[0]
	if first > 96 && first < 123 {
		tmp[0] = first - 32
		return string(tmp)
	}
	return input
}

func FirstToLower(input string) string {
	if input == "" {
		return ""
	}
	tmp := []byte(input)
	first := tmp[0]
	if first > 64 && first < 91 {
		tmp[0] = first + 32
		return string(tmp)
	}
	return input
}

func CamelToLine(str string) string {
	splits := []rune(str)
	join := ""
	for index, split := range splits {
		chart := string(split)
		lower := strings.ToLower(chart)
		if lower != chart {
			if index == 0 {
				join = lower
			} else {
				join += "-" + lower
			}
		} else {
			join += lower
		}
	}
	return join
}
