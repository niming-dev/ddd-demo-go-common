package test

import (
	"log"
	"strings"
	"testing"

	"github.com/niming-dev/ddd-demo/go-common/expression"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type testContext2 struct {
}

func (c *testContext2) Get(name string) (*expression.Data, error) {
	switch name {
	case "ten":
		return expression.NewInt(10), nil
	case "one":
		return expression.NewInt(1), nil
	case "five":
		return expression.NewInt(5), nil
	default:
		return expression.NewInt(0), nil
	}
}

func (c *testContext2) Call(name string, args []*expression.Data) (*expression.Data, error) {
	return expression.NewString("lisi"), nil
}

type TestRule struct {
	rule   string
	expect string
}

func Test_DeclareFunc(t *testing.T) {
	rules := []TestRule{
		{
			rule: `$func() {
			return "15";
		}`,
			expect: "15",
		},

		{
			rule: `$func() {
				switch ( ${ten} ) {
					default: {
						return "1";
					}
				}
			}`,
			expect: "1",
		},

		{
			rule: `$func() {
				switch ( ${ten} ) {
					case 1: {
						return "15";
					}
				}
				return "0";
			}`,
			expect: "0",
		},

		{
			rule: `$func() {
				switch ( ${ten} ) {
					case 1: {
						return "15";
					}
					default: {
						return "1";
					}
				}
				return "0";
			}`,
			expect: "1",
		},

		{
			rule: `$func() {
				switch ( ${ten} ) {
					case 10: {
						return "10";
					}
					default: {
						return "1";
					}
				}
				return "0";
			}`,
			expect: "10",
		},

		{
			rule: `$func() {
				switch ( ${ten} ) {
					case <=15: {
						return "15";
					}
					case <=20: {
						return "20";
					}
				}
			}`,
			expect: "15",
		},

		{
			rule: `$func() {
				switch ( ${ten} ) {
					case <=5: {
						return "5";
					}
					case <=20: {
						return "20";
					}
				}
			}`,
			expect: "20",
		},
	}

	for i, rule := range rules {
		exec, err := expression.Parse(strings.NewReader(rule.rule))
		if nil != err {
			t.Fatal(err, rule)
		}

		if i == 6 {
			log.Println(exec.Dump(0))
		}
		data, err := exec.Evaluate(&testContext2{})
		if nil != err {
			t.Fatal(err, data, rule)
		}
		b, err := expression.Equal(data, expression.NewString(rule.expect))
		if nil != err {
			t.Fatal(err, data, rule)
		}
		if !b.Bool() {
			t.Fatal("not equal to expect", data, rule)
		}

		log.Println("------------- testcase", i, "ok -------------")
	}
}
