package app

import (
	"errors"
)

// DepositSession allows customer to perform a transaction
type DepositSession struct {
	depositPlans []DepositPlan
	deposits     []float32
}

// Customer is the portfolio owner
type Customer struct {
	ID             string
	portfolios     []*Portfolio
	DepositSession *DepositSession
}

// NewCustomer instantiate a new customer with no portfolios
func NewCustomer(id string) Customer {
	return Customer{ID: id, portfolios: []*Portfolio{}}
}

// AddPortfolio adds portfolio after determining validity
func (c *Customer) AddPortfolio(name string) error {
	for _, p := range c.portfolios {
		if p.Name == name {
			return errors.New("portfolio with specfied name already added")
		}
	}

	newP, err := NewPortfolio(name)
	if err != nil {
		return err
	}
	c.portfolios = append(c.portfolios, &newP)
	return nil
}

// StartSession starts a deposit session
func (c *Customer) StartSession() error {
	if c.DepositSession != nil {
		return errors.New("another transaction is still active")
	}
	c.DepositSession = &DepositSession{depositPlans: []DepositPlan{}, deposits: []float32{}}
	return nil
}

// Deposit represents the amount the customer has deposit
func (c *Customer) Deposit(amount float32) error {
	if c.DepositSession == nil {
		return errors.New("no active session")
	}
	if amount < 0 {
		return errors.New("amount is negative")
	}
	c.DepositSession.deposits = append(c.DepositSession.deposits, amount)
	return nil
}

// PayDepositPlan represnts the plans the customer would like to pay
func (c *Customer) PayDepositPlan(plan DepositPlan) error {
	if c.DepositSession == nil {
		return errors.New("no active session")
	}

	for _, dp := range c.DepositSession.depositPlans {
		if dp.Name() == plan.Name() {
			return errors.New("duplicate plan name in session")
		}
	}

	c.DepositSession.depositPlans = append(c.DepositSession.depositPlans, plan)
	return nil
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
