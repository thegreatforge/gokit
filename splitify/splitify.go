package splitify

type Condition func() (bool, error)

type Handler interface{}

type RuleConfig struct {
	Handler    Handler
	Weight     int
	Conditions []Condition
}

type ISplitify interface {
	AddRule(*RuleConfig) error
	RemoveAllRules() error
	Next() (Handler, error)
}
