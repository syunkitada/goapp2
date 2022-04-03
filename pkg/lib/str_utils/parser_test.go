package str_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLastValue(t *testing.T) {
	a := assert.New(t)
	{
		result := ParseLastValue("voluntary_ctxt_switches:        14415")
		a.Equal("14415", result)
	}

	{
		result := ParseLastValue("voluntary_ctxt_switches:       \t14415")
		a.Equal("14415", result)
	}

	{
		result := ParseLastValue("14415")
		a.Equal("14415", result)
	}

	{
		result := ParseLastValue("")
		a.Equal("", result)
	}

	{
		result := ParseLastValue("14415 ")
		a.Equal("", result)
	}

	{
		result := ParseLastValue("14415\t")
		a.Equal("", result)
	}
}

func TestParseLastSecondValue(t *testing.T) {
	a := assert.New(t)
	{
		result := ParseLastSecondValue("VmSize:    15736 kB")
		a.Equal("15736", result)
	}

	{
		result := ParseLastSecondValue("15736 kB")
		a.Equal("15736", result)
	}

	{
		result := ParseLastSecondValue("15736")
		a.Equal("", result)
	}
}

func TestParseRangeFormatStr(t *testing.T) {
	a := assert.New(t)
	{
		result := ParseRangeFormatStr("0-1,3-5")
		a.Equal([]int{0, 1, 3, 4, 5}, result)
	}
	{
		result := ParseRangeFormatStr("3,5")
		a.Equal([]int{3, 5}, result)
	}
	{
		result := ParseRangeFormatStr("3")
		a.Equal([]int{3}, result)
	}
	{
		result := ParseRangeFormatStr("a-1,3-5")
		var expected []int
		a.Equal(expected, result)
	}
	{
		result := ParseRangeFormatStr("0-a,3-5")
		var expected []int
		a.Equal(expected, result)
	}
}

func TestSplitSpace(t *testing.T) {
	a := assert.New(t)
	{
		result := SplitSpace("VmSize:    15736 kB")
		a.Equal([]string{"VmSize:", "15736", "kB"}, result)
	}
	{
		result := SplitSpace("VmSize:\t15736 kB")
		a.Equal([]string{"VmSize:", "15736", "kB"}, result)
	}
	{
		result := SplitSpace("VmSize:  \t  15736 kB")
		a.Equal([]string{"VmSize:", "15736", "kB"}, result)
	}
}

func TestSplitColon(t *testing.T) {
	a := assert.New(t)
	{
		result := SplitColon("   com-1-ex:    1426    245  ")
		a.Equal([]string{"com-1-ex", "    1426    245  "}, result)
	}
	{
		result := SplitColon("\t   com-1-ex: \t 1426 \t 245  ")
		a.Equal([]string{"com-1-ex", " \t 1426 \t 245  "}, result)
	}
	{
		result := SplitColon("")
		var expected []string
		a.Equal(expected, result)
	}
}

func TestSplitSpaceColon(t *testing.T) {
	a := assert.New(t)
	{
		result := SplitSpaceColon("core id         : 6")
		a.Equal([]string{"core id", "6"}, result)
	}
	{
		result := SplitSpaceColon("core id\t: 6")
		a.Equal([]string{"core id", "6"}, result)
	}
	{
		result := SplitSpaceColon("")
		var expected []string
		a.Equal(expected, result)
	}
}
