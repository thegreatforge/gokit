package splitify

import (
	"errors"
	"testing"
)

func TestConditional_Next(t *testing.T) {
	c := NewConditionalSplit("default")

	splitterConditionofA := func(arg interface{}) (bool, error) {
		b, ok := arg.(bool)
		if !ok {
			t.Fatal("invalid argument")
			return false, errors.New("invalid argument")
		}
		if b {
			return true, nil
		}
		return false, nil
	}

	splitterConditionofB := func(arg interface{}) (bool, error) {
		b, ok := arg.(bool)
		if !ok {
			t.Fatal("invalid argument")
			return false, errors.New("invalid argument")
		}
		if b {
			return false, nil
		}
		return true, nil
	}

	c.AddRule(&Rule{
		Handler: "a",
		Conditions: []Condition{
			splitterConditionofA,
		},
	})

	c.AddRule(&Rule{
		Handler: "b",
		Conditions: []Condition{
			splitterConditionofB,
		},
	})

	x := true
	resp := make([]string, 10)
	for i := 0; i < 10; i++ {
		h, err := c.Next(x)
		if err != nil {
			t.Fatal(err)
		}
		resp[i] = h.(string)
	}

	for _, v := range resp {
		if v != "a" {
			t.Fatalf("expected 'a', got %s", v)
		}
	}
}
