package app

import (
	"errors"
)

// Customer is the portfolio owner
type Customer struct {
	portfolios []*Portfolio
}

// PerformDeposit split the passed in deposit into the respective portfolio
func (c *Customer) PerformDeposit(depositPlans []DepositPlan, deposits []float32) error {
	var totalDeposit float32
	for _, v := range deposits {
		totalDeposit += v
	}

	var totalNeeded float32
	for _, dp := range depositPlans {
		ok := true
		for k := range dp.PortfolioRatio() {
			var exist bool
			for _, p := range c.portfolios {
				if k == p.Name {
					exist = true
					break
				}
			}
			if !exist {
				ok = false
				break
			}
		}

		if !ok {
			return errors.New("deposit plan does not match customer portfolio")
		}

		totalNeeded += dp.DepositTotal()
	}

	if totalNeeded != totalDeposit {
		return errors.New("deposits does not match the plan amounts")
	}

	for _, dp := range depositPlans {
		for _, p := range c.portfolios {
			v := dp.PortfolioRatio()[p.Name]
			p.Deposit(v)
		}
	}

	return nil
}
