package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	t.Log("TestStringCheck")
	str := []string{"I", "am", "stupid", "and", "weak"}
	Target := []string{"I", "am", "smart", "and", "strong"}
	result1 := GerneralForLoop(str)
	result2 := IndexForRangeLoop(str)
	result3 := ValueForRangeLoop(str)
	result4 := ByMap(str)

	// Assert check
	assert.Equal(t, result1, Target)
	assert.Equal(t, result2, Target)
	assert.Equal(t, result3, Target)
	assert.Equal(t, result4, Target)
}
