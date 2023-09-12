package splitify

type ConditionalSplit struct {
	rules          []*RuleConfig
	defaultHandler Handler
}

func NewConditionalSplit(defaultHandler Handler) ISplitify {
	return &ConditionalSplit{
		defaultHandler: defaultHandler,
	}
}

func (c *ConditionalSplit) AddRule(rule *RuleConfig) error {
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

func (c *ConditionalSplit) Next() (Handler, error) {
	if c.rules == nil {
		return nil, ErrNoRules
	}
	for _, rule := range c.rules {
		for _, condition := range rule.Conditions {
			if ok, err := condition(); err != nil {
				return c.defaultHandler, err
			} else if ok {
				return rule.Handler, nil
			}
		}
	}
	return c.defaultHandler, nil
}
