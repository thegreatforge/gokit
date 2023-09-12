package splitify

import "testing"

func TestConditional_Next(t *testing.T) {
	c := NewConditionalSplit("default")

	c.AddRule(&Rule{
		Handler: "a",
		Conditions: []Condition{
			func() (bool, error) {
				return true, nil
			},
		},
	})

	c.AddRule(&Rule{
		Handler: "b",
		Conditions: []Condition{
			func() (bool, error) {
				return false, nil
			},
		},
	})

	resp := make([]string, 10)
	for i := 0; i < 10; i++ {
		h, err := c.Next()
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
