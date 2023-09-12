package splitify

import "sync"

type WeightedSplit struct {
	rules         []*Rule
	rulesCount    int
	gcd           int
	maxWeight     int
	i             int
	currentWeight int
	mutex         sync.Mutex
}

func NewWeightedSplit() ISplitify {
	return &WeightedSplit{}
}

func (w *WeightedSplit) AddRule(rule *Rule) error {
	if rule.Weight > 0 {
		if w.gcd == 0 {
			w.gcd = rule.Weight
			w.maxWeight = rule.Weight
			w.i = -1
			w.currentWeight = 0
		} else {
			w.gcd = gcd(w.gcd, rule.Weight)
			if w.maxWeight < rule.Weight {
				w.maxWeight = rule.Weight
			}
		}
	}
	w.rules = append(w.rules, rule)
	w.rulesCount++
	return nil
}

func (w *WeightedSplit) RemoveAllRules() error {
	w.rules = w.rules[:0]
	w.rulesCount = 0
	w.gcd = 0
	w.maxWeight = 0
	w.i = -1
	w.currentWeight = 0
	return nil
}

func (w *WeightedSplit) Next() (Handler, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.rulesCount == 0 {
		return nil, ErrNoRules
	}

	if w.rulesCount == 1 {
		return w.rules[0].Handler, nil
	}

	for {
		w.i = (w.i + 1) % w.rulesCount
		if w.i == 0 {
			w.currentWeight = w.currentWeight - w.gcd
			if w.currentWeight <= 0 {
				w.currentWeight = w.maxWeight
				if w.currentWeight == 0 {
					return nil, nil
				}
			}
		}

		if w.rules[w.i].Weight >= w.currentWeight {
			return w.rules[w.i].Handler, nil
		}
	}
}

func gcd(x, y int) int {
	var t int
	for {
		t = (x % y)
		if t > 0 {
			x = y
			y = t
		} else {
			return y
		}
	}
}
