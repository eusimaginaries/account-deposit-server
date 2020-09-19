package app

import (
	"errors"
)

// Portfolio stores the state of the product
type Portfolio struct {
	Name    string
	Balance float32
}

// NewPortfolio instantiate and returns a portfolio witht the specified name.
func NewPortfolio(name string) (Portfolio, error) {
	if name == "" {
		return Portfolio{}, errors.New("invalid name")
	}
	return Portfolio{Name: name}, nil
}

// Deposit add to portfolio balance by the specified amount.
func (p *Portfolio) Deposit(amount float32) error {
	if amount < 0 {
		return errors.New("amount is negative")
	}
	p.Balance += amount
	return nil
}

// Withdraw subtract from portfolio balance by the specified amount.
func (p *Portfolio) Withdraw(amount float32) error {
	if amount < 0 {
		return errors.New("amount is negative")
	}
	if p.Balance < amount {
		return errors.New("withdrawal amount more than balance")
	}
	p.Balance -= amount
	return nil
}
