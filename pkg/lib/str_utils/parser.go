package str_utils

import (
	"strconv"
	"strings"
)

func ParseLastValue(s string) string {
	// Parse 'voluntary_ctxt_switches:        14415', and return 14415
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ' ' || s[i:i+1] == "\t" {
			return s[i+1 : len(s)]
		}
	}
	return s
}

func ParseLastSecondValue(s string) string {
	// Parse 'VmSize:    15736 kB' and return 1536
	lastIndex := 0
	i := len(s) - 1
	for ; i >= 0; i-- {
		if s[i] == ' ' || s[i:i+1] == "\t" {
			if lastIndex == 0 || lastIndex == i+1 {
				lastIndex = i
			} else {
				return s[i+1 : lastIndex]
			}
		}
	}
	if lastIndex != 0 {
		return s[i+1 : lastIndex]
	}
	return ""
}

func ParseRangeFormatStr(str string) (ints []int) {
	// parse 0-1,3-5, and return [0,1,3,4,5]
	splitedComma := strings.Split(strings.TrimRight(str, "\n"), ",")
	for _, c := range splitedComma {
		splitedRange := strings.Split(c, "-")
		firstNum, tmpErr := strconv.Atoi(splitedRange[0])
		if tmpErr != nil {
			return
		}
		if len(splitedRange) == 1 {
			ints = append(ints, firstNum)
			continue
		}
		secondNum, tmpErr := strconv.Atoi(splitedRange[1])
		if tmpErr != nil {
			return
		}
		for ; firstNum <= secondNum; firstNum++ {
			ints = append(ints, firstNum)
		}
	}
	return
}

func SplitSpace(str string) (strs []string) {
	lenstr := len(str)
	tmpBytes := []byte{}
	for i := 0; i < lenstr; i++ {
		if str[i] == ' ' {
			if len(tmpBytes) > 0 {
				strs = append(strs, string(tmpBytes))
				tmpBytes = []byte{}
			}
			continue
		}
		if i+1 < lenstr && str[i:i+1] == "\t" {
			if len(tmpBytes) > 0 {
				strs = append(strs, string(tmpBytes))
				tmpBytes = []byte{}
			}
			continue
		}
		tmpBytes = append(tmpBytes, str[i])
	}

	if len(tmpBytes) > 0 {
		strs = append(strs, string(tmpBytes))
	}
	return
}

func SplitColon(str string) (strs []string) {
	// parse '   com-1-ex:    1426    245  ', and return ["com-1-ex", "    1426    245  "]

	lenstr := len(str)
	lastSpaceIndex := 0
	lastStrIndex := 0
	for i := 0; i < lenstr; i++ {
		if str[i] == ' ' {
			lastSpaceIndex = i + 1
			continue
		} else if i+1 < lenstr && str[i:i+1] == "\t" {
			i++
			lastSpaceIndex = i + 1
			continue
		} else if str[i] == ':' {
			strs = append(strs, str[lastSpaceIndex:lastStrIndex+1])
			if i+1 < lenstr {
				strs = append(strs, str[i+1:lenstr])
			}
			return
		}
		lastStrIndex = i
	}

	return
}

func SplitSpaceColon(str string) (strs []string) {
	// parse 'core id         : 6', and return ["core id", "6"]

	lenstr := len(str)
	lastStrIndex := 0
	for i := 0; i < lenstr; i++ {
		if str[i] == ' ' {
			continue
		} else if str[i] == ':' {
			strs = append(strs, str[0:lastStrIndex+1])
			if i+2 < lenstr {
				strs = append(strs, str[i+2:lenstr])
			}
			return
		} else if i+1 < lenstr && str[i:i+1] == "\t" {
			continue
		}
		lastStrIndex = i
	}

	return
}
