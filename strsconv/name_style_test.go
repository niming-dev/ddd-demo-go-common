package strsconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamel2Snake(t *testing.T) {
	testTable := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "test1",
			in:   "TestName",
			out:  "test_name",
		},
		{
			name: "test2",
			in:   "testName",
			out:  "test_name",
		},
		{
			name: "test3",
			in:   "testName3",
			out:  "test_name3",
		},
		{
			name: "test4",
			in:   "test4Name",
			out:  "test4_name",
		},
		{
			name: "test5",
			in:   "TestName在",
			out:  "test_name在",
		},
	}

	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.out, Camel2Snake(v.in))
		})
	}
}

func TestSnake2Camel(t *testing.T) {
	testTable := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "test1",
			in:   "test_name",
			out:  "testName",
		},
		{
			name: "test3",
			in:   "test_name3",
			out:  "testName3",
		},
		{
			name: "test4",
			in:   "test4_name",
			out:  "test4Name",
		},
		{
			name: "test5",
			in:   "TestName",
			out:  "TestName",
		},
		{
			name: "test5",
			in:   "testName",
			out:  "testName",
		},
	}

	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.out, Snake2Camel(v.in))
		})
	}
}

func TestSnake2Pascal(t *testing.T) {
	testTable := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "test1",
			in:   "test_name",
			out:  "TestName",
		},
		{
			name: "test2",
			in:   "test_name",
			out:  "TestName",
		},
		{
			name: "test3",
			in:   "test_name3",
			out:  "TestName3",
		},
		{
			name: "test4",
			in:   "test4_name",
			out:  "Test4Name",
		},
		{
			name: "test5",
			in:   "TestName",
			out:  "TestName",
		},
	}

	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.out, Snake2Pascal(v.in))
		})
	}
}

func TestCamel2Pascal(t *testing.T) {
	testTable := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "test1",
			in:   "testName",
			out:  "TestName",
		},
		{
			name: "test2",
			in:   "uaccountFullAccess",
			out:  "UaccountFullAccess",
		},
		{
			name: "test2",
			in:   "u",
			out:  "U",
		},
	}

	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.out, Camel2Pascal(v.in))
		})
	}
}

func TestPascal2Camel(t *testing.T) {
	testTable := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "test1",
			in:   "TestName",
			out:  "testName",
		},
		{
			name: "test2",
			in:   "UaccountFullAccess",
			out:  "uaccountFullAccess",
		},
		{
			name: "test3",
			in:   "A",
			out:  "a",
		},
		{
			name: "test4",
			in:   "a",
			out:  "a",
		},
		{
			name: "test4",
			in:   "a1a",
			out:  "a1a",
		},
	}

	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.out, Pascal2Camel(v.in))
		})
	}
}

func TestConvertNameStyle(t *testing.T) {
	testTable := []struct {
		name      string
		srcStyle  NameStyle
		destStyle NameStyle
		in        string
		out       string
	}{
		{
			name:      "test1",
			srcStyle:  NameStylePascal,
			destStyle: NameStyleCamel,
			in:        "TestName",
			out:       "testName",
		},
		{
			name:      "test2",
			srcStyle:  NameStylePascal,
			destStyle: NameStyleSnake,
			in:        "UaccountFullAccess",
			out:       "uaccount_full_access",
		},
		{
			name:      "test3",
			srcStyle:  NameStyleSnake,
			destStyle: NameStylePascal,
			in:        "aa_bb_cc_dd",
			out:       "AaBbCcDd",
		},
		{
			name:      "test4",
			srcStyle:  NameStyleSnake,
			destStyle: NameStyleCamel,
			in:        "aa_bb_cc_dd",
			out:       "aaBbCcDd",
		},
		{
			name:      "test5",
			srcStyle:  NameStyleCamel,
			destStyle: NameStyleSnake,
			in:        "aaBbCcDd",
			out:       "aa_bb_cc_dd",
		},
		{
			name: "test6",
			in:   "aaBbCcDd",
			out:  "aaBbCcDd",
		},
	}

	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.out, ConvertNameStyle(v.in, v.srcStyle, v.destStyle))
		})
	}
}

func TestNameStyleAnyToOne(t *testing.T) {
	testTable := []struct {
		name      string
		destStyle NameStyle
		in        string
		out       string
	}{
		{
			name:      "test1",
			destStyle: NameStyleCamel,
			in:        "TestName",
			out:       "testName",
		},
		{
			name:      "test2",
			destStyle: NameStyleSnake,
			in:        "UaccountFullAccess",
			out:       "uaccount_full_access",
		},
		{
			name:      "test3",
			destStyle: NameStylePascal,
			in:        "aa_bb_cc_dd",
			out:       "AaBbCcDd",
		},
		{
			name:      "test4",
			destStyle: NameStyleCamel,
			in:        "aa_bb_cc_dd",
			out:       "aaBbCcDd",
		},
		{
			name:      "test5",
			destStyle: NameStyleSnake,
			in:        "aaBbCcDd",
			out:       "aa_bb_cc_dd",
		},
		{
			name: "test6",
			in:   "aaBbCcDd",
			out:  "aaBbCcDd",
		},
	}

	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.out, NameStyleAnyToOne(v.in, v.destStyle))
		})
	}
}

func TestCheckNameStyle(t *testing.T) {
	testTable := []struct {
		name  string
		in    string
		style NameStyle
	}{
		{
			name:  "test1",
			style: NameStylePascal,
			in:    "TestName",
		},
		{
			name:  "test2",
			style: NameStyleCamel,
			in:    "uaccountAccess",
		},
		{
			name:  "test3",
			style: NameStyleSnake,
			in:    "aa_bb_cc_dd",
		},
		{
			name:  "test4",
			style: NameStyleUnspecified,
			in:    "",
		},
	}

	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.style, CheckNameStyle(v.in))
		})
	}
}

func BenchmarkCamel2Snake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Camel2Snake("TestName")
	}
}

func BenchmarkSnake2Camel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Snake2Pascal("test_name")
	}
}
