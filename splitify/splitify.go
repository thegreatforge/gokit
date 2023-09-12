package splitify

type Condition func(interface{}) (bool, error)

type Handler interface{}

type Rule struct {
	Handler    Handler
	Weight     int
	Conditions []Condition
}
