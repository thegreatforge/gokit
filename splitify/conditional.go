package splitify

type ConditionalSplit struct {
	rules          []*Rule
	defaultHandler Handler
}

func NewConditionalSplit(defaultHandler Handler) *ConditionalSplit {
	return &ConditionalSplit{
		defaultHandler: defaultHandler,
	}
}

func (c *ConditionalSplit) AddRule(rule *Rule) error {
	if rule.Conditions == nil {
		return ErrNoConditions
	}
	c.rules = append(c.rules, rule)
	return nil
}

func (c *ConditionalSplit) RemoveAllRules() error {
	c.rules = c.rules[:0]
	return nil
}

func (c *ConditionalSplit) Next(arg interface{}) (Handler, error) {
	if c.rules == nil {
		return nil, ErrNoRules
	}
	for _, rule := range c.rules {
		for _, condition := range rule.Conditions {
			if ok, err := condition(arg); err != nil {
				return c.defaultHandler, err
			} else if ok {
				return rule.Handler, nil
			}
		}
	}
	return c.defaultHandler, nil
}
