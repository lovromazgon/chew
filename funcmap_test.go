package chew

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndent_Expected(t *testing.T) {
	testCases := []struct {
		IndentSize int
		Str        string
		Expected   string
	}{
		{1, "", " "},
		{0, "", ""},
		{5, "test", "     test"},
		{0, "   test", "   test"},
		{3, "multi\nline\nstring", "   multi\n   line\n   string"},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.Expected, Indent(tc.IndentSize, tc.Str))
	}
}

func TestIndent_Panic(t *testing.T) {
	testCases := []struct {
		IndentSize int
		Str        string
	}{
		{-1, ""},
		{-1, "test"},
	}

	for _, tc := range testCases {
		assert.Panics(t, func() {
			Indent(tc.IndentSize, tc.Str)
		})
	}
}

func TestOffset_Expected(t *testing.T) {
	testCases := []struct {
		OffsetSize int
		Str        interface{}
		Expected   string
	}{
		{0, "", ""},
		{1, "", " "},
		{5, "12345", ""},
		{5, "test", " "},
		{5, 12345, ""},
		{5, 1234, " "},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.Expected, Offset(tc.OffsetSize, tc.Str))
	}
}

func TestOffset_Panic(t *testing.T) {
	testCases := []struct {
		OffsetSize int
		Str        string
	}{
		{-1, ""},
		{3, "test"},
	}

	for _, tc := range testCases {
		assert.Panics(t, func() {
			Offset(tc.OffsetSize, tc.Str)
		})
	}
}

func TestMaxLength_Expected(t *testing.T) {
	testCases := []struct {
		Slice    []string
		Expected int
	}{
		{[]string{}, 0},
		{[]string{""}, 0},
		{[]string{"1", "2"}, 1},
		{[]string{"a", "ab", "abc"}, 3},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.Expected, MaxLength(tc.Slice))
	}
}

func TestMaxLength_Panic(t *testing.T) {
	testCases := []struct {
		Slice interface{}
	}{
		{nil},
		{[]interface{}{}},
		{ 1},
		{"test"},
	}

	for _, tc := range testCases {
		assert.Panics(t, func() {
			MaxLength(tc.Slice)
		})
	}

}
