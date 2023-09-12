package splitify

type Condition func() (bool, error)

type Handler interface{}

type Rule struct {
	Handler    Handler
	Weight     int
	Conditions []Condition
}

type ISplitify interface {
	AddRule(*Rule) error
	RemoveAllRules() error
	Next() (Handler, error)
}
